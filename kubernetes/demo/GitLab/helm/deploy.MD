## Prepare
### Prepare metallb
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

# 自签泛域名证书
helm upgrade --install gitlab gitlab/gitlab --create-namespace --namespace=gitlab-system \
  --set global.edition=ce \
  --set global.hosts.domain=gitops.com \
  --set certmanager.install=false \
  --set global.ingress.class=gitlab-nginx \
  --set global.ingress.configureCertmanager=false \
  --set global.ingress.tls.secretName=gitlab-site-tls \
  --set postgresql.image.tag=13.6.0 \
  --set gitlab-runner.runners.privileged=true \
  --set gitlab-runner.install=false
  
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

## Login to GitLab
### Initial login

```shell 
kubectl get secret gitlab-gitlab-initial-root-password -ojsonpath='{.data.password}' | base64 --decode ; echo
bFy0AsRDOFbI01NyWKZfyW6JIQiTEGvxRTcF4nNPjdbQb2Wa7mj42cBcfslNHSn4

# Update CA Root cert on client（self-signed）
cp ca.crt /usr/local/share/ca-certificates/
update-ca-certificates 

# 访问 https://gitlab.gitops.com
```

### But Invalid login or password.

```shell
# Wanna Change Password ?
kubectl exec -it -n gitlab-system gitlab-toolbox-5479cb9f8f-tfjwd -- bash
gitlab-rails console
--------------------------------------------------------------------------------
 Ruby:         ruby 2.7.7p221 (2022-11-24 revision 168ec2b1e5) [x86_64-linux]
 GitLab:       15.8.1 (383efe57adf) FOSS
 GitLab Shell: 14.15.0
 PostgreSQL:   13.6
------------------------------------------------------------[ booted in 11.16s ]
Loading production environment (Rails 6.1.6.1)
irb(main):001:0> u=User.where(id:1).first
=> #<User id:1 @root>
irb(main):008:0> u.password_confirmation='bFy0AsRDOFbI01NyWKZfyW6JIQiTEGvxRTcF4nNPjdbQb2Wa7mj42cBcfslNHSn4'
=> "bFy0AsRDOFbI01NyWKZfyW6JIQiTEGvxRTcF4nNPjdbQb2Wa7mj42cBcfslNHSn4"
irb(main):009:0> u.password='bFy0AsRDOFbI01NyWKZfyW6JIQiTEGvxRTcF4nNPjdbQb2Wa7mj42cBcfslNHSn4'
=> "bFy0AsRDOFbI01NyWKZfyW6JIQiTEGvxRTcF4nNPjdbQb2Wa7mj42cBcfslNHSn4"
irb(main):010:0> u.save
=> true
irb(main):011:0> exit

# Now you can login
```