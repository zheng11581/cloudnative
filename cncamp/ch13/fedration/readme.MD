## install kubefed by kind
### 下载federation代码
```
git clone https://github.com/kubernetes-sigs/kubefed.git
```
### 选择HostCluster，确认kubeconfig符合federatio命名规范，用vi编辑kubeconfig，确保context属性没用@字符
```
vi ~/.kube/config 
```
```
contexts:
- context:
  cluster: kubernetes
  user: kubernetes-admin
  name: `cluster1`
  current-context: `cluster1`
```


### 安装
```
kind create cluster --image kindest/node:v1.23.5 --name federation
cd kubefed
cp bin/kubefedctl-linux-amd64 /usr/local/bin/kubefedctl
helm repo add kubefed-charts https://raw.githubusercontent.com/kubernetes-sigs/kubefed/master/charts
helm --namespace kube-federation-system upgrade -i kubefed kubefed-charts/kubefed --version=0.9.1 --create-namespace
kubectl get all -n kube-federation-system

# 配置
kubectl config set-cluster --kubeconfig=config-demo gke-us-01 --server=https://35.192.141.53
kubectl config set-credentials --kubeconfig=config-demo gke-us-01 --auth-provider=gcp \
--auth-provider-arg=access-token="ya29.A0ARrdaM-LYmpQyYHbi9Ak5-fZz2bmjDWts8TWq2xMqhi5NurGod9Aik4Kt7YXrb9Xdr2Go-1KP_Y6Dr88HUlCBW8hg41msXgIVisE5rbOR6F1ohcvVaO6Emorjc4eRQVx3ZNzlSofoKso4IpWpFjfjpSpaHJiRGoopVvUAw" \
--auth-provider-arg=cmd-args="config config-helper --format=json" \
--auth-provider-arg=cmd-path=/usr/local/google-cloud-sdk/bin/gcloud \
--auth-provider-arg=expiry-key='{.credential.token_expiry}' \
--auth-provider-arg=token-key='{.credential.access_token}' 
kubectl config set-context --kubeconfig=config-demo gke-us-01 --cluster=gke-us-01 --user=gke-us-01
kubectl config use-context --kubeconfig=config-demo gke-us-01

https://kubernetes.io/zh/docs/tasks/access-application-cluster/configure-access-multiple-clusters/
https://betterprogramming.pub/build-a-federation-of-multiple-kubernetes-clusters-with-kubefed-v2-8d2f7d9e198a

kubectl -n kube-federation-system get kubefedcluster -oyaml
```
### fix
```
k get po -n kube-system kube-apiserver-kind-control-plane -owide
kube-system              kube-apiserver-kind-control-plane             1/1     Running   0          24m     172.18.0.2    kind-control-plane   <none>           <none>
172.18.0.2
```
#### edit kubefedcluster and change apiEndpoint to 172.18.0.2:6443
```
kubectl -n kube-federation-system edit kubefedcluster
```
### 安装完成后查看federation
```
kubectl get all -n kube-federation-system
```
### 将namespace设置为联邦对象
```
kubectl create ns federate-me
./bin/kubefedctl federate ns federate-me

```
### 创建FederatedDeployment
```
kubectl apply -f test-deployment.yaml -n federate-me
```
### 查看membercluster中的deployment对象
```
kubectl get deployment test-deployment -n federate-me
```
### 创建RSP对象
```
kubectl apply -f test-deployment-rsp.yaml -n federate-me
```
### 查看membercluster中的deployment对象
```
kubectl get deployment test-deployment -n federate-me
```
### 查看federatedployment，会发现因为rsp，federateddeployment的override属性被更新了
```
kubectl get federateddeployment test-deployment -n federate-me
```
### 删除rsp，发现override没有被拿掉