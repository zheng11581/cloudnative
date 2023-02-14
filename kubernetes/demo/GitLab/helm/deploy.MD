### Deploy using Helm
```shell
helm repo add gitlab https://charts.gitlab.io/
helm repo update
helm upgrade --install gitlab gitlab/gitlab \
  --timeout 600s \
  --set global.hosts.domain=gitops.com \
  --set global.hosts.externalIP=10.10.10.10 \
  --set certmanager-issuer.email=me@gitops.com \
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