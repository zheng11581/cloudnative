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
```


#### 2. Because the hostname of certification is harbor 
```shell
vim /etc/hosts
192.168.110.72 harbor
```

