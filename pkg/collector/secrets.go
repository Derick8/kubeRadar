package collector

import (
	"context"
	"kubeRadar/pkg/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Collector) collectSecretInfo(ctx context.Context) (models.SecretAssessment, error) {
	secretAssessment := models.SecretAssessment{}

	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return secretAssessment, err
	}

	for _, ns := range namespaces.Items {
		secrets, err := c.client.CoreV1().Secrets(ns.Name).List(ctx, metav1.ListOptions{})
		if err != nil {
			continue
		}

		for _, secret := range secrets.Items {
			secretAssessment.Secrets = append(secretAssessment.Secrets, models.SecretInfo{
				CommonInfo: models.CommonInfo{
					Name:      secret.Name,
					Namespace: secret.Namespace,
					Labels:    secret.Labels,
					CreatedAt: secret.CreationTimestamp.String(),
				},
				Type: string(secret.Type),
			})
		}
	}

	return secretAssessment, nil
}
