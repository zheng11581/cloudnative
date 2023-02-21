## Install with Helm

```shell
helm repo add harbor https://helm.goharbor.io

helm repo update

kubectl apply -f harbor-site-tls.yaml
kubectl apply -f notary-site-tls.yaml

helm install harbor harbor/harbor -f values.yaml --namespace harbor --create-namespace
```