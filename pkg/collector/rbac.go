package collector

import (
	"context"
	"kubeRadar/pkg/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Collector) collectRBACInfo(ctx context.Context) (models.RBACAssessment, error) {
	rbac := models.RBACAssessment{}

	// Collect ClusterRoles
	clusterRoles, err := c.client.RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return rbac, err
	}

	for _, cr := range clusterRoles.Items {
		rules := make([]models.PolicyRule, 0)
		for _, rule := range cr.Rules {
			rules = append(rules, models.PolicyRule{
				APIGroups:     rule.APIGroups,
				Resources:     rule.Resources,
				ResourceNames: rule.ResourceNames,
				Verbs:         rule.Verbs,
			})
		}

		rbac.ClusterRoles = append(rbac.ClusterRoles, models.RoleInfo{
			Name:        cr.Name,
			Namespace:   "", // ClusterRoles are cluster-scoped
			ClusterRole: true,
			Rules:       rules,
			CreatedAt:   cr.CreationTimestamp.String(),
		})
	}

	// Collect ClusterRoleBindings
	clusterRoleBindings, err := c.client.RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return rbac, err
	}

	for _, crb := range clusterRoleBindings.Items {
		subjects := make([]models.Subject, 0)
		for _, subject := range crb.Subjects {
			subjects = append(subjects, models.Subject{
				Kind:      subject.Kind,
				Name:      subject.Name,
				Namespace: subject.Namespace,
			})
		}

		rbac.ClusterRoleBindings = append(rbac.ClusterRoleBindings, models.BindingInfo{
			Name:      crb.Name,
			Namespace: "", // ClusterRoleBindings are cluster-scoped
			RoleRef:   crb.RoleRef.Name,
			Subjects:  subjects,
			CreatedAt: crb.CreationTimestamp.String(),
		})
	}

	// Collect Roles and RoleBindings from all namespaces
	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return rbac, err
	}

	for _, ns := range namespaces.Items {
		// Collect Roles
		roles, err := c.client.RbacV1().Roles(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, role := range roles.Items {
			rules := make([]models.PolicyRule, 0)
			for _, rule := range role.Rules {
				rules = append(rules, models.PolicyRule{
					APIGroups:     rule.APIGroups,
					Resources:     rule.Resources,
					ResourceNames: rule.ResourceNames,
					Verbs:         rule.Verbs,
				})
			}

			rbac.Roles = append(rbac.Roles, models.RoleInfo{
				Name:        role.Name,
				Namespace:   role.Namespace,
				ClusterRole: false,
				Rules:       rules,
				CreatedAt:   role.CreationTimestamp.String(),
			})
		}

		// Collect RoleBindings
		roleBindings, err := c.client.RbacV1().RoleBindings(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, rb := range roleBindings.Items {
			subjects := make([]models.Subject, 0)
			for _, subject := range rb.Subjects {
				subjects = append(subjects, models.Subject{
					Kind:      subject.Kind,
					Name:      subject.Name,
					Namespace: subject.Namespace,
				})
			}

			rbac.RoleBindings = append(rbac.RoleBindings, models.BindingInfo{
				Name:      rb.Name,
				Namespace: rb.Namespace,
				RoleRef:   rb.RoleRef.Name,
				Subjects:  subjects,
				CreatedAt: rb.CreationTimestamp.String(),
			})
		}

		// Collect ServiceAccounts from all namespaces
		serviceAccounts, err := c.client.CoreV1().ServiceAccounts(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}
		for _, sa := range serviceAccounts.Items {
			secrets := make([]string, 0)
			for _, s := range sa.Secrets {
				secrets = append(secrets, s.Name)
			}
			imagePullSecrets := make([]string, 0)
			for _, ips := range sa.ImagePullSecrets {
				imagePullSecrets = append(imagePullSecrets, ips.Name)
			}
			rbac.ServiceAccounts = append(rbac.ServiceAccounts, models.ServiceAccountInfo{
				Name:             sa.Name,
				Namespace:        sa.Namespace,
				Labels:           sa.Labels,
				CreatedAt:        sa.CreationTimestamp.String(),
				Secrets:          secrets,
				ImagePullSecrets: imagePullSecrets,
			})
		}
	}
	return rbac, nil
}
