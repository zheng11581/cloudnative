# 1. Node information and Node healthy
## Node && Lease && cAdvisor
```shell
# 节点信息kubelet通过Node对象上报
kubectl describe Node 192.168.101.69
...
Lease:
  HolderIdentity:  192.168.101.69
  AcquireTime:     <unset>
  RenewTime:       Sun, 06 Mar 2022 09:05:13 +0800
Capacity:
  cpu:                16
  ephemeral-storage:  494819Mi
  hugepages-1Gi:      0
  hugepages-2Mi:      0
  memory:             32778484Ki
  pods:               256
Allocatable:
  cpu:                15800m
  ephemeral-storage:  466969794197
  hugepages-1Gi:      0
  hugepages-2Mi:      0
  memory:             31652084Ki
  pods:               256
...


# 节点健康状态Kubelet通过Lease对象上报
k get lease -n kube-node-lease <nodeName> -oyaml

apiVersion: coordination.k8s.io/v1
kind: Lease
metadata:
  creationTimestamp: "2022-01-20T08:37:12Z"
  name: 192.168.101.69
  namespace: kube-node-lease
  ownerReferences:
  - apiVersion: v1
    kind: Node
    name: 192.168.101.69
    uid: 8f2ce016-6900-41c1-bbe6-14213fa85ffe
  resourceVersion: "36978738"
  uid: ecb32cea-f732-4f60-b70a-e50cb907ab0d
spec:
  holderIdentity: 192.168.101.69
  leaseDurationSeconds: 40
  renewTime: "2022-03-06T00:50:44.975545Z"

# 节点资源使用情况通过cAdvisor上报

```

### Reserved resource
```shell
--system-reserved
--kube-reserved
```


## 2. nodefs

### containerd & docker

`/var/lib/kubelet/pods/ecbbaf5e-c7f1-4f00-8627-454f03cfbb4d`

## 3. imagefs

### containerd

`/var/lib/containerd/io.containerd.snapshotter.v1.overlayfs`

### docker

`/var/lib/docker/overlay2/`

## Eviction configuration

```yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
evictionHard:
  memory.available: '500Mi'
  nodefs.available: '1Gi'
  imagefs.available: '100Gi'
evictionMinimumReclaim:
  memory.available: '0Mi'
  nodefs.available: '500Mi'
  imagefs.available: '2Gi'
```

### KubeletConfiguration file
```shell
/var/lib/kubelet/config.yaml
```

### Default KubeletConfiguration
```shell
kubeadm config print init-defaults --component-configs KubeletConfiguration
```