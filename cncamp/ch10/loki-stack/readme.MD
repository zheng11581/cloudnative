### Add grafana repo

```sh
helm repo add grafana https://grafana.github.io/helm-charts
```

### Install loki-stack

```sh
helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```

### If you get the following error, that means your k8s version is too new to install

```
Error: unable to build kubernetes objects from release manifest: [unable to recognize "": no matches for kind "ClusterRole" in version "rbac.authorization.k8s.io/v1beta1", unable to recognize "": no matches for kind "ClusterRoleBinding" in version "rbac.authorization.k8s.io/v1beta1", unable to recognize "": no matches for kind "Role" in version "rbac.authorization.k8s.io/v1beta1", unable to recognize "": no matches for kind "RoleBinding" in version "rbac.authorization.k8s.io/v1beta1"]
```

### Download loki-stack

```sh
helm pull grafana/loki-stack
tar -xvf loki-stack-*.tgz
cd loki-stack
```

### Replace all `rbac.authorization.k8s.io/v1beta1` with `rbac.authorization.k8s.io/v1` by
```sh
grep -rl "rbac.authorization.k8s.io/v1beta1" . | xargs sed -i 's/rbac.authorization.k8s.io\/v1beta1/rbac.authorization.k8s.io\/v1/g'
```

### Install loki locally

```sh
helm upgrade --install loki ./loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false
```

### Get kube-state-metrics image from docker.io
Go to the node where loki-kube-state-metrics scheduled

```shell
crictl pull bitnami/kube-state-metrics:2.2.4
ctr -n=k8s.io image tag docker.io/bitnami/kube-state-metrics:2.2.4 k8s.gcr.io/kube-state-metrics/kube-state-metrics:v2.2.4
```

### Change the grafana service to NodePort type and access it

```sh
kubectl edit svc loki-grafana -oyaml -n default
```

And change ClusterIP type to NodePort.

Login password is in secret `loki-grafana`

```sh
kubectl get secret loki-grafana -oyaml -n default
```

Find admin-password: `xxx`

```sh
echo 'xxx' | base64 -d
```

Then you will get grafana login password, the login username is 'admin' on default.
