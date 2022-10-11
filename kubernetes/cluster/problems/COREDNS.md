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
$ kubectl create -f ndsutils.yaml -n kube-system
```

## 问题分析

### 进入 DNS 工具 Pod 的命令行

上面 DNS 工具已经部署完成，我们可也通过 Kubectl 工具进入 Pod 命令行，然后，使用里面的一些工具进行问题分析，命令如下：

```shell
$ kubectl exec -it dnsutils /bin/sh -n kube-system
```
### 通过 Ping 和 Nsloopup 命令测试

进入容器 sh 命令行界面后，先使用 ping 命令来分别探测观察是否能够 ping 通集群内部和集群外部的地址，观察到的信息如下：

```shell
## Ping 集群外部，例如这里 ping 一下百度
$ ping www.baidu.com
ping: bad address 'www.baidu.com'

## Ping 集群内部 kube-apiserver 的 Service 地址
$ ping kubernetes.default
ping: bad address 'kubernetes.default'

## 使用 nslookup 命令查询域名信息
$ nslookup kubernetes.default
;; connection timed out; no servers could be reached

## 直接 Ping 集群内部 kube-apiserver 的 IP 地址
$ ping 10.96.0.1
PING 10.96.0.1 (10.96.0.1): 56 data bytes
64 bytes from 10.96.0.1: seq=0 ttl=64 time=0.096 ms
64 bytes from 10.96.0.1: seq=1 ttl=64 time=0.050 ms
64 bytes from 10.96.0.1: seq=2 ttl=64 time=0.068 ms

## 退出 dnsutils Pod 命令行
$ exit
```
可以观察到两次 ping 域名都不能 ping 通，且使用 nsloopup 分析域名信息时超时。然而，使用 ping kube-apiserver 的 IP 地址 “10.96.0.1” 则可以正常通信，所以，排除网络插件（Flannel、Calico 等）的问题。初步判断，很可能是 CoreDNS 组件的错误引起的某些问题，所以接下来我们测试 CoreDNS 是否正常。

### 检测 CoreDNS 应用是否正常运行

首先，检查 CoreDNS Pod 是否正在运行，如果 READY 为 0，则显示 CoreDNS 组件有问题：

```shell
$ kubectl get pods -l k8s-app=kube-dns -n kube-system
NAME                       READY   STATUS    RESTARTS   AGE
coredns-669f77d7cc-8pkpw   1/1     Running   2          6h5m
coredns-669f77d7cc-jk9wk   1/1     Running   2          6h5m
```

可也看到 CoreDNS 两个 Pod 均正常启动，所以再查看两个 Pod 中的日志信息，看看有无错误日志：

```shell
$ for p in $(kubectl get pods --namespace=kube-system -l k8s-app=kube-dns -o name); \
do kubectl logs --namespace=kube-system $p; done

.:53
[INFO] plugin/reload: Running configuration MD5 = 4e235fcc3696966e76816bcd9034ebc7
CoreDNS-1.6.7
linux/amd64, go1.13.6, da7f65b
.:53
[INFO] plugin/reload: Running configuration MD5 = 4e235fcc3696966e76816bcd9034ebc7
CoreDNS-1.6.7
linux/amd64, go1.13.6, da7f65b

```

通过上面信息可以观察到，日志中信息也是正常启动没有问题。再接下来，查看下 CoreDNS 的入口 Service "kube-dns" 是否存在：

```shell
$ kubectl get service -n kube-system | grep kube-dns
NAME                        TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)               
kube-dns                    ClusterIP   10.96.0.10      <none>        53/UDP,53/TCP,9153/TCP
# kube-dns 的 IP 为 10.96.0.10，集群内的 Pod 都是通过该 IP 与 DNS 组件进行交互，查询 DNS 信息
```


上面显示 Service “kube-dns” 也存在，但是 Service 是通过 endpoints 和 Pod 进行绑定的，所以看看这个 CoreDNS 的 Endpoints 是否存在，及信息是否正确：

```shell
$ kubectl get endpoints kube-dns -n kube-system
NAME       ENDPOINTS                                      
kube-dns   10.244.0.21:53,d10.244.2.82:53,10.244.0.21:9153

```

可以看到 Endpoints 配置也是正常的，正确的与两个 CporeDNS Pod 进行了关联。

经过上面一系列检测 CoreDNS 组件确实是正常运行。接下来，观察 CoreDNS 域名解析日志，进而确定 Pod 中的域名解析请求是否能够正常进入 CoreDNS。


### 观察 CoreDNS 域名解析日志信息

使用 kubectl edit 命令来修改存储于 Kubernetes ConfigMap 中的 CoreDNS 配置参数信息，添加 log 参数，让 CoreDNS 日志中显示域名解析信息：

CoreDNS 配置参数说明：
- errors：输出错误信息到控制台。
- health：CoreDNS 进行监控检测，检测地址为 http://localhost:8080/health 如果状态为不健康则让 Pod 进行重启。
- ready：全部插件已经加载完成时，将通过 endpoints 在 8081 端口返回 HTTP 状态 200。
- kubernetes：CoreDNS 将根据 Kubernetes 服务和 pod 的 IP 回复 DNS 查询。
- prometheus：是否开启 CoreDNS Metrics 信息接口，如果配置则开启，接口地址为 http://localhost:9153/metrics
- forward：任何不在Kubernetes 集群内的域名查询将被转发到预定义的解析器 (/etc/resolv.conf)。
- cache：启用缓存，30 秒 TTL。
- loop：检测简单的转发循环，如果找到循环则停止 CoreDNS 进程。
- reload：监听 CoreDNS 配置，如果配置发生变化则重新加载配置。
- loadbalance：DNS 负载均衡器，默认 round_robin。

```shell
## 编辑 coredns 配置
$ kubectl edit configmap coredns -n kube-system

apiVersion: v1
data:
Corefile: |
.:53 {
    log            #添加log
    errors
    health {
       lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
       pods insecure
       fallthrough in-addr.arpa ip6.arpa
       ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
} 
```

保存更改后 CoreDNS 会自动重新加载配置信息，不过可能需要等上一两分钟才能将这些更改传播到 CoreDNS Pod。等一段时间后，再次查看 CoreDNS 日志：

```shell
$ for p in $(kubectl get pods --namespace=kube-system -l k8s-app=kube-dns -o name); \
do kubectl logs --namespace=kube-system $p; done

.:53
[INFO] plugin/reload: Running configuration MD5 = 6434d0912b39645ed0707a3569fd69dc
CoreDNS-1.6.7
linux/amd64, go1.13.6, da7f65b
[INFO] Reloading
[INFO] plugin/health: Going into lameduck mode for 5s
[INFO] 127.0.0.1:47278 - 55171 "HINFO IN 4940754309314083739.5160468069505858354. udp 57 false 512" NXDOMAIN qr,rd,ra 57 0.040844011s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete

.:53
[INFO] plugin/reload: Running configuration MD5 = 6434d0912b39645ed0707a3569fd69dc
CoreDNS-1.6.7
linux/amd64, go1.13.6, da7f65b
[INFO] Reloading
[INFO] plugin/health: Going into lameduck mode for 5s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete
[INFO] 127.0.0.1:32896 - 49064 "HINFO IN 1027842207973621585.7098421666386159336. udp 57 false 512" NXDOMAIN qr,rd,ra 57 0.044098742s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete
```

可以看到 CoreDNS 已经重新加载了配置，我们再次进入 dnsuitls Pod 中执行 ping 命令：

```shell
## 进入 DNSutils Pod 命令行
$ kubectl exec -it dnsutils /bin/sh -n kube-system

## 执行 ping 命令
$ ping www.baidu.com

## 退出 dnsutils Pod 命令行
$ exit
```

然后，再次查看 CoreDNS 的日志信息：

```shell
$ for p in $(kubectl get pods --namespace=kube-system -l k8s-app=kube-dns -o name); \
do kubectl logs --namespace=kube-system $p; done

.:53
[INFO] plugin/reload: Running configuration MD5 = 6434d0912b39645ed0707a3569fd69dc
CoreDNS-1.6.7
linux/amd64, go1.13.6, da7f65b
[INFO] Reloading
[INFO] plugin/health: Going into lameduck mode for 5s
[INFO] 127.0.0.1:47278 - 55171 "HINFO IN 4940754309314083739.5160468069505858354. udp 57 false 512" NXDOMAIN qr,rd,ra 57 0.040844011s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete

.:53
[INFO] plugin/reload: Running configuration MD5 = 6434d0912b39645ed0707a3569fd69dc
CoreDNS-1.6.7
linux/amd64, go1.13.6, da7f65b
[INFO] Reloading
[INFO] plugin/health: Going into lameduck mode for 5s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete
[INFO] 127.0.0.1:32896 - 49064 "HINFO IN 1027842207973621585.7098421666386159336. udp 57 false 512" NXDOMAIN qr,rd,ra 57 0.044098742s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete
```

发现和之前没有执行 ping 命令时候一样，没有 DNS 域名解析的日志信息，说明 Pod 执行域名解析时，请求并没有进入 CoreDNS 中。接下来在查看 Pod 中 DNS 配置信息，进而分析问题。

### 查看 Pod 中的 DNS 配置信息

一般 Pod 中的 DNS 策略默认为 ClusterFirst，该参数起到的作用是，优先从 Kubernetes DNS 插件地址读取 DNS 配置。所以，我们正常创建的 Pod 中，DNS 配置 DNS 服务器地址应该指定为 Kubernetes 集群的 DNS 插件 Service 提供的虚拟 IP 地址。

注：其中 DNS 策略（dnsPolicy）支持四种类型：
- Default： 从 DNS 所在节点继承 DNS 配置，即该 Pod 的 DNS 配置与宿主机完全一致。
- ClusterFirst：预先从 Kubenetes 的 DNS 插件中进行 DNS 解析，如果解析不成功，才会使用宿主机的 DNS 进行解析。
- ClusterFirstWithHostNet：Pod 是用 HOST 模式启动的（hostnetwork），用 HOST 模式表示 Pod 中的所有容器，都使用宿主机的 /etc/resolv.conf 配置进行 DNS 解析，但如果使用了 HOST 模式，还继续使用 Kubernetes 的 DNS 服务，那就将 dnsPolicy 设置为 ClusterFirstWithHostNet。
- None：完全忽略 kubernetes 系统提供的 DNS，以 Pod Spec 中 dnsConfig 配置为主。

为了再分析原因，我们接着进入 dnsutils Pod 中，查看 Pod 中 DNS 配置文件 /etc/resolv.conf 配置参数是否正确：

resolv.conf 配置参数说明：
- search： 指明域名查询顺序。
- nameserver： 指定 DNS 服务器的 IP 地址，可以配置多个 nameserver。

```shell
## 进入 dnsutils Pod 内部 sh 命令行
$ kubectl exec -it dnsutils /bin/sh -n kube-system

## 查看 resolv.conf 配置文件
$ cat /etc/resolv.conf

nameserver 10.96.0.10
search kube-system.svc.cluster.local svc.cluster.local cluster.local
options ndots:5

## 退出 DNSutils Pod 命令行
$ exit
```

可以看到 Pod 内部的 resolv.conf 内容，其中 nameserver 指定 DNS 解析服务器 IP 为 “10.96.0.10” ，这个 IP 地址正是本人 Kubernetes 集群 CoreDNS 的 Service “kube-dns” 的 cluterIP，说明当 Pod 内部进行域名解析时，确实是将查询请求发送到 Service “kube-dns” 提供的虚拟 IP 进行域名解析。

那么，既然 Pod 中 DNS 配置文件没问题，且 CoreDNS 也没问题，会不会是 Pod 本身域名解析不正常呢？或者 Service “kube-dns” 是否能够正常转发域名解析请求到 CoreDNS Pod 中？


### 进行观察来定位问题所在

上面怀疑是 Pod 本身解析域名有问题，不能正常解析域名。或者 Pod 没问题，但是请求域名解析时将请求发送到 Service “kube-dns” 后不能正常转发请求到 CoreDNS Pod。 为了验证这两点，我们可以修改 Pod 中的 /etc/resolv.conf 配置来进行测试验证。


修改 resolv.conf 中 DNS 解析请求地址为 阿里云 DNS 服务器地址，然后执行 ping 命令验证是否为 Pod 解析域名是否有问题：

```shell
## 进入 dnsutils Pod 内部 sh 命令行
$ kubectl exec -it dnsutils /bin/sh -n kube-system

## 编辑 /etc/resolv.conf 文件，修改 nameserver 参数为阿里云提供的 DNS 服务器地址
$ vi /etc/resolv.conf

nameserver 223.5.5.5
#nameserver 10.96.0.10
search kube-system.svc.cluster.local svc.cluster.local cluster.local
options ndots:5

## 修改完后再进行 ping 命令测试，看看是否能够解析 www.mydlq.club 网址
$ ping www.mydlq.club

PING www.mydlq.club (140.143.8.181) 56(84) bytes of data.
64 bytes from 140.143.8.181 (140.143.8.181): icmp_seq=1 ttl=128 time=9.70 ms
64 bytes from 140.143.8.181 (140.143.8.181): icmp_seq=2 ttl=128 time=9.21 ms

## 退出 DNSutils Pod 命令行
$ exit
```
上面可也观察到 Pod 中更换 DNS 服务器地址后，域名解析正常，说明 Pod 本身域名解析是没有问题的。

接下来再修改 resolv.conf 中 DNS 解析请求地址为 CoreDNS Pod 的 IP 地址，这样让 Pod 直接连接 CoreDNS Pod 的 IP，而不通过 Service 进行转发，再进行 ping 命令测试，进而判断 Service kube-dns 是否能够正常转发请求到 CoreDNS Pod 的问题：

```shell
## 查看 CoreDNS Pod 的 IP 地址
$ kubectl get pods -n kube-system -o wide | grep coredns

coredns-669f77d7cc-rss5f     1/1     Running   0     10.244.2.155   k8s-node-2-13
coredns-669f77d7cc-rt8l6     1/1     Running   0     10.244.1.163   k8s-node-2-12

## 进入 dnsutils Pod 内部 sh 命令行
$ kubectl exec -it dnsutils /bin/sh -n kube-system

## 编辑 /etc/resolv.conf 文件，修改 nameserver 参数为阿里云提供的 DNS 服务器地址
$ vi /etc/resolv.conf

nameserver 10.244.2.155
nameserver 10.244.1.163
#nameserver 10.96.0.10
search kube-system.svc.cluster.local svc.cluster.local cluster.local
options ndots:5

## 修改完后再进行 ping 命令测试，看看是否能够解析 www.mydlq.club 网址
$ ping www.mydlq.club

PING www.baidu.com (39.156.66.18): 56 data bytes
64 bytes from 39.156.66.18: seq=0 ttl=127 time=6.054 ms
64 bytes from 39.156.66.18: seq=1 ttl=127 time=4.678 ms

## 退出 DNSutils Pod 命令行
$ exit

## 观察 CoreDNS 日志信息，查看有无域名解析相关日志
$ for p in $(kubectl get pods --namespace=kube-system -l k8s-app=kube-dns -o name); \
do kubectl logs --namespace=kube-system $p; done

.:53
[INFO] plugin/reload: Running configuration MD5 = 6434d0912b39645ed0707a3569fd69dc
CoreDNS-1.6.7
linux/amd64, go1.13.6, da7f65b
[INFO] Reloading
[INFO] plugin/health: Going into lameduck mode for 5s
[INFO] 127.0.0.1:47278 - 55171 "HINFO IN 4940754309314083739.5160468069505858354. udp 57 false 512" NXDOMAIN qr,rd,ra 57 0.040844011s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete
[INFO] 10.244.1.162:40261 - 21083 "AAAA IN www.mydlq.club.kube-system.svc.cluster.local. udp 62 false 512" NXDOMAIN qr,aa,rd 155 0.000398875s
[INFO] 10.244.1.162:40261 - 20812 "A IN www.mydlq.club.kube-system.svc.cluster.local. udp 62 false 512" NXDOMAIN qr,aa,rd 155 0.000505793s
[INFO] 10.244.1.162:55066 - 53460 "AAAA IN www.mydlq.club.svc.cluster.local. udp 50 false 512" NXDOMAIN qr,aa,rd 143 0.000215384s
[INFO] 10.244.1.162:55066 - 53239 "A IN www.mydlq.club.svc.cluster.local. udp 50 false 512" NXDOMAIN qr,aa,rd 143 0.000267642s

.:53
[INFO] plugin/reload: Running configuration MD5 = 6434d0912b39645ed0707a3569fd69dc
CoreDNS-1.6.7
linux/amd64, go1.13.6, da7f65b
[INFO] Reloading
[INFO] plugin/health: Going into lameduck mode for 5s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete
[INFO] 127.0.0.1:32896 - 49064 "HINFO IN 1027842207973621585.7098421666386159336. udp 57 false 512" NXDOMAIN qr,rd,ra 57 0.044098742s
[INFO] plugin/reload: Running configuration MD5 = a4809ab99f6713c362194263016e6fac
[INFO] Reloading complete
[INFO] 10.244.1.162:40261 - 21083 "AAAA IN www.mydlq.club.kube-system.svc.cluster.local. udp 62 false 512" NXDOMAIN qr,aa,rd 155 0.000217299s
[INFO] 10.244.1.162:40261 - 20812 "A IN www.mydlq.club.kube-system.svc.cluster.local. udp 62 false 512" NXDOMAIN qr,aa,rd 155 0.000264552s
[INFO] 10.244.1.162:55066 - 53460 "AAAA IN www.mydlq.club.svc.cluster.local. udp 50 false 512" NXDOMAIN qr,aa,rd 143 0.000144795s
[INFO] 10.244.1.162:55066 - 53239 "A IN www.mydlq.club.svc.cluster.local. udp 50 false 512" N
```

经过上面两个测试，已经可以得知，如果 Pod DNS 配置中直接修改 DNS 服务器地址为 CoreDNS Pod 的 IP 地址，DNS 解析确实没有问题，能够正常解析。不过，正常的情况下 Pod 中 DNS 配置的服务器地址一般是 CoreDNS 的 Service 地址，不直接绑定 Pod IP（因为 Pod 每次重启 IP 都会发生变化）。 所以问题找到了，正是在 Pod 向 CoreDNS 的 Service “kube-dns” 进行域名解析请求转发时，出现了问题，一般 Service 的问题都跟 Kube-proxy 组件有关，接下来观察该组件是否存在问题。

### 分析 Kube-Proxy 是否存在问题

观察 Kube-proxy 的日志，查看是否存在问题：

```shell
## 查看 kube-proxy Pod 列表
$ kubectl get pods -n kube-system | grep kube-proxy

kube-proxy-6kdj2          1/1     Running   3          9h
kube-proxy-lw2q6          1/1     Running   3          9h
kube-proxy-mftlt          1/1     Running   3          9h

## 选择一个 kube-proxy Pod，查看最后 5 条日志内容
$ kubectl logs kube-proxy-6kdj2 --tail=5  -n kube-system

E0326 15:20:23.159364  1 proxier.go:1950] Failed to list IPVS destinations, error: parseIP Error ip=[10 96 0 10 0 0 0 0 0 0 0 0 0 0 0 0]
E0326 15:20:23.159388  1 proxier.go:1192] Failed to sync endpoint for service: 10.8.0.10:53/UPD, err: parseIP Error ip=[10 96 0 16 0 0 0 0 0 0 0 0 0 0 0 0]
E0326 15:20:23.159479  1 proxier.go:1950] Failed to list IPVS destinations, error: parseIP Error ip=[10 96 0 10 0 0 0 0 0 0 0 0 0 0 0 0]
E0326 15:20:23.159501  1 proxier.go:1192] Failed to sync endpoint for service: 10.8.0.10:53/TCP, err: parseIP Error ip=[10 96 0 16 0 0 0 0 0 0 0 0 0 0 0 0]
E0326 15:20:23.159595  1 proxier.go:1950] Failed to list IPVS destinations, error: parseIP
```

通过 kube-proxy Pod 的日志可以看到，里面有很多 Error 级别的日志信息，根据关键字 IPVS、parseIP Error 可知，可能是由于 IPVS 模块对 IP 进行格式化导致出现问题。

因为这个问题是升级到 Kubernetes 1.18 版本才出现的，所以去 Kubernetes GitHub 查看相关 issues，发现有人在升级 Kubernetes 版本到 1.18 后，也遇见了相同的问题，经过 issue 中 Kubernetes 维护人员讨论，分析出原因可能为新版 Kubernetes 使用的 IPVS 模块是比较新的，需要系统内核版本支持，本人使用的是 CentOS 系统，内核版本为 3.10，里面的 IPVS 模块比较老旧，缺少新版 Kubernetes IPVS 所需的依赖。

根据该 issue 讨论结果，解决该问题的办法是，更新内核为新的版本。

注：该 issues 地址为：https://github.com/kubernetes/ ... 89520

## 解决问题

### 升级系统内核版本

升级 Kubernetes 集群各个节点的 CentOS 系统内核版本：
```shell
## 载入公钥
$ rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org

## 安装 ELRepo 最新版本
$ yum install -y https://www.elrepo.org/elrepo-release-7.el7.elrepo.noarch.rpm

## 列出可以使用的 kernel 包版本
$ yum list available --disablerepo=* --enablerepo=elrepo-kernel

## 安装指定的 kernel 版本：
$ yum install -y kernel-lt-4.4.218-1.el7.elrepo --enablerepo=elrepo-kernel

## 查看系统可用内核
$ cat /boot/grub2/grub.cfg | grep menuentry

menuentry 'CentOS Linux (3.10.0-1062.el7.x86_64) 7 (Core)' --class centos （略）
menuentry 'CentOS Linux (4.4.218-1.el7.elrepo.x86_64) 7 (Core)' --class centos ...（略）

## 设置开机从新内核启动
$ grub2-set-default "CentOS Linux (4.4.218-1.el7.elrepo.x86_64) 7 (Core)"

## 查看内核启动项
$ grub2-editenv list
saved_entry=CentOS Linux (4.4.218-1.el7.elrepo.x86_64) 7 (Core)
```
重启系统使内核生效：
```shell
$ reboot
```

启动完成查看内核版本是否更新：

```shell
$ uname -r

4.4.218-1.el7.elrepo.x86_64
```

### 测试 Pod 中 DNS 是否能够正常解析

进入 Pod 内部使用 ping 命令测试 DNS 是否能正常解析：
```shell
## 进入 dnsutils Pod 内部 sh 命令行
$ kubectl exec -it dnsutils /bin/sh -n kube-system

## Ping 集群外部，例如这里 ping 一下百度
$ ping www.baidu.com
64 bytes from 39.156.66.14 (39.156.66.14): icmp_seq=1 ttl=127 time=7.20 ms
64 bytes from 39.156.66.14 (39.156.66.14): icmp_seq=2 ttl=127 time=6.60 ms
64 bytes from 39.156.66.14 (39.156.66.14): icmp_seq=3 ttl=127 time=6.38 ms

## Ping 集群内部 kube-api 的 Service 地址
$ ping kubernetes.default
64 bytes from kubernetes.default.svc.cluster.local (10.96.0.1): icmp_seq=1 ttl=64 time=0.051 ms
64 bytes from kubernetes.default.svc.cluster.local (10.96.0.1): icmp_seq=2 ttl=64 time=0.051 ms
64 bytes from kubernetes.default.svc.cluster.local (10.96.0.1): icmp_seq=3 ttl=64 time=0.064 ms
```

可以看到 Pod 中的域名解析已经恢复正常。