### Setup tekton

```sh
# Tekton Pipelines
kubectl apply --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml

# Tekton Triggers 
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml
kubectl apply --filename https://storage.googleapis.com/tekton-releases/triggers/latest/interceptors.yaml

# Tekton Dashboard
kubectl apply --filename https://github.com/tektoncd/dashboard/releases/latest/download/tekton-dashboard-release.yaml
```

### Setup tekton cli tkn

```sh
curl -LO https://github.com/tektoncd/cli/releases/download/v0.20.0/tektoncd-cli-0.20.0_Linux-64bit.deb
sudo dpkg -i ./tektoncd-cli-0.20.0_Linux-64bit.deb
which tkn
```

### Docs
```shell
https://tekton.dev/docs/getting-started/
```

