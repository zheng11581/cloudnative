### Install istio

```sh
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.13.2
cp bin/istioctl /usr/local/bin
istioctl install --set profile=demo -y
```

### Istio monitoring

- Grafana dashboard `7639`