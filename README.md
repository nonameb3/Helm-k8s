## Prerequisites
```bash
# Check Kubernetes cluster
kubectl cluster-info
kubectl get nodes
```

## Deploy Application
```bash
# Test deployment (dry-run) - shows what will be created/changed
helm upgrade --install health-service devops/helm/ --namespace health-service --create-namespace --dry-run --debug

# Preview generated templates
helm template health-service devops/helm/

# Deploy for real (after dry-run looks good)
helm upgrade --install health-service devops/helm/ --namespace health-service --create-namespace

# Check deployment
kubectl get pods -n health-service
kubectl get services -n health-service
kubectl get serviceaccount -n health-service
```

## Setup External Access (NGINX Ingress)
```bash
# Install NGINX Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml

# Wait for ingress controller to be ready
kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s

# Add local DNS entry
echo "127.0.0.1 health-service.local" | sudo tee -a /etc/hosts

# Test the application
curl http://health-service.local/health
```

## Useful Commands
```bash
# View all resources in namespace
kubectl get all -n health-service

# Check ingress status
kubectl get ingress -n health-service

# View pod logs
kubectl logs -n health-service -l service=health-service

# Remove deployment
helm uninstall health-service --namespace health-service
```

## Troubleshooting & Validation
```bash
# Validate Helm chart before deployment
helm lint devops/helm/

# Check what Helm will deploy (safe preview)
helm template health-service devops/helm/

# Dry-run to see changes without applying
helm upgrade --install health-service devops/helm/ --namespace health-service --dry-run

# Check pod details and events
kubectl describe pod -n health-service -l service=health-service

# Check service endpoints
kubectl get endpoints -n health-service

# Validate ingress configuration
kubectl describe ingress -n health-service

# Check if ingress controller is running
kubectl get pods -n ingress-nginx

# Test connectivity inside cluster
kubectl run test-pod --rm -i --tty --image=busybox -- /bin/sh
# Then inside pod: wget -qO- http://health-service-service.health-service.svc.cluster.local/health
```