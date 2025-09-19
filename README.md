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
```

### Staging Environment
```bash
# Deploy staging environment
helm upgrade --install health-service-staging devops/helm/ --namespace health-service-staging --create-namespace --values devops/helm/environments/values.nonprod.yaml

# Check deployment
kubectl get pods -n health-service-staging
kubectl get services -n health-service-staging
kubectl get ingress -n health-service-staging
```

### Production Environment
```bash
# Deploy production environment
helm upgrade --install health-service-prod devops/helm/ --namespace health-service-prod --create-namespace --values devops/helm/environments/values.prod.yaml

# Check deployment
kubectl get pods -n health-service-prod
kubectl get services -n health-service-prod
kubectl get ingress -n health-service-prod
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

## Environment Configuration Overview
- **Development**: 1 replica, minimal resources, local images
- **Staging**: 2 replicas, production-like resources
- **Production**: 3 replicas, maximum resources and CPU limits