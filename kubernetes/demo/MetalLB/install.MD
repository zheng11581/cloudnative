## Config ARP
ipvs配置打开strictARP

```shell
# 查看 kube-proxy 中的 strictARP 配置
kubectl get configmap -n kube-system kube-proxy -o yaml | grep strictARP

# 手动修改 strictARP 配置为true
$ kubectl edit configmap -n kube-system kube-proxy

# 重启 kube-proxy 确保配置生效
$ kubectl rollout restart ds kube-proxy -n kube-system

# 确认配置生效
$ kubectl get configmap -n kube-system kube-proxy -o yaml | grep strictARP
      strictARP: true
```

## Install MetalLB

```shell
helm repo add metallb https://metallb.github.io/metallb

helm update repo

helm install metallb metallb/metallb --create-namespace --namespace=metallb-system

```

## Config MetalLB

### Configure the lmit of IPAddressPools

```shell
kubectl apply -f examples/metallb-configmap.yaml

```

## Create IPAddressPools

```shell
kubectl apply -f examples/metallb-ipaddresspools.yaml
```

## Test MetalLB

```shell


```
