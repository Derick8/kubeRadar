package collector

import (
	"context"
	"kubeRadar/pkg/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Collector) collectNetworkInfo(ctx context.Context) (models.NetworkAssessment, error) {
	network := models.NetworkAssessment{}

	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return network, err
	}

	for _, ns := range namespaces.Items {
		// Collect Services
		services, err := c.client.CoreV1().Services(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, svc := range services.Items {
			ports := make([]models.ServicePort, 0)
			for _, port := range svc.Spec.Ports {
				ports = append(ports, models.ServicePort{
					Port:       port.Port,
					Protocol:   string(port.Protocol),
					TargetPort: port.TargetPort.IntVal,
				})
			}

			network.Services = append(network.Services, models.ServiceInfo{
				Name:        svc.Name,
				Namespace:   svc.Namespace,
				Labels:      svc.Labels,
				CreatedAt:   svc.CreationTimestamp.String(),
				Type:        string(svc.Spec.Type),
				ClusterIP:   svc.Spec.ClusterIP,
				ExternalIPs: svc.Spec.ExternalIPs,
				Ports:       ports,
				Size:        0,
			})
		}

		// Collect NetworkPolicies
		netpols, err := c.client.NetworkingV1().NetworkPolicies(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, netpol := range netpols.Items {
			policyTypes := make([]string, 0)
			for _, ptype := range netpol.Spec.PolicyTypes {
				policyTypes = append(policyTypes, string(ptype))
			}

			network.NetworkPolicies = append(network.NetworkPolicies, models.NetworkPolicyInfo{
				Name:        netpol.Name,
				Namespace:   netpol.Namespace,
				Labels:      netpol.Labels,
				CreatedAt:   netpol.CreationTimestamp.String(),
				PodSelector: netpol.Spec.PodSelector.String(),
				PolicyTypes: policyTypes,
			})
		}

		// Collect Ingresses
		ingresses, err := c.client.NetworkingV1().Ingresses(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, ing := range ingresses.Items {
			ingressRules := make([]models.IngressRule, 0)
			for _, rule := range ing.Spec.Rules {
				if rule.Host != "" {
					paths := make([]models.IngressPath, 0)
					if rule.HTTP != nil {
						for _, path := range rule.HTTP.Paths {
							paths = append(paths, models.IngressPath{
								Path:        path.Path,
								ServiceName: path.Backend.Service.Name,
								ServicePort: path.Backend.Service.Port.Number,
							})
						}
					}
					ingressRules = append(ingressRules, models.IngressRule{
						Host:  rule.Host,
						Paths: paths,
					})
				}
			}

			tlsHosts := make([]string, 0)
			for _, tls := range ing.Spec.TLS {
				tlsHosts = append(tlsHosts, tls.Hosts...)
			}

			network.Ingresses = append(network.Ingresses, models.IngressInfo{
				Name:      ing.Name,
				Namespace: ing.Namespace,
				Labels:    ing.Labels,
				CreatedAt: ing.CreationTimestamp.String(),
				Rules:     ingressRules,
				TLS:       tlsHosts,
			})
		}
	}

	return network, nil
}
