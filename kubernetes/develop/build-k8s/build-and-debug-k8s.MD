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
# apt install resolvconf
# systemctl enable --now resolvconf.service
# vi /etc/resolvconf/resolv.conf.d/head
# resolvconf -u

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

### Startup local cluster

```shell
# ./hack/local-up-cluster.sh -O
```

### Local Debug

```shell
# ps -ef |grep kube-apiserver
# APISERVER_CMD=/root/go/src/k8s.io/kubernetes/_output/bin/kube-apiserver
# APISERVER_ARGS="--authorization-mode=Node,RBAC  --cloud-provider= --cloud-config=   --v=3 --vmodule= --audit-policy-file=/tmp/kube-audit-policy-file --audit-log-path=/tmp/kube-apiserver-audit.log --authorization-webhook-config-file= --authentication-token-webhook-config-file= --cert-dir=/var/run/kubernetes --egress-selector-config-file=/tmp/kube_egress_selector_configuration.yaml --client-ca-file=/var/run/kubernetes/client-ca.crt --kubelet-client-certificate=/var/run/kubernetes/client-kube-apiserver.crt --kubelet-client-key=/var/run/kubernetes/client-kube-apiserver.key --service-account-key-file=/tmp/kube-serviceaccount.key --service-account-lookup=true --service-account-issuer=https://kubernetes.default.svc --service-account-jwks-uri=https://kubernetes.default.svc/openid/v1/jwks --service-account-signing-key-file=/tmp/kube-serviceaccount.key --enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,Priority,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota,NodeRestriction --disable-admission-plugins= --admission-control-config-file= --bind-address=0.0.0.0 --secure-port=6443 --tls-cert-file=/var/run/kubernetes/serving-kube-apiserver.crt --tls-private-key-file=/var/run/kubernetes/serving-kube-apiserver.key --storage-backend=etcd3 --storage-media-type=application/vnd.kubernetes.protobuf --etcd-servers=http://127.0.0.1:2379 --service-cluster-ip-range=10.0.0.0/24 --feature-gates=AllAlpha=false --external-hostname=localhost --requestheader-username-headers=X-Remote-User --requestheader-group-headers=X-Remote-Group --requestheader-extra-headers-prefix=X-Remote-Extra- --requestheader-client-ca-file=/var/run/kubernetes/request-header-ca.crt --requestheader-allowed-names=system:auth-proxy --proxy-client-cert-file=/var/run/kubernetes/client-auth-proxy.crt --proxy-client-key-file=/var/run/kubernetes/client-auth-proxy.key --cors-allowed-origins=/127.0.0.1(:[0-9]+)?$,/localhost(:[0-9]+)?$"

# dlv exec $APISERVER_CMD -- $APISERVER_ARGS
```

### Remote Debug

```shell
# dlv exec /root/go/src/k8s.io/kubernetes/_output/bin/kube-apiserver --headless --listen=:12345 --log --log-output=debugger,gdbwire,lldbout,debuglineerr,rpc,dap,fncall,minidump --log-dest=/root/k8s-dev/apiserver.log --api-version=2 -- "--authorization-mode=Node,RBAC  --cloud-provider= --cloud-config=   --v=3 --vmodule= --audit-policy-file=/tmp/kube-audit-policy-file --audit-log-path=/tmp/kube-apiserver-audit.log --authorization-webhook-config-file= --authentication-token-webhook-config-file= --cert-dir=/var/run/kubernetes --egress-selector-config-file=/tmp/kube_egress_selector_configuration.yaml --client-ca-file=/var/run/kubernetes/client-ca.crt --kubelet-client-certificate=/var/run/kubernetes/client-kube-apiserver.crt --kubelet-client-key=/var/run/kubernetes/client-kube-apiserver.key --service-account-key-file=/tmp/kube-serviceaccount.key --service-account-lookup=true --service-account-issuer=https://kubernetes.default.svc --service-account-jwks-uri=https://kubernetes.default.svc/openid/v1/jwks --service-account-signing-key-file=/tmp/kube-serviceaccount.key --enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,Priority,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota,NodeRestriction --disable-admission-plugins= --admission-control-config-file= --bind-address=0.0.0.0 --secure-port=6443 --tls-cert-file=/var/run/kubernetes/serving-kube-apiserver.crt --tls-private-key-file=/var/run/kubernetes/serving-kube-apiserver.key --storage-backend=etcd3 --storage-media-type=application/vnd.kubernetes.protobuf --etcd-servers=http://127.0.0.1:2379 --service-cluster-ip-range=10.0.0.0/24 --feature-gates=AllAlpha=false --external-hostname=localhost --requestheader-username-headers=X-Remote-User --requestheader-group-headers=X-Remote-Group --requestheader-extra-headers-prefix=X-Remote-Extra- --requestheader-client-ca-file=/var/run/kubernetes/request-header-ca.crt --requestheader-allowed-names=system:auth-proxy --proxy-client-cert-file=/var/run/kubernetes/client-auth-proxy.crt --proxy-client-key-file=/var/run/kubernetes/client-auth-proxy.key --cors-allowed-origins=/127.0.0.1(:[0-9]+)?$,/localhost(:[0-9]+)?$"

# dlv connect :12345
```