## Install with Helm

```shell
helm repo add harbor https://helm.goharbor.io

helm repo update

kubectl apply -f harbor-site-tls.yaml

helm update harbor harbor/harbor -f values.yaml --namespace harbor --create-namespace
```

## Check harbor-core Pod configuration data.EXT_ENDPOINT

```shell
# kubectl edit cm -n harbor harbor-core
EXT_ENDPOINT: https://harbor.gitops.com:31137/ # This value defined in values.yaml: externalURL, pls update to that value
```

