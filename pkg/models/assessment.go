package models

// ClusterInfo represents basic information about the Kubernetes cluster
// | Version | NodeCount | APIServer | Platform | Components | Nodes | Namespaces |
type ClusterInfo struct {
	Version    string
	NodeCount  int
	APIServer  string
	Platform   string
	Nodes      []NodeInfo
	Namespaces []NamespaceInfo
}

// NodeInfo represents detailed information about a node
// | Name | Version | Architecture | OS | ContainerRuntime | CPU | Memory | Ready | Labels |
type NodeInfo struct {
	Name             string
	Version          string
	Architecture     string
	OS               string
	ContainerRuntime string
	CPU              string
	Memory           string
	Ready            bool
	Labels           map[string]string
}

// NamespaceInfo represents detailed information about a namespace
// | Name | Status | CreatedAt | Labels |
type NamespaceInfo struct {
	Name      string
	Status    string
	CreatedAt string
	Labels    map[string]string
}

// RBACAssessment contains RBAC-related security information
// | ClusterRoles | ClusterRoleBindings | Roles | RoleBindings | ServiceAccounts |
type RBACAssessment struct {
	ClusterRoles        []RoleInfo
	ClusterRoleBindings []BindingInfo
	Roles               []RoleInfo
	RoleBindings        []BindingInfo
	ServiceAccounts     []ServiceAccountInfo
}

// ClusterRolesOnly returns all roles that are cluster roles (ClusterRole == true)
func (r *RBACAssessment) ClusterRolesOnly() []RoleInfo {
	roles := make([]RoleInfo, 0)
	for _, role := range r.ClusterRoles {
		if role.ClusterRole {
			roles = append(roles, role)
		}
	}
	for _, role := range r.Roles {
		if role.ClusterRole {
			roles = append(roles, role)
		}
	}
	return roles
}

// RoleInfo contains information about RBAC roles and cluster roles
// | Name | Namespace | ClusterRole | Rules | CreatedAt |
type RoleInfo struct {
	Name        string
	Namespace   string
	ClusterRole bool // true if this is a ClusterRole
	Rules       []PolicyRule
	CreatedAt   string
}

// BindingInfo contains information about RBAC bindings
// | Name | Namespace | RoleRef | Subjects | CreatedAt |
type BindingInfo struct {
	Name      string
	Namespace string
	RoleRef   string
	Subjects  []Subject
	CreatedAt string
}

// PolicyRule represents an RBAC policy rule
// | APIGroups | Resources | ResourceNames | Verbs |
type PolicyRule struct {
	APIGroups     []string
	Resources     []string
	ResourceNames []string
	Verbs         []string
}

// Subject represents a binding subject
// | Kind | Name | Namespace |
type Subject struct {
	Kind      string
	Name      string
	Namespace string
}

// DeploymentInfo contains information about deployments
// | Name | Namespace | Replicas | UpdateStrategy | Labels | CreatedAt |
type DeploymentInfo struct {
	Name           string
	Namespace      string
	Replicas       int32
	UpdateStrategy string
	Labels         map[string]string
	CreatedAt      string
}

// StatefulSetInfo contains information about stateful sets
// | Name | Namespace | Replicas | UpdateStrategy | Labels | CreatedAt |
type StatefulSetInfo struct {
	Name           string
	Namespace      string
	Replicas       int32
	UpdateStrategy string
	Labels         map[string]string
	CreatedAt      string
}

// DaemonSetInfo contains information about daemon sets
// | Name | Namespace | UpdateStrategy | Labels | CreatedAt |
type DaemonSetInfo struct {
	Name           string
	Namespace      string
	UpdateStrategy string
	Labels         map[string]string
	CreatedAt      string
}

// WorkloadAssessment contains information about workloads
// | Pods | Deployments | StatefulSets | DaemonSets |
type WorkloadAssessment struct {
	Pods         []PodInfo
	Deployments  []DeploymentInfo
	StatefulSets []StatefulSetInfo
	DaemonSets   []DaemonSetInfo
}

// CommonInfo for all resources
// | Name | Namespace | Labels | CreatedAt |
type CommonInfo struct {
	Name      string
	Namespace string
	Labels    map[string]string
	CreatedAt string
}

// PodInfo contains pod-level information including security context
// | Name | Namespace | NodeName | ServiceAccount | Labels | CreatedAt | SecurityContext | Containers | AutomountServiceAccountToken |
type PodInfo struct {
	Name                         string
	Namespace                    string
	NodeName                     string
	ServiceAccount               string
	Labels                       map[string]string
	CreatedAt                    string
	SecurityContext              PodSecurityInfo
	Containers                   []ContainerInfo
	AutomountServiceAccountToken *bool
}

// ContainerInfo contains security-relevant information about containers
// | Name | Image | SecurityContext | Resources | EnvVars |
type ContainerInfo struct {
	Name            string
	Image           string
	SecurityContext ContainerSecurityInfo
	Resources       ResourceRequirements
	EnvVars         []string
}

// PodSecurityInfo contains pod-level security context information
// | RunAsUser | RunAsGroup | FSGroup | HostNetwork | HostPID | HostIPC |
type PodSecurityInfo struct {
	RunAsUser   *int64
	RunAsGroup  *int64
	FSGroup     *int64
	HostNetwork bool
	HostPID     bool
	HostIPC     bool
}

// ContainerSecurityInfo contains container-level security context
// | Capabilities | RunAsUser | RunAsNonRoot | ReadOnlyRoot | Privileged | AllowPrivilegeEscalation |
type ContainerSecurityInfo struct {
	Capabilities             []string
	RunAsUser                *int64
	RunAsNonRoot             *bool
	ReadOnlyRoot             bool
	Privileged               bool
	AllowPrivilegeEscalation *bool
}

// ResourceRequirements contains container resource constraints
// | Limits | Requests |
type ResourceRequirements struct {
	Limits   ResourceList
	Requests ResourceList
}

// ResourceList contains resource quantities
// | CPU | Memory |
type ResourceList struct {
	CPU    string
	Memory string
}

// NetworkAssessment contains networking-related security information
// | Services | NetworkPolicies | Ingresses |
type NetworkAssessment struct {
	Services        []ServiceInfo
	NetworkPolicies []NetworkPolicyInfo
	Ingresses       []IngressInfo
}

// ServiceInfo represents a Kubernetes Service
// | Name | Namespace | Labels | CreatedAt | Type | ClusterIP | ExternalIPs | Ports | Size |
type ServiceInfo struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	CreatedAt   string
	Type        string
	ClusterIP   string
	ExternalIPs []string
	Ports       []ServicePort
	Size        int64
}

// ServicePort represents a service port configuration
// | Port | TargetPort | Protocol |
type ServicePort struct {
	Port       int32
	TargetPort int32
	Protocol   string
}

// NetworkPolicyInfo represents a Kubernetes NetworkPolicy
// | Name | Namespace | Labels | CreatedAt | PodSelector | PolicyTypes |
type NetworkPolicyInfo struct {
	Name        string
	Namespace   string
	Labels      map[string]string
	CreatedAt   string
	PodSelector string
	PolicyTypes []string
}

// IngressInfo represents a Kubernetes Ingress
// | Name | Namespace | Labels | CreatedAt | Rules | TLS |
type IngressInfo struct {
	Name      string
	Namespace string
	Labels    map[string]string
	CreatedAt string
	Rules     []IngressRule
	TLS       []string
}

// IngressRule represents a rule in an Ingress resource
// | Host | Paths |
type IngressRule struct {
	Host  string
	Paths []IngressPath
}

// IngressPath represents a path in an Ingress rule
// | Path | ServiceName | ServicePort |
type IngressPath struct {
	Path        string
	ServiceName string
	ServicePort int32
}

// ServiceAccountInfo represents a Kubernetes ServiceAccount
// | Name | Namespace | Labels | CreatedAt | Secrets | ImagePullSecrets |
type ServiceAccountInfo struct {
	Name             string
	Namespace        string
	Labels           map[string]string
	CreatedAt        string
	Secrets          []string
	ImagePullSecrets []string
}

// SecretInfo represents a Kubernetes Secret
// | Name | Namespace | Labels | CreatedAt | Type |
type SecretInfo struct {
	CommonInfo
	Type string
}

// SecretAssessment contains information about Kubernetes Secrets
// | Secrets |
type SecretAssessment struct {
	Secrets []SecretInfo
}

// AssessmentData represents all collected assessment data
// | ClusterInfo | RBAC | Workloads | Network | Secrets |
type AssessmentData struct {
	ClusterInfo ClusterInfo
	RBAC        RBACAssessment
	Workloads   WorkloadAssessment
	Network     NetworkAssessment
	Secrets     SecretAssessment
}
