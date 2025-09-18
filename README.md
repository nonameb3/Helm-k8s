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