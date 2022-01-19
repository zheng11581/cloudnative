## 网络配置
```shell
vim /etc/netplan/00-installer-config.yaml
    # This is the network config written by 'subiquity'
    network:
      ethernets:
        ens18: # device
          dhcp4: false # static
          addresses: [10.95.102.111/24]
          optional: true
          gateway4: 10.95.102.1
          nameservers:
                  addresses: [114.114.114.114]
      version: 2

netplan apply

```

## 更改镜像源
```shell
cp /etc/apt/sources.list /etc/apt/sources.list.bak
vim /etc/apt/sources.list
    deb http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse
    deb-src http://mirrors.aliyun.com/ubuntu/ focal main restricted universe multiverse
    
    deb http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse
    deb-src http://mirrors.aliyun.com/ubuntu/ focal-security main restricted universe multiverse
    
    deb http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse
    deb-src http://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse
    
    deb http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse
    deb-src http://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse
    
    deb http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
    deb-src http://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse
```