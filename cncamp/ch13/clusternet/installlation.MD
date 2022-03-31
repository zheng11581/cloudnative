## Install kubectl clusternet plugin
### Install krew
```shell
(
  set -x; cd "$(mktemp -d)" &&
  OS="$(uname | tr '[:upper:]' '[:lower:]')" &&
  ARCH="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')" &&
  KREW="krew-${OS}_${ARCH}" &&
  curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" &&
  tar zxvf "${KREW}.tar.gz" &&
  ./"${KREW}" install krew
)

```
### Setup krew
```shell
export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"
kubectl krew
```
### Install Clusternet plugin
```shell
kubectl krew update
kubectl krew install clusternet
kubectl clusternet version
```
## Setup Clusternet control plane

### Install Clusternet Hub to PARENT CLUSTER
```shell
https://github.com/clusternet/charts/tree/main/charts/clusternet-hub

helm repo add clusternet https://clusternet.github.io/charts
helm install clusternet-hub -n clusternet-system --create-namespace clusternet/clusternet-hub
kubectl apply -f https://raw.githubusercontent.com/clusternet/clusternet/main/manifests/samples/cluster_bootstrap_token.yaml

helm delete clusternet-hub -n clusternet-system
```

### Install Clusternet Scheduler to PARENT CLUSTER
```shell
https://github.com/clusternet/charts/tree/main/charts/clusternet-scheduler

helm repo add clusternet https://clusternet.github.io/charts
helm install clusternet-scheduler -n clusternet-system --create-namespace \
  clusternet/clusternet-scheduler
  
helm delete clusternet-scheduler -n clusternet-system
```

### Install Clusternet Agent to CHILD CLUSTER ,registrationToken is installed in kube-system
```shell
https://github.com/clusternet/charts/tree/main/charts/clusternet-agent

helm install clusternet-agent -n clusternet-system --create-namespace \
--set parentURL=https://xxx:443 \
--set registrationToken=07401b.f395accd246ae52d \
clusternet/clusternet-agent

helm delete clusternet-agent -n clusternet-system
```
### check managed clusters and get cluster-id
```shell
kubectl get managedcluster -A
```
