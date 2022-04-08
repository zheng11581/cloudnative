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
[https://kubernetes.io/zh/docs/reference/access-authn-authz/rbac/#command-line-utilities](https://kubernetes.io/zh/docs/reference/access-authn-authz/rbac/#command-line-utilities)

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
[https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#drain](https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands#drain)

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

### 1.4 etcd备份恢复

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
ETCDCTL_API=3 etcdctl --endpoints=https://127.0.0.1:2379 \
--cacert=/etc/kubernetes/pki/etcd/ca.crt \
--cert=/etc/kubernetes/pki/etcd/peer.crt --key=/etc/kubernetes/pki/etcd/peer.key \
snapshot save /srv/data/etcd-snapshot.db

# restore
mkdir /opt/backup
cd /etc/kubernetes/manifests && mv kube-* /opt/backup

ETCDCTL_API=3 etcdctl --endpoints=https://127.0.0.1:2379 \
--cacert=/etc/kubernetes/pki/etcd/ca.crt \
--cert=/etc/kubernetes/pki/etcd/peer.crt --key=/etc/kubernetes/pki/etcd/peer.key \
snapshot restore /srv/data/etcd-snapshot.db --data-dir=/var/lib/etcd-restore

vim etcd.yaml
# 将 volume 配置的 path: /var/lib/etcd 改成/var/lib/etcd-restore
volumes:
- hostPath:
path: /etc/kubernetes/pki/etcd
type: DirectoryOrCreate
name: etcd-certs
- hostPath:
path: /var/lib/etcd-restore

# 还原 k8s 组件
mv /opt/backup/* /etc/kubernetes/manifests
systemctl restart kubelet
```

### 1.5 网络策略 NetworkPolicy

#### Task
```text
创建一个名字为 all-port-from-namespace 的 NetworkPolicy，这个 NetworkPolicy 允许 internal
命名空间下的 Pod 访问该命名空间下的 9000 端口。
并且
不允许非 internal 命令空间的下的 Pod 访问
不允许访问没有监听 9000 端口的 Pod。
```


#### Ref
[https://kubernetes.io/zh/docs/concepts/services-networking/network-policies/](https://kubernetes.io/zh/docs/concepts/services-networking/network-policies/)

#### Operation
```shell
vim 5-psp.yaml
kubectl apply -f 5-psp.yaml
```

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: all-port-from-namespace
  namespace: internal
spec:
  podSelector: {}
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector: {}
      ports:
        - protocol: TCP
          port: 9000
```

### 1.6 四层负载均衡 service

#### Task
```text
1. 重新配置已经存在的deployment front-end，并且新增一个名为http的端口来暴露现有容器nginx 80/tcp
2. 创建一个service front-end-svc，暴露容器nginx名为http的端口，并且 service 的类型为 NodePort。
```

#### Ref
[https://kubernetes.io/zh/docs/concepts/services-networking/connect-applications-service/](https://kubernetes.io/zh/docs/concepts/services-networking/connect-applications-service/)


#### Operation
```shell
vim 6-service.yaml
kubectl apply -f 6-service.yaml
kubectl expose deploy front-end --name=front-end-svc --port=80 --target-port=http --type=NodePort

```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: front-end
spec:
  selector:
    matchLabels:
      run: front-end
  replicas: 2
  template:
    metadata:
      labels:
        run: front-end
    spec:
      containers:
        - name: front-end
          image: nginx
          ports:
            - containerPort: 80
              name: http
              protocal: TCP
```

### 1.7 七层负载均衡 ingress

#### Task
```text
新建一个nginx ingress资源对象：
1. 名字叫 pong
2. 命名空间是 ing-internal
3. 通过路径:5678/hi来暴露service hi
```

#### Ref
[https://kubernetes.io/zh/docs/concepts/services-networking/ingress/](https://kubernetes.io/zh/docs/concepts/services-networking/ingress/)

#### Operation
```shell
vim 7-ingress.yaml
kubectl apply -f 7-ingress.yaml
```

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: pong
  namespace: ing-internal
  annotations:
    nginx.ingress.kubernetes.io/class: nginx-ingress
spec:
  rules:
    - http:
        paths:
          - path: /hi
            pathType: Prefix
            backend:
              service:
                name: hi
                port:
                  number: 5678

```

### 1.8 Deployment 扩缩容

#### Task
```text
扩容名字为 loadbalancer 的 deployment 的副本数为 6
```

#### Ref

#### Operation
```shell
kubect scale --replicas=6 deployment loadbalancer 
```
