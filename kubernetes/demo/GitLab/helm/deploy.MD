## Prepare
### Prepare metallb
```shell

```

### Prepare self-signed certification
```shell
openssl x509 -req -extfile ./cert/v3.ext -days 7300 -in ./cert/server.csr -CA ./cert/ca.crt -CAkey ./cert/ca.key -CAcreateserial -out ./cert/server.crt
```
## Deploy
### Create TLS secret
```shell
kubectl -n gitlab-system create secret tls gitops.com --key ./cert/server.key --cert ./cert/server.crt
```

### Deploy using Helm
```shell

helm repo add gitlab https://charts.gitlab.io/
helm repo update

# certmanger+letsencyipt
helm upgrade --install gitlab gitlab/gitlab --create-namespace --namespace=gitlab-system \
  --timeout 600s \
  --set global.edition=ce \
  --set global.hosts.domain=gitops.com \
  --set global.hosts.externalIP=192.68.3.101 \
  --set certmanager-issuer.email=me@gitops.com \
  --set postgresql.image.tag=13.6.0 \
  --set gitlab-runner.runners.privileged=true \ 
  --set nginx-ingress.defaultBackend.service.loadBalancerIP=192.68.3.101

# 自签泛域名证书
helm upgrade --install gitlab gitlab/gitlab --create-namespace --namespace=gitlab-system \
  --set global.edition=ce \
  --set global.hosts.domain=gitops.com \
  --set global.hosts.externalIP=192.68.3.101 \
  --set certmanager.install=false \
  --set global.ingress.configureCertmanager=false \
  --set global.ingress.tls.secretName=gitops.com \
  --set postgresql.image.tag=13.6.0 \
  --set gitlab-runner.runners.privileged=true \
  --set nginx-ingress.defaultBackend.service.loadBalancerIP=192.68.3.101
  
  ```

### Monitoring the Deployment

```shell
# Check if Pods are Ready
kubectl get pod -n gitlab-system

# Check if PVC are Bounded
kubectl get pvc -n gitlab-system

# Check nginx-ingress service LB
kubectl get svc -n gitlab-system -l app=nginx-ingress
```

### Initial login

```shell 
kubectl get secret gitlab-gitlab-initial-root-password -ojsonpath='{.data.password}' | base64 --decode ; echo
bFy0AsRDOFbI01NyWKZfyW6JIQiTEGvxRTcF4nNPjdbQb2Wa7mj42cBcfslNHSn4

# Update CA Root cert on client（self-signed）
cp ca.crt /usr/local/share/ca-certificates/
update-ca-certificates 

# 访问 https://gitlab.gitops.com
```
