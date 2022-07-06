### Clone kube-prometheus

```shell
git clone https://github.com/prometheus-operator/kube-prometheus.git
```

### Create the monitoring stack using the config in the manifests directory

```shell
# Create the namespace and CRDs, and then wait for them to be available before creating the remaining resources
# kubectl apply --server-side -f manifests/setup
# until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done
# kubectl apply -f manifests/
```

### And to teardown the stack

```shell
# kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup
```

### Access Prometheus UIs

```shell
# kubectl --namespace monitoring port-forward svc/prometheus-k8s --address 0.0.0.0 9090
# Then access via http://localhost:9090
```

### Access Grafana UIs

```shell
# kubectl --namespace monitoring port-forward svc/grafana --address 0.0.0.0 3000
# Then access via http://localhost:3000 and use the default grafana user:password of admin:admin
```

### Access Alert Manager UIs

```shell
# kubectl --namespace monitoring port-forward svc/alertmanager-main --address 0.0.0.0 9093
# Then access via http://localhost:9093
```
