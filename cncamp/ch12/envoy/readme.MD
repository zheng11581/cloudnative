### Deploy simple

```sh
kubectl create -f simple.yaml
```

### Deploy envoy

```sh
kubectl create configmap envoy-config --from-file=envoy.yaml
kubectl create -f envoy-deploy.yaml
```

### curl envoy pod ip
```shell
kubectl run -it curl --image=zheng11581/curl --restart=Never -- bash
root@curl:/# curl 10.0.1.12:10000
Hello, [stranger]
Reading request headers to response...
X-Request-Id=[0f9e54a8-e934-4a64-9786-14fac0d3d569]
X-Envoy-Expected-Rq-Timeout-Ms=[15000]
User-Agent=[curl/7.35.0]
Accept=[*/*]
X-Forwarded-Proto=[http]
```

### Envoy doc
[https://www.envoyproxy.io/docs/envoy/latest/start/start](https://www.envoyproxy.io/docs/envoy/latest/start/start)