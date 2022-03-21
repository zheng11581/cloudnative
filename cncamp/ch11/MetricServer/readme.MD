### 1. Install 
https://artifacthub.io/packages/helm/metrics-server/metrics-server

```shell
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
helm upgrade --install metrics-server metrics-server/metrics-server
```