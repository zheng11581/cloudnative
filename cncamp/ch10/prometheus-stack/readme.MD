### Get Helm Repository Info

```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

### Install Helm Chart

```shell
helm install demo prometheus-community/kube-prometheus-stack
```

### Confirm The Deploy Completed

```shell
# kubectl get pod
NAME                                    READY   STATUS    RESTARTS   AGE
alertmanager-kube-prom-alertmanager-0   2/2     Running   0          73m
prometheus-kube-prom-prometheus-0       2/2     Running   0          73m
grafana-85df8545b8-mh8jc                3/3     Running   0          73m
kube-prom-operator-b876bcd4b-x684m      1/1     Running   0          73m
kube-state-metrics-5b98d5b7c6-5g78h     1/1     Running   0          73m
prometheus-node-exporter-jv9z4          1/1     Running   0          73m

```

### Expose Grafana Pod

```shell
kubectl port-forward grafana-85df8545b8-mh8jc --address 0.0.0.0 3000:3000
```

### Get Grafana User/Passwd

```shell
# kubectl get secrets prometheus-stack-grafana -oyaml
apiVersion: v1
data:
  admin-password: cHJvbS1vcGVyYXRvcg==
  admin-user: YWRtaW4=
  ldap-toml: ""
kind: Secret
metadata:
  annotations:
    meta.helm.sh/release-namespace: default
  creationTimestamp: "2022-07-05T08:25:49Z"
  labels:
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: grafana
    app.kubernetes.io/version: 9.0.1
    helm.sh/chart: grafana-6.31.0
  name: prometheus-stack-grafana
  namespace: default
  resourceVersion: "957574"
  uid: f6234f82-1d17-4691-a9a7-4c0271406c76
type: Opaque

# echo "cHJvbS1vcGVyYXRvcg==" |base64 -d
prom-operator

# echo "YWRtaW4=" |base64 -d
admin
```

### Login From http://<nodeip>:3000