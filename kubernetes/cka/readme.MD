### bash-completion
```shell
kubectl completion -h

  # Installing bash completion on Linux
  ## If bash-completion is not installed on Linux, install the 'bash-completion' package
  ## via your distribution's package manager.
  ## Load the kubectl completion code for bash into the current shell
  source <(kubectl completion bash)
  ## Write bash completion code to a file and source it from .bash_profile
  kubectl completion bash > ~/.kube/completion.bash.inc
  printf "
  # Kubectl shell completion
  source '$HOME/.kube/completion.bash.inc'
  " >> $HOME/.bash_profile
  source $HOME/.bash_profile

  # Load the kubectl completion code for zsh[1] into the current shell
  source <(kubectl completion zsh)
  # Set the kubectl completion code for zsh[1] to autoload on startup
  kubectl completion zsh > "${fpath[1]}/_kubectl"
```

### 1.1 基于角色的访问控制-RBAC

#### Task 
`创建一个名为 deployment-clusterrole 的 clusterrole，该 clusterrole 只允许创建 Deployment、
Daemonset、Statefulset 的 create 操作
在名字为 app-team1 的 namespace 下创建一个名为 cicd-token 的 serviceAccount，并且将上一步创
建 clusterrole 的权限绑定到该 serviceAccount`

#### Ref
[https://kubernetes.io/docs/reference/access-authn-authz/rbac/#command-line-utilities](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#command-line-utilities)

#### Operation
```shell
kubectl create clusterrole deployment-clusterrole --verb=create --
resource=deployments,statefulsets,daemonsets
kubectl -n app-team1 create serviceaccount cicd-token
kubectl -n app-team1 create rolebinding cicd-token-binding --clusterrole=deployment-clusterrole --serviceaccount=app-team1:cicd-token

```

### 1.2 节点维护-指定 node 节点不可用

#### Task
`将 ek8s-node-1 节点设置为不可用，然后重新调度该节点上的所有 Pod`

#### Ref

#### Operation
```shell
kubectl config use-context ek8s
kubectl cordon ek8s-node-1 #设置节点是不可调度状态
kubectl drain ek8s-node-1 --delete-emptydir-data --ignore-daemonsets --force
```

### 1.3 K8s 版本升级

#### Task

`有的 Kubernetes 集权正在运行的版本是 1.23.1，仅将主节点上的所有 kubernetes 控制面板和
组件升级到版本 1.23.4 另外，在主节点上升级 kubelet 和 kubectl`

#### Ref
[https://kubernetes.io/zh/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/](https://kubernetes.io/zh/docs/tasks/administer-cluster/kubeadm/kubeadm-upgrade/)

#### Operation
```shell
kubectl cordon cn-node1
kubectl drain cn-node1 --delete-emptydir-data --ignore-daemonsets --force

apt update
apt-cache policy kubeadm |grep 1.23.4
apt-get install kubeadm=1.23.4-00
kubeadm upgrade plan
kubeadm upgrade apply v1.23.4
apt-get install kubelet=1.23.4-00 kubectl=1.23.4-00
systemctl daemon-reload
systemctl restart kubelet.service

kubectl uncordon cn-node1

ERROR: kubeadm upgrade plan 遇到这个错误, 通过reset重新加入集群
[[upgrade/config] FATAL: failed to getAPIEndpoint: could not retrieve API endpoints for node "cn-node1" using pod annotations: timed out waiting for the condition]
```

### 1.4 etcd升级备份

#### Task
```text
针对 etcd 实例 https://127.0.0.1:2379 创建一个快照，保存到/srv/data/etcd-snapshot.db。在创
建快照的过程中，如果卡住了，就键入 ctrl+c 终止，然后重试。
然后恢复一个已经存在的快照： /var/lib/backup/etcd-snapshot-previous.db
执行 etcdctl 命令的证书存放在：
ca 证书：/opt/KUIN00601/ca.crt
客户端证书：/opt/KUIN00601/etcd-client.crt
客户端密钥：/opt/KUIN00601/etcd-client.key
```

#### Ref
[https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/#backing-up-an-etcd-cluster](https://kubernetes.io/docs/tasks/administer-cluster/configure-upgrade-etcd/#backing-up-an-etcd-cluster)

#### Operation
```shell
# backup
ETCDCTL_API=3 
etcdctl --endpoints=https://127.0.0.1:2379 \
--cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/etcd/peer.crt --key=/etc/kubernetes/pki/etcd/peer.key \
snapshot save /srv/data/etcd-snapshot.db

# restore
mkdir /opt/backup
cd /etc/kubernetes/manifests && mv kube-* /opt/backup
etcdctl --endpoints="https://127.0.0.1:2379" --cacert=/opt/KUIN000601/ca.crt --
cert=/opt/KUIN000601/etcd-client.crt --key=/opt/KUIN000601/etcd-client.key snapshot 
restore /var/lib/backup/etcd-snapshot-previous.db --data-dir=/var/lib/etcd-restore
```
