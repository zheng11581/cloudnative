## Prepare Dev enviroment

### Isntall Etcd
```shell
# source /etc/profile.d/etcd.sh

```

### Install Containerd
```shell
# tar -zxf cri-containerd-cni-1.6.19-linux-amd64.tar.gz -C /
# containerd config default | sudo tee /etc/containerd/config.toml
# vi /etc/containerd/config.toml

      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."docker.io"]
          endpoint = ["https://xt4zv57n.mirror.aliyuncs.com", "https://registry-1.docker.io"]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."k8s.gcr.io"]
          endpoint = ["https://xt4zv57n.mirror.aliyuncs.com"]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."registry.k8s.io"]
          endpoint = ["https://xt4zv57n.mirror.aliyuncs.com"]

# systemctl daemon-reload
# systemctl containerd restart
 
```
### Install Golang
```shell
# wget https://go.dev/dl/go1.19.7.linux-amd64.tar.gz -o go1.19.7.linux-amd64.tar.gz
# source /etc/profile.d/golang.sh
# go env -w GOPROXY="https://goproxy.cn,direct"
# go env -w GO111MODULE="on"
```
### Run local cluster
```shell
# crictl pull docker.io/coredns/coredns:1.9.3
# crictl images
# ctr images tag docker.io/coredns/coredns:1.9.3 registry.k8s.io/coredns/coredns:v1.9.3

# mkdir -p $GOPATH/bin
# mkdir -p $GOPATH/pkg
# mkdir -p $GOPATH/src/k8s.io
# cd $GOPATH/src/k8s.io
# git clone git@github.com:kubernetes/kubernetes.git
# cd kubernetes
# git checkout -b my-v1.26.2 v1.26.2
# vi hack/lib/golang.sh
    # 注释
    #if [[ "${DBG:-}" == 1 ]]; then
    #    # Debugging - disable optimizations and inlining and trimPath
    #    gogcflags="${GOGCFLAGS:-} all=-N -l"
    #    goasmflags=""
    #fi
    # 新增
    gogcflags="${GOGCFLAGS:-} all=-N -l"

    goldflags="all=$(kube::version::ldflags) ${GOLDFLAGS:-}"
    # 注释
    #if [[ "${DBG:-}" != 1 ]]; then
    #    # Not debugging - disable symbols and DWARF.
    #    goldflags="${goldflags} -s -w"
    #fi

# make all GOFLAGS=-v GOGCFLAGS="-N -l"
# ./hack/local-up-cluster.sh -O

# export KUBERNETES_PROVIDER=local
# cluster/kubectl.sh config set-cluster local --server=https://localhost:6443 --certificate-authority=/var/run/kubernetes/server-ca.crt
# source /etc/profile.d/kubernetes.sh 
# kubectl.sh get node

```

## Debug Kubernetes

### Install go-delve

```shell
# go install github.com/go-delve/delve/cmd/dlv@latest

```
