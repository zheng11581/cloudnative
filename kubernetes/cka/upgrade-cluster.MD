
# 确定要升级到哪个版本

```shell
yum list --showduplicates kubeadm --disableexcludes=kubernetes # master only

# 在列表中查找最新的 1.24 版本
# 它看起来应该是 1.24.x-0，其中 x 是最新的补丁版本
```

# 升级计划

```shell
yum install -y kubeadm-1.24.x-0 --disableexcludes=kubernetes

kubeadm upgrade plan # master only

sudo kubeadm upgrade apply v1.24.x # master only
sudo kubeadm upgrade node # node only

```

# 排空

```shell
kubectl cordon <nodename>

kubectl drain <nodename> --delete-emptydir-data --ignore-daemonsets --force

kubectl get node -owide
```

# 升级kubelet

```shell
yum install -y kubelet-1.24.x-0 kubectl-1.24.x-0 --disableexcludes=kubernetes

sudo systemctl daemon-reload

sudo systemctl restart kubelet
```

# 恢复节点

```shell
kubectl uncordon <nodename>

kubectl get node -owide
```