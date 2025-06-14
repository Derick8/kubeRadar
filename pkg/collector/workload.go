package collector

import (
	"context"
	"fmt"
	"kubeRadar/pkg/models"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Collector) collectWorkloadInfo(ctx context.Context) (models.WorkloadAssessment, error) {
	workloads := models.WorkloadAssessment{
		Pods:         make([]models.PodInfo, 0),
		Deployments:  make([]models.DeploymentInfo, 0),
		StatefulSets: make([]models.StatefulSetInfo, 0),
		DaemonSets:   make([]models.DaemonSetInfo, 0),
	}

	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return workloads, err
	}

	for _, ns := range namespaces.Items {
		// Collect Pods
		pods, err := c.client.CoreV1().Pods(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, pod := range pods.Items {
			containers := make([]models.ContainerInfo, 0)
			for _, container := range pod.Spec.Containers {
				securityContext := container.SecurityContext
				var containerSecInfo models.ContainerSecurityInfo

				if securityContext != nil {
					containerSecInfo = models.ContainerSecurityInfo{
						Capabilities:             getCapabilities(securityContext.Capabilities),
						RunAsUser:                securityContext.RunAsUser,
						RunAsNonRoot:             securityContext.RunAsNonRoot,
						ReadOnlyRoot:             securityContext.ReadOnlyRootFilesystem != nil && *securityContext.ReadOnlyRootFilesystem,
						Privileged:               securityContext.Privileged != nil && *securityContext.Privileged,
						AllowPrivilegeEscalation: securityContext.AllowPrivilegeEscalation,
					}
				} // Collect environment variables
				envVars := make([]string, 0)
				for _, env := range container.Env {
					envVars = append(envVars, env.Name)
				}

				containers = append(containers, models.ContainerInfo{
					Name:            container.Name,
					Image:           container.Image,
					SecurityContext: containerSecInfo,
					Resources: models.ResourceRequirements{
						Limits: models.ResourceList{
							CPU:    container.Resources.Limits.Cpu().String(),
							Memory: container.Resources.Limits.Memory().String(),
						},
						Requests: models.ResourceList{
							CPU:    container.Resources.Requests.Cpu().String(),
							Memory: container.Resources.Requests.Memory().String(),
						},
					},
					EnvVars: envVars,
				})
			}

			podSecurity := pod.Spec.SecurityContext
			var podSecInfo models.PodSecurityInfo
			if podSecurity != nil {
				podSecInfo = models.PodSecurityInfo{
					RunAsUser:   podSecurity.RunAsUser,
					RunAsGroup:  podSecurity.RunAsGroup,
					FSGroup:     podSecurity.FSGroup,
					HostNetwork: pod.Spec.HostNetwork,
					HostPID:     pod.Spec.HostPID,
					HostIPC:     pod.Spec.HostIPC}
			}

			workloads.Pods = append(workloads.Pods, models.PodInfo{
				Name:                         pod.Name,
				Namespace:                    pod.Namespace,
				ServiceAccount:               pod.Spec.ServiceAccountName,
				SecurityContext:              podSecInfo,
				Containers:                   containers,
				NodeName:                     pod.Spec.NodeName,
				CreatedAt:                    pod.CreationTimestamp.String(),
				Labels:                       pod.Labels,
				AutomountServiceAccountToken: pod.Spec.AutomountServiceAccountToken,
			})
		}

		// Collect Deployments
		deployments, err := c.client.AppsV1().Deployments(ns.Name).List(ctx, metav1.ListOptions{})
		if err == nil {
			for _, deploy := range deployments.Items {
				workloads.Deployments = append(workloads.Deployments, models.DeploymentInfo{
					Name:           deploy.Name,
					Namespace:      deploy.Namespace,
					Replicas:       *deploy.Spec.Replicas,
					UpdateStrategy: string(deploy.Spec.Strategy.Type),
					Labels:         deploy.Labels,
					CreatedAt:      deploy.CreationTimestamp.String(),
				})
			}
		}

		// Collect StatefulSets
		statefulSets, err := c.client.AppsV1().StatefulSets(ns.Name).List(ctx, metav1.ListOptions{})
		if err == nil {
			for _, sts := range statefulSets.Items {
				workloads.StatefulSets = append(workloads.StatefulSets, models.StatefulSetInfo{
					Name:           sts.Name,
					Namespace:      sts.Namespace,
					Replicas:       *sts.Spec.Replicas,
					UpdateStrategy: string(sts.Spec.UpdateStrategy.Type),
					Labels:         sts.Labels,
					CreatedAt:      sts.CreationTimestamp.String(),
				})
			}
		}

		// Collect DaemonSets
		daemonSets, err := c.client.AppsV1().DaemonSets(ns.Name).List(ctx, metav1.ListOptions{})
		if err == nil {
			for _, ds := range daemonSets.Items {
				workloads.DaemonSets = append(workloads.DaemonSets, models.DaemonSetInfo{
					Name:           ds.Name,
					Namespace:      ds.Namespace,
					UpdateStrategy: string(ds.Spec.UpdateStrategy.Type),
					Labels:         ds.Labels,
					CreatedAt:      ds.CreationTimestamp.String(),
				})
			}
		}
	}

	return workloads, nil
}

func getCapabilities(capabilities *corev1.Capabilities) []string {
	if capabilities == nil {
		return nil
	}

	var caps []string
	if capabilities.Add != nil {
		for _, cap := range capabilities.Add {
			caps = append(caps, fmt.Sprintf("+%s", cap))
		}
	}
	if capabilities.Drop != nil {
		for _, cap := range capabilities.Drop {
			caps = append(caps, fmt.Sprintf("-%s", cap))
		}
	}
	return caps
}
