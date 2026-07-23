# Website Kubernetes Operator

A production-ready Kubernetes Operator built with **Go** and **Kubebuilder** that automates website deployment and lifecycle management using a custom Kubernetes resource.

The operator follows the Kubernetes Operator pattern and continuously reconciles the desired state of a website by creating and managing the required Kubernetes resources automatically.

---

## Features

- Custom Resource Definition (CRD) for Website
- Automatic Deployment creation
- Automatic Service creation
- Automatic Ingress creation
- Automatic TLS certificate provisioning
- cert-manager integration
- Prometheus ServiceMonitor support
- Horizontal Pod Autoscaler (HPA)
- NetworkPolicy generation
- Status updates on Custom Resources
- Kubernetes Events
- Finalizers for cleanup
- Leader Election support
- RBAC configuration
- Production-ready reconciliation loop

---

## Architecture

```
                    Website Custom Resource
                              │
                              ▼
                   Website Controller (Operator)
                              │
      ┌─────────────┬──────────┼──────────┬─────────────┐
      │             │          │          │             │
      ▼             ▼          ▼          ▼             ▼
 Deployment      Service    Ingress   Certificate   NetworkPolicy
      │             │          │          │             │
      └─────────────┴──────────┴──────────┴─────────────┘
                              │
                              ▼
                       Running Website
```

---

## Custom Resource Example

```yaml
apiVersion: apps.homelab.io/v1alpha1
kind: Website
metadata:
  name: demo-site

spec:
  image: nginx:latest

  replicas: 2

  port: 80

  hostname: demo.example.com

  tls:
    enabled: true

  monitoring:
    enabled: true
```

---

## Resources Managed

The operator automatically creates and manages:

- Deployment
- Service
- Ingress
- TLS Certificate
- Secret
- HorizontalPodAutoscaler
- ServiceMonitor
- NetworkPolicy

---

## Tech Stack

- Go
- Kubebuilder
- controller-runtime
- Kubernetes API
- Custom Resource Definitions (CRDs)
- cert-manager
- Prometheus Operator
- Docker
- Helm
- Kind / Minikube / Kubernetes

---

## Project Structure

```
.
├── api/
├── cmd/
├── config/
├── controllers/
├── hack/
├── internal/
├── pkg/
├── test/
├── Dockerfile
├── Makefile
└── README.md
```

---

## Getting Started

### Prerequisites

- Go 1.24+
- Docker
- Kubectl
- Kind / Minikube / Kubernetes Cluster
- Kubebuilder
- cert-manager
- Prometheus Operator

### Clone Repository

```bash
git clone https://github.com/<your-username>/website-operator.git

cd website-operator
```

### Install CRDs

```bash
make install
```

### Run Locally

```bash
make run
```

### Build Image

```bash
make docker-build IMG=<your-image>
```

### Deploy Operator

```bash
make deploy IMG=<your-image>
```

---

## Operator Workflow

1. User creates a Website Custom Resource.
2. The controller detects the new resource.
3. Deployment is created.
4. Service is created.
5. Ingress is configured.
6. TLS certificate is requested.
7. ServiceMonitor is created.
8. HPA is configured.
9. NetworkPolicy is applied.
10. Status is continuously reconciled.

---

## Roadmap

### Phase 1

- Website Operator
- Deployment automation
- Service automation
- Ingress automation
- TLS support

### Phase 2

- ConfigMap management
- Secret management
- Vault integration
- Health monitoring

### Phase 3

- AI-powered security scanning
- Kubernetes security recommendations
- Compliance reports
- Risk scoring
- Web dashboard

### Phase 4

- Multi-agent architecture
- OIDC authentication
- Keycloak integration
- React frontend
- PostgreSQL backend
- AI-powered Kubernetes assistant

---

## Learning Objectives

This project demonstrates:

- Kubernetes Operators
- Controller Runtime
- Reconciliation Loop
- Custom Resources
- Status Management
- Owner References
- Finalizers
- Watches
- Events
- Leader Election
- Admission Webhooks
- Production-grade Kubernetes development

---

## License

MIT