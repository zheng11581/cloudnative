### containerd pull images (Public repo)

#### Prepared: You have installed Harbor and configured Harbor

#### 1. config containerd
```shell
          
vim /etc/crictl.yaml 

runtime-endpoint: unix:///run/containerd/containerd.sock
image-endpoint: unix:///run/containerd/containerd.sock
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
  
kubectl create secret docker-registry regcred -n devops \
  --docker-server=registry.cn-beijing.aliyuncs.com \
  --docker-username=gobroadyun \
  --docker-password=xxx
```

#### 4. 
