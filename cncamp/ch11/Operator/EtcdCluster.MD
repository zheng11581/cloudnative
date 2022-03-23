### kubebuilder 安装
```shell
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
```

### CR设计

```yaml
# size 就是集群大小，image是我们需要制定的image
apiVersion: etcd.fmeng.io/v1
kind: EtcdCluster
metadata: 
  name: etcd-fmeng
spec:  
  size: 5  
  image: quay.io/coreos/etcd:v3.3.8
```


### 初始化项目

```shell
--domain cncamp.io examle.com
--owner fmeng
--repo gitee.com/zheng11581/etcd-cluster

kubebuilder init --domain cncamp.io --repo gitee.com/zheng11581/etcd-cluster --owner fmeng
```

### 创建API

```shell
kubebuilder create api --group apps --version v1beta1 --kind EtcdCluster
```

### Coding

1. api/v1beta1/etcdcluster_types.go

```go
// EtcdCluster is the Schema for the etcdclusters API
type EtcdCluster struct {
metav1.TypeMeta   `json:",inline"`
metav1.ObjectMeta `json:"metadata,omitempty"`

Spec   EtcdClusterSpec   `json:"spec,omitempty"`
Status EtcdClusterStatus `json:"status,omitempty"`
}
```

2. controllers/resource.go

sts和headless svc的逻辑创建etcd集群要做两个事情，1、建statefulset，2、建headless svc
核心思路是引入apps和core的包，然后对照我们之前的两个spec向里面填充相对应的内容
```go
appsv1 "k8s.io/api/apps/v1"
corev1 "k8s.io/api/core/v1"
```

Reconcile当创建逻辑写好之后，就要用reconcile逻辑去不断对比实际状态和spec状态
主体就两个逻辑，CreateOrUpdate Service / CreateOrUpdate sts

### Testing
```shell
make install
kubectl get crd | grep etcd
make run
```

```yaml
apiVersion: etcd.cncamp.io/v1beta1
kind: EtcdCluster
metadata:
  name: etcd-fmeng
  spec:
    size: 3
    image: quay.io/coreos/etcd:v3.3.8
```

