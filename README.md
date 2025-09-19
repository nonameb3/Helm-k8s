## Kube
```
kubectl cluster-info
kubectl get nodes
```
# Deploy
```
helm install health-service devops/helm/ --namespace health-service --create-namespace
kubectl get pods
kubectl get services
kubectl get namespace
```
# NGINX Ingress
```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s
helm upgrade health-service devops/helm/ --namespace health-service
echo "127.0.0.1 health-service.local" | sudo tee -a /etc/hosts
curl http://health-service.local/health
```