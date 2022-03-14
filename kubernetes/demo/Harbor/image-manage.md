### containerd pull images (Public repo)

#### Prepared: You have installed Harbor and configured Harbor

#### 1. config containerd
```shell
vim /etc/containerd/config.toml

      [plugins."io.containerd.grpc.v1.cri".registry.configs]
        [plugins."io.containerd.grpc.v1.cri".registry.configs."192.168.110.72".tls]
          insecure_skip_verify = true

      [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
        [plugins."io.containerd.grpc.v1.cri".registry.mirrors."192.168.110.72"]
          endpoint = ["https://192.168.110.72"]
          
vim /etc/crictl.yaml 

runtime-endpoint: unix:///var/run/containerd/containerd.sock
image-endpoint: unix:///var/run/containerd/containerd.sock
timeout: 2
debug: false
pull-image-on-create: false

```


#### 2. Because the hostname of certification is goharbor.com 
```shell
vim /etc/hosts
192.168.110.72 goharbor
```

#### 3. Create credential for CICD pipeline

```shell
kubectl create secret docker-registry regcred -n devops \
  --docker-server=192.168.110.72 \
  --docker-username=admin \
  --docker-password=xxx
```

#### 4. 
