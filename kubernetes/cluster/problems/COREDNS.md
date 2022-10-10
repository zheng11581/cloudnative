# 故障排查：Kubernetes 中 Pod 无法正常解析域名

## 问题描述

- Kubernetes 升级到 1.18.1 版本
- 工作节点的部分 Pod 无法启动，查看消息全是 connetion timeout 的问题
  - 域名方式连接集群内部服务
  - 域名方式连接集群外部地址
  - 通过 IP 进行远程连接的应用倒是没有问题
- 初步怀疑很可能是 DNS 出现了问题

## 部署 DNS 调试工具

为了探针是否为 DNS 问题，这里需要提前部署用于测试 DNS 问题的 dnsutils 镜像，该镜像中包含了用于测试 DNS 问题的工具包，非常利于我们分析与发现问题。接下来，我们将在 Kubernetes 中部署这个工具镜像。

### 创建 DNS 工具 Pod 部署文件

创建用于部署的 Deployment 资源文件 ndsutils.yaml：

ndsutils.yaml
```yaml
apiVersion: v1
kind: Pod
metadata:
name: dnsutils
spec:
containers:
- name: dnsutils
image: mydlqclub/dnsutils:1.3
imagePullPolicy: IfNotPresent
command: ["sleep","3600"]
```

### 通过 Kubectl 工具部署 DNS 工具镜像

通过 Kubectl 工具，将对上面 DNS 工具镜像部署到 Kubernetes 中：

```shell
# kubectl create -f ndsutils.yaml -n kube-system
```

## 问题分析

### 进入 DNS 工具 Pod 的命令行

### 