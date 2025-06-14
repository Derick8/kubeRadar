# kubeRadar - Kubernetes Reconnaissance Tool
![image](https://github.com/user-attachments/assets/a4b12b90-713c-45f7-91c4-2625a05ad8b9)


kubeRadar is a tool that performs reconnaissance of Kubernetes clusters and generates a comprehensive Excel report. The report provides detailed visibility into cluster configuration, RBAC, workloads, networking, and secrets, helping you inventory and review your cluster setup.

## Features

kubeRadar provides comprehensive visibility into your Kubernetes cluster by collecting metadata, RBAC configuration, workloads, networking, and secrets. It generates a multi-sheet Excel report that helps you inventory, review, and audit your cluster's setup, including nodes, namespaces, workloads, network policies, RBAC roles and bindings, service accounts, and secrets. This enables security reviews and operational insight in a single, professional report.

## Prerequisites

- Go 1.21 or higher
- Access to a Kubernetes cluster
- kubeconfig file with proper credentials

## Installation

```bash
# Clone the repository
git clone https://github.com/derick8/kubeRadar.git
cd kubeRadar
# Build the project
make build
```

## Usage

Run the tool from the root folder with your kubeconfig:

```bash
# Run directly with Go
 go run ./main.go --kubeconfig <path-to-kubeconfig> --output <output-file.xlsx>

# Or build and run the executable
 go build -o kubeRadar.exe main.go
 ./kubeRadar.exe --kubeconfig <path-to-kubeconfig> --output <output-file.xlsx>
```

- `--kubeconfig` (optional): Path to your kubeconfig file. Defaults to `~/.kube/config`.
- `--output` (optional): Output Excel file path. Defaults to `k8s_assessment.xlsx`.

During execution, the tool prints status messages to the console (stderr) to indicate progress (e.g., collecting data, generating report, writing Excel file).


## Excel Report Structure

The generated Excel report contains the following worksheets, each with detailed columns:

- **Nodes**: Name, Version, Architecture, OS, Container Runtime, CPU, Memory, Ready, Labels
- **Namespaces**: Name, Status, Created At, Labels
- **Pods**: Name, Namespace, Node, Service Account, Privileged, Host Network, Host PID, Host IPC, Run As User, Run As Non Root, Auto Mount SA Token, Container Names, Container Images, Capabilities, Resources, Sysctls, Environment Variables, Created At, Labels
- **Deployments**: Name, Namespace, Replicas, Update Strategy, Created At, Labels
- **StatefulSets**: Name, Namespace, Replicas, Update Strategy, Created At, Labels
- **DaemonSets**: Name, Namespace, Update Strategy, Created At, Labels
- **Services**: Name, Namespace, Type, Cluster IP, External IPs, Ports
- **Network Policies**: Name, Namespace, Pod Selector, Policy Types, Created At, Labels
- **Ingresses**: Name, Namespace, Rules, TLS, Created At, Labels
- **RBAC Roles**: Name, Namespace, Created At, Rules
- **Role Bindings**: Name, Namespace, Role Ref, Subjects, Created At
- **Cluster Roles**: Name, Created At, Rules
- **Cluster Role Bindings**: Name, Role Ref, Subjects, Created At
- **Service Accounts**: Name, Namespace, Secrets, Image Pull Secrets, Created At, Labels
- **Secrets**: Name, Namespace, Type, Created At

## Logo Symbolism

- **Hexagon**: A nod to the Kubernetes logo, representing a container cluster.
- **Ship Wheel**: Derived from the original Kubernetes symbol, symbolizing control over distributed systems.
- **Radar Needle + Dots**: Placed inside the wheel, these elements symbolise scanning or reconnaissance, representing the tool's core functionality: scanning Kubernetes environments.
## License

This project is licensed under the MIT License - see the LICENSE file for details.
