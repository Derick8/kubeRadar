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

Each worksheet in the Excel report contains the following columns:

- **Nodes**: Name, Version, Architecture, OS, Container Runtime, CPU, Memory, Ready, Labels
- **Namespaces**: Name, Status, Created At, Labels
- **Pods**: Name, Namespace, Node, Service Account, Host Network, Host PID, Host IPC, Created At
- **Deployments**: Name, Namespace, Replicas, Update Strategy, Created At, Labels
- **StatefulSets**: Name, Namespace, Replicas, Update Strategy, Created At, Labels
- **DaemonSets**: Name, Namespace, Update Strategy, Created At, Labels
- **Services**: Name, Namespace, Type, Cluster IP, External IPs, Ports
- **Network Policies**: Name, Namespace, Pod Selector, Policy Types, Created At, Labels
- **Ingresses**: Name, Namespace, Rules, TLS, Created At, Labels
- **RBAC Roles**: Name, Namespace, Created At, Rules
- **Role Bindings**: Name, Namespace, Role Ref, Subjects, Created At
- **Service Accounts**: Name, Namespace, Secrets, Image Pull Secrets, Created At, Labels
- **Secrets**: Name, Namespace, Type, Created At

## Project Structure

- `pkg/collector/` - Kubernetes resource collectors
- `pkg/excel/` - Excel report generation logic
- `pkg/models/` - Data models
- `internal/` - Internal utilities
- `build/` - Build artifacts and generated reports

## Development

- Build the project: `make build`
- Clean build artifacts: `make clean`
- Install dependencies: `make deps`

## Notes

- Temporary/lock files like `~$k8s_assessment.xlsx` can be ignored or deleted.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
