### 1. Add the Helm repository

```shell
helm repo add jetstack https://charts.jetstack.io

```

### 2. Update your local Helm chart repository cache

```shell
helm repo update
```

### 3. Install cert-manager with CRDs

```shell
# Example: disabling prometheus using a Helm parameter
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.11.0 \
  --set prometheus.enabled=false \
  --set installCRDs=true


helm template \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.11.0 \
  --set prometheus.enabled=false \   # Example: disabling prometheus using a Helm parameter
  --set installCRDs=true \           # Uncomment to also template CRDs
  > cert-manager.custom.yaml
```

### 4. Uninstalling

```shell
# Ensure that all cert-manager resources
kubectl get Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges --all-namespaces

# Uninstalling with Helm
helm --namespace cert-manager delete cert-manager

# Delete the cert-manager namespace
kubectl delete namespace cert-manager

# Delete the cert-manager CustomResourceDefinitions
kubectl delete -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.crds.yaml

# Namespace Stuck in Terminating State
kubectl delete apiservice v1beta1.webhook.cert-manager.io
```
