### Prepare self-signed certification
```shell
openssl x509 -req -extfile ./cert/v3.ext -days 7300 -in ./cert/server.csr -CA ./cert/ca.crt -CAkey ./cert/ca.key -CAcreateserial -out ./cert/server.crt
```

### Create TLS secret
```shell
kubectl create secret tls gitops.com --key ./cert/server.key --cert ./cert/server.crt
```

### Deploy using Helm
```shell

helm repo add gitlab https://charts.gitlab.io/
helm repo update

# certmanger+letsencyipt
helm upgrade --install gitlab gitlab/gitlab \
  --timeout 600s \
  --set global.edition=ce \
  --set global.hosts.domain=gitops.com \
  --set global.hosts.externalIP=192.68.3.101 \
  --set certmanager-issuer.email=me@gitops.com \
  --set postgresql.image.tag=13.6.0

# 自签泛域名证书
helm upgrade --install gitlab gitlab/gitlab \
  --set global.edition=ce \
  --set global.hosts.domain=gitops.com \
  --set global.hosts.externalIP=192.68.3.101 \
  --set certmanager.install=false \
  --set global.ingress.configureCertmanager=false \
  --set global.ingress.tls.secretName=gitops.com \
  --set postgresql.image.tag=13.6.0
  ```

### Monitoring the Deployment

```shell
helm status gitlab
```

### Initial login

```shell 
kubectl get secret gitlab-gitlab-initial-root-password -ojsonpath='{.data.password}' | base64 --decode ; echo
bFy0AsRDOFbI01NyWKZfyW6JIQiTEGvxRTcF4nNPjdbQb2Wa7mj42cBcfslNHSn4
```