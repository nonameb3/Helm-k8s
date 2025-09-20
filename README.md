## Prerequisites
```bash
# Check Kubernetes cluster
kubectl cluster-info
kubectl get nodes
```

## Multi-Environment Deployment

### Development Environment
```bash
# Test deployment (dry-run)
helm upgrade --install health-service-dev devops/helm/ --namespace health-service-dev --create-namespace --values devops/helm/environments/values.dev.yaml --dry-run --debug

# Deploy development environment
helm upgrade --install health-service-dev devops/helm/ --namespace health-service-dev --create-namespace --values devops/helm/environments/values.dev.yaml

# Check deployment
kubectl get pods -n health-service-dev
kubectl get services -n health-service-dev
kubectl get ingress -n health-service-dev
kubectl get hpa -n health-service-dev
```

### Staging Environment
```bash
# Deploy staging environment
helm upgrade --install health-service-staging devops/helm/ --namespace health-service-staging --create-namespace --values devops/helm/environments/values.nonprod.yaml

# Check deployment
kubectl get pods -n health-service-staging
kubectl get services -n health-service-staging
kubectl get ingress -n health-service-staging
kubectl get hpa -n health-service-staging
```

### Production Environment
```bash
# Deploy production environment
helm upgrade --install health-service-prod devops/helm/ --namespace health-service-prod --create-namespace --values devops/helm/environments/values.prod.yaml

# Check deployment
kubectl get pods -n health-service-prod
kubectl get services -n health-service-prod
kubectl get ingress -n health-service-prod
kubectl get hpa -n health-service-prod
```

## Setup External Access (NGINX Ingress)
```bash
# Install NGINX Ingress Controller (one-time setup)
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml

# Wait for ingress controller to be ready
kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s

# Add local DNS entries for all environments
echo "127.0.0.1 health-service-dev.local" | sudo tee -a /etc/hosts
echo "127.0.0.1 health-service-staging.local" | sudo tee -a /etc/hosts
echo "127.0.0.1 health-service-prod.local" | sudo tee -a /etc/hosts

# Test all environments
curl http://health-service-dev.local/health
curl http://health-service-staging.local/health
curl http://health-service-prod.local/health
```

## Useful Commands

### View Resources by Environment
```bash
# Development
kubectl get all -n health-service-dev
kubectl get ingress -n health-service-dev
kubectl logs -n health-service-dev -l service=health-service-dev

# Staging
kubectl get all -n health-service-staging
kubectl get ingress -n health-service-staging
kubectl logs -n health-service-staging -l service=health-service-staging

# Production
kubectl get all -n health-service-prod
kubectl get ingress -n health-service-prod
kubectl logs -n health-service-prod -l service=health-service-prod
```

### Remove Deployments
```bash
# Remove specific environment
helm uninstall health-service-dev --namespace health-service-dev
helm uninstall health-service-staging --namespace health-service-staging
helm uninstall health-service-prod --namespace health-service-prod

# Remove all environments at once
kubectl delete namespace health-service-dev health-service-staging health-service-prod
```

## Troubleshooting & Validation

### Helm Chart Validation
```bash
# Validate Helm chart (requires environment values)
helm lint devops/helm/ --values devops/helm/environments/values.dev.yaml

# Preview templates for specific environment
helm template health-service-dev devops/helm/ --values devops/helm/environments/values.dev.yaml
helm template health-service-staging devops/helm/ --values devops/helm/environments/values.nonprod.yaml
helm template health-service-prod devops/helm/ --values devops/helm/environments/values.prod.yaml
```

### Environment-Specific Troubleshooting
```bash
# Check pod details (replace ENVIRONMENT with dev/staging/prod)
kubectl describe pod -n health-service-ENVIRONMENT -l service=health-service-ENVIRONMENT

# Check service endpoints
kubectl get endpoints -n health-service-ENVIRONMENT

# Validate ingress configuration
kubectl describe ingress -n health-service-ENVIRONMENT

# Check pod events for issues
kubectl get events -n health-service-ENVIRONMENT --sort-by='.lastTimestamp'
```

### General Troubleshooting
```bash
# Check if ingress controller is running
kubectl get pods -n ingress-nginx

# Test internal cluster connectivity
kubectl run test-pod --rm -i --tty --image=busybox -- /bin/sh
# Inside pod, test each environment:
# wget -qO- http://health-service-dev-service.health-service-dev.svc.cluster.local/health
# wget -qO- http://health-service-staging-service.health-service-staging.svc.cluster.local/health
# wget -qO- http://health-service-prod-service.health-service-prod.svc.cluster.local/health

# View all resources across environments
kubectl get pods,services,ingress --all-namespaces | grep health-service
```

## Horizontal Pod Autoscaler (HPA)

### HPA Configuration by Environment
| Environment | Min Replicas | Max Replicas | CPU Target | Memory Target | Scale Up Policy | Scale Down Policy |
|-------------|--------------|--------------|------------|---------------|----------------|------------------|
| Development | 1 | 3 | 70% | 80% | Conservative (50%/1 pod per 60s) | Slow (10% per 60s) |
| Staging | 2 | 5 | 60% | 70% | Moderate (100%/2 pods per 30s) | Medium (20% per 60s) |
| Production | 3 | 10 | 50% | 60% | Aggressive (100%/3 pods per 15s) | Careful (25% per 60s) |

### HPA Monitoring Commands
```bash
# Check HPA status across all environments
kubectl get hpa --all-namespaces | grep health-service

# Watch HPA scaling in real-time
kubectl get hpa -n health-service-dev -w

# Detailed HPA information
kubectl describe hpa -n health-service-dev health-service-dev-hpa

# Check current resource usage
kubectl top pods --all-namespaces | grep health-service
```

### Load Testing HPA

#### Method 1: Built-in Load Endpoint (Recommended)
```bash
# Generate sustained CPU load to trigger scaling
curl http://health-service-dev.local/load/30000 &  # 30-second CPU load
curl http://health-service-dev.local/load/30000 &  # 30-second CPU load
curl http://health-service-dev.local/load/30000 &  # 30-second CPU load

# Test multiple concurrent requests to see load distribution
for i in {1..6}; do
  curl http://health-service-dev.local/load/20000 &
done

# Watch scaling in real-time
kubectl get pods -n health-service-dev -w
kubectl get hpa -n health-service-dev -w

# Monitor which pods handle requests
kubectl logs -n health-service-dev -l service=health-service-dev --tail=20 -f
```

#### Method 2: External Load Testing
```bash
# Install Apache Bench for load testing
# macOS: brew install httpd
# Ubuntu: sudo apt install apache2-utils

# Generate load to trigger scaling
ab -n 10000 -c 50 http://health-service-dev.local/health

# Monitor HPA decisions
kubectl get events -n health-service-dev --sort-by='.lastTimestamp' | grep HorizontalPodAutoscaler
```

### Troubleshooting HPA
```bash
# Check if metrics-server is running
kubectl get pods -n kube-system | grep metrics-server
kubectl top nodes

# Verify pod resource requests (required for HPA)
kubectl describe pod -n health-service-dev -l service=health-service-dev

# Check HPA events and scaling decisions
kubectl describe hpa -n health-service-dev

# View HPA controller logs
kubectl logs -n kube-system -l k8s-app=metrics-server
```

## Environment Configuration Overview
- **Development**: 1-3 replicas, minimal resources, conservative HPA
- **Staging**: 2-5 replicas, production-like resources, moderate HPA
- **Production**: 3-10 replicas, maximum resources, aggressive HPA