package collector

import (
	"context"
	"kubeRadar/pkg/models"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Collector struct {
	client *kubernetes.Clientset
	config *rest.Config
}

func NewCollector(kubeconfigPath string) (*Collector, error) {
	if kubeconfigPath == "" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfigPath = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Collector{
		client: clientset,
		config: config,
	}, nil
}

func (c *Collector) CollectAll() (*models.AssessmentData, error) {
	ctx := context.Background()

	clusterInfo, err := c.collectClusterInfo(ctx)
	if err != nil {
		return nil, err
	}

	rbac, err := c.collectRBACInfo(ctx)
	if err != nil {
		return nil, err
	}

	workloads, err := c.collectWorkloadInfo(ctx)
	if err != nil {
		return nil, err
	}

	network, err := c.collectNetworkInfo(ctx)
	if err != nil {
		return nil, err
	}

	secrets, err := c.collectSecretInfo(ctx)
	if err != nil {
		return nil, err
	}

	return &models.AssessmentData{
		ClusterInfo: clusterInfo,
		RBAC:        rbac,
		Workloads:   workloads,
		Network:     network,
		Secrets:     secrets,
	}, nil
}

func (c *Collector) collectClusterInfo(ctx context.Context) (models.ClusterInfo, error) {
	version, err := c.client.Discovery().ServerVersion()
	if err != nil {
		return models.ClusterInfo{}, err
	}

	nodes, err := c.client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return models.ClusterInfo{}, err
	}

	// Collect namespace information
	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return models.ClusterInfo{}, err
	}

	// Collect node details
	nodeDetails := make([]models.NodeInfo, 0)
	for _, node := range nodes.Items {
		nodeInfo := models.NodeInfo{
			Name:             node.Name,
			Version:          node.Status.NodeInfo.KubeletVersion,
			Architecture:     node.Status.NodeInfo.Architecture,
			OS:               node.Status.NodeInfo.OperatingSystem,
			ContainerRuntime: node.Status.NodeInfo.ContainerRuntimeVersion,
			CPU:              node.Status.Capacity.Cpu().String(),
			Memory:           node.Status.Capacity.Memory().String(),
			Ready:            false,
			Labels:           node.Labels,
		}

		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				nodeInfo.Ready = condition.Status == "True"
				break
			}
		}
		nodeDetails = append(nodeDetails, nodeInfo)
	}

	// Collect namespace details
	namespaceDetails := make([]models.NamespaceInfo, 0)
	for _, ns := range namespaces.Items {
		nsInfo := models.NamespaceInfo{
			Name:      ns.Name,
			Status:    string(ns.Status.Phase),
			CreatedAt: ns.CreationTimestamp.String(),
			Labels:    ns.Labels,
		}
		namespaceDetails = append(namespaceDetails, nsInfo)
	}

	// Get platform info from nodes
	platform := ""
	if len(nodes.Items) > 0 {
		platform = nodes.Items[0].Status.NodeInfo.OperatingSystem
	}

	return models.ClusterInfo{
		Version:    version.String(),
		NodeCount:  len(nodes.Items),
		Platform:   platform,
		Nodes:      nodeDetails,
		Namespaces: namespaceDetails,
		APIServer:  c.config.Host,
	}, nil
}
