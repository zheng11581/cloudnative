## 高级方案
更高级的方案已经有人写好了，项目地址是：https://github.com/ysde/grafana-backup-tool

### 该备份工具支持以下几种配置：

目录

Dashboard

数据源

Grafana 告警频道（Alert Channel）

组织（Organization）

用户（User）

使用方法很简单，跑个容器就好了嘛，不过作者提供的 Dockerfile 我不是很满意，自己修改了点内容：

```dockerfile
FROM alpine:latest

LABEL maintainer="grafana-backup-tool Docker Maintainers https://fuckcloudnative.io"

ENV ARCHIVE_FILE ""

RUN echo "@edge http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories; \
apk --no-cache add python3 py3-pip py3-cffi py3-cryptography ca-certificates bash git; \
git clone https://github.com/ysde/grafana-backup-tool /opt/grafana-backup-tool; \
cd /opt/grafana-backup-tool; \
pip3 --no-cache-dir install .; \
chown -R 1337:1337 /opt/grafana-backup-tool

WORKDIR /opt/grafana-backup-tool

USER 1337
```

只有 Dockerfile 不行，还得通过 CI/CD 自动构建并推送到 docker.io。不要问我用什么，当然是白嫖 GitHub Action，workflow 内容如下：
```yaml
#=================================================
# https://github.com/yangchuansheng/docker-image
# Description: Build and push grafana-backup-tool Docker image
# Lisence: MIT
# Author: Ryan
# Blog: https://fuckcloudnative.io
#=================================================

name: Build and push grafana-backup-tool Docker image

# Controls when the action will run. Triggers the workflow on push or pull request
# events but only for the master branch
on:
   push:
      branches: [ master ]
      paths:
         - 'grafana-backup-tool/Dockerfile'
         - '.github/workflows/grafana-backup-tool.yml'
   pull_request:
      branches: [ master ]
      paths:
         - 'grafana-backup-tool/Dockerfile'
      #watch:
      #types: started

# A workflow run is made up of one or more jobs that can run sequentially or in parallel
jobs:
   # This workflow contains a single job called "build"
   build:
      # The type of runner that the job will run on
      runs-on: ubuntu-latest

      # Steps represent a sequence of tasks that will be executed as part of the job
      steps:
         # Checks-out your repository under $GITHUB_WORKSPACE, so your job can access it
         - uses: actions/checkout@v2

         - name: Set up QEMU
           uses: docker/setup-qemu-action@v1

         - name: Set up Docker Buildx
           uses: docker/setup-buildx-action@v1

         - name: Login to DockerHub
           uses: docker/login-action@v1
           with:
              username: ${{ secrets.DOCKER_USERNAME }}
              password: ${{ secrets.DOCKER_PASSWORD }}

         - name: Login to GitHub Package Registry
           env:
              username: ${{ github.repository_owner }}
              password: ${{ secrets.GHCR_TOKEN }}
           run: echo ${{ env.password }} | docker login ghcr.io -u ${{ env.username }} --password-stdin

         # Runs a single command using the runners shell
         - name: Build and push Docker images to docker.io and ghcr.io
           uses: docker/build-push-action@v2
           with:
              file: 'grafana-backup-tool/Dockerfile'
              platforms: linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64,linux/ppc64le,linux/s390x
              context: grafana-backup-tool
              push: true
              tags: |
                 yangchuansheng/grafana-backup-tool:latest
                 ghcr.io/yangchuansheng/grafana-backup-tool:latest

              #- name: Update repo description
              #uses: peter-evans/dockerhub-description@v2
              #env:
              #DOCKERHUB_USERNAME: ${{ secrets.DOCKER_USERNAME }}
              #DOCKERHUB_PASSWORD: ${{ secrets.DOCKER_PASSWORD }}
              #DOCKERHUB_REPOSITORY: yangchuansheng/grafana-backup-tool
              #README_FILEPATH: grafana-backup-tool/readme.md
```

你可以直接关注我的仓库：
https://github.com/yangchuansheng/docker-image

构建好镜像后，就可以直接运行容器来进行备份和恢复操作了。如果你想在集群内操作，可以通过 Deployment 或 Job 来实现；如果你想在本地或 k8s 集群外操作，可以选择 docker run，我不反对，你也可以选择 docker-compose，这都没问题。但我要告诉你一个更骚的办法，可以骚到让你无法自拔。

首先需要在本地或集群外安装 Podman，如果操作系统是 Win10，可以考虑通过 WSL 来安装；如果操作系统是 Linux，那就不用说了；如果操作系统是 MacOS，请参考我的上篇文章：????在 macOS 中使用 Podman。

装好了 Podman 之后，就可以进行骚操作了，请睁大眼睛。

先编写一个 Deployment 配置清单（什么？Deployment？是的，你没听错）：

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
   name: grafana-backup
   labels:
      app: grafana-backup
spec:
   replicas: 1
   selector:
      matchLabels:
         app: grafana-backup
   template:
      metadata:
         labels:
            app: grafana-backup
      spec:
         containers:
            - name: grafana-backup
              image: yangchuansheng/grafana-backup-tool:latest
              imagePullPolicy: IfNotPresent
              command:
                 - "/bin/bash"
              tty: true
              stdin: true
              env:
                 - name: GRAFANA_TOKEN
                   value: "eyJrIjoiOFFiRkNxbkJodE01M0lkM05QUXplSmdnZUxVVnN0TGQiLCJuIjoiZ3JhZmFuYV9iYWNrdXAiLCJpZCI6MX0="
                 - name: GRAFANA_URL
                   value: "http://kube-prometheus-grafana.cpaas-system"
                 - name: GRAFANA_ADMIN_ACCOUNT
                   value: "admin"
                 - name: GRAFANA_ADMIN_PASSWORD
                   value: "admin"
                 - name: VERIFY_SSL
                   value: "False"
              volumeMounts:
                 - mountPath: /opt/grafana-backup-tool
                   name: data
         volumes:
            - name: data
              hostPath:
                 path: /mnt/manifest/grafana/backup
```

这里面的环境变量根据自己的实际情况修改，一定不要照抄我的！

不要一脸懵逼，我先来解释一下为什么要准备这个 Deployment 配置清单，因为 Podman 可以直接通过这个配置清单运行容器，命令如下：
```shell
$ podman play kube grafana-backup-deployment.yaml
```

细心的读者应该能发现上面的配置清单有点奇怪，Dockerfile 也有点奇怪。Dockerfile 中没有写 CMD 或 ENTRYPOINT，Deployment 中直接将启动命令设置为 bash，这是因为在我之前测试的过程中发现该镜像启动的容器有点问题，它会陷入一个循环，备份完了之后又会继续备份，不断重复，导致备份目录下生成了一坨压缩包。目前还没找到比较好的解决办法，只能将容器的启动命令设置为 bash，等容器运行后再进入容器进行备份操作：
```shell
$ podman pod ls
POD ID        NAME                  STATUS   CREATED        # OF CONTAINERS  INFRA ID
728aec216d66  grafana-backup-pod-0  Running  3 minutes ago  2                92aa0824fe7d

$ podman ps
CONTAINER ID  IMAGE                                      COMMAND    CREATED        STATUS            PORTS   NAMES
b523fa8e4819  yangchuansheng/grafana-backup-tool:latest  /bin/bash  3 minutes ago  Up 3 minutes ago          grafana-backup-pod-0-grafana-backup
92aa0824fe7d  k8s.gcr.io/pause:3.2                                  3 minutes ago  Up 3 minutes ago          728aec216d66-infra

$ podman exec -it grafana-backup-pod-0-grafana-backup bash
bash-5.0$ grafana-backup save
...
...
########################################

backup folders at: _OUTPUT_/folders/202012111556
backup datasources at: _OUTPUT_/datasources/202012111556
backup dashboards at: _OUTPUT_/dashboards/202012111556
backup alert_channels at: _OUTPUT_/alert_channels/202012111556
backup organizations at: _OUTPUT_/organizations/202012111556
backup users at: _OUTPUT_/users/202012111556

created archive at: _OUTPUT_/202012111556.tar.gz
```


默认情况下会备份所有的组件，你也可以指定备份的组件：
```shell
$ grafana-backup save --components=<folders,dashboards,datasources,alert-channels,organizations,users>
```

比如，我只想备份 Dashboards 和 Folders：
```shell
$ grafana-backup save --components=folders,dashboards
```

当然，你也可以全部备份，恢复的时候再选择自己想恢复的组件：
```shell
$ grafana-backup restore --components=folders,dashboards
```

至此，再也不用怕 Dashboard 被改掉或删除啦。

最后提醒一下，Prometheus Operator 项目中的 Grafana 通过 Provisioning[4] 的方式 预导入 了一些默认的 Dashboards，这本来没有什么问题，但 grafana-backup-tool 工具无法忽略跳过已经存在的配置，如果恢复的过程中遇到已经存在的配置，会直接报错退出。本来这也很好解决，一般情况下到 Grafana Web 界面中删除所有的 Dashboard 就好了，但通过 Provisioning 导入的 Dashboard 是无法删除的，这就很尴尬了。

在作者修复这个 bug 之前，要想解决这个问题，有两个办法：

第一个办法是在恢复之前将 Grafana Deployment 中关于  Provisioning 的配置全部删除，就是这些配置：
```shell
        volumeMounts:
        - mountPath: /etc/grafana/provisioning/datasources
          name: grafana-datasources
          readOnly: false
        - mountPath: /etc/grafana/provisioning/dashboards
          name: grafana-dashboards
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/apiserver
          name: grafana-dashboard-apiserver
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/cluster-total
          name: grafana-dashboard-cluster-total
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/controller-manager
          name: grafana-dashboard-controller-manager
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/k8s-resources-cluster
          name: grafana-dashboard-k8s-resources-cluster
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/k8s-resources-namespace
          name: grafana-dashboard-k8s-resources-namespace
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/k8s-resources-node
          name: grafana-dashboard-k8s-resources-node
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/k8s-resources-pod
          name: grafana-dashboard-k8s-resources-pod
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/k8s-resources-workload
          name: grafana-dashboard-k8s-resources-workload
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/k8s-resources-workloads-namespace
          name: grafana-dashboard-k8s-resources-workloads-namespace
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/kubelet
          name: grafana-dashboard-kubelet
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/namespace-by-pod
          name: grafana-dashboard-namespace-by-pod
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/namespace-by-workload
          name: grafana-dashboard-namespace-by-workload
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/node-cluster-rsrc-use
          name: grafana-dashboard-node-cluster-rsrc-use
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/node-rsrc-use
          name: grafana-dashboard-node-rsrc-use
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/nodes
          name: grafana-dashboard-nodes
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/persistentvolumesusage
          name: grafana-dashboard-persistentvolumesusage
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/pod-total
          name: grafana-dashboard-pod-total
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/prometheus-remote-write
          name: grafana-dashboard-prometheus-remote-write
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/prometheus
          name: grafana-dashboard-prometheus
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/proxy
          name: grafana-dashboard-proxy
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/scheduler
          name: grafana-dashboard-scheduler
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/statefulset
          name: grafana-dashboard-statefulset
          readOnly: false
        - mountPath: /grafana-dashboard-definitions/0/workload-total
          name: grafana-dashboard-workload-total
          readOnly: false
...
...
      volumes:
      - name: grafana-datasources
        secret:
          secretName: grafana-datasources
      - configMap:
          name: grafana-dashboards
        name: grafana-dashboards
      - configMap:
          name: grafana-dashboard-apiserver
        name: grafana-dashboard-apiserver
      - configMap:
          name: grafana-dashboard-cluster-total
        name: grafana-dashboard-cluster-total
      - configMap:
          name: grafana-dashboard-controller-manager
        name: grafana-dashboard-controller-manager
      - configMap:
          name: grafana-dashboard-k8s-resources-cluster
        name: grafana-dashboard-k8s-resources-cluster
      - configMap:
          name: grafana-dashboard-k8s-resources-namespace
        name: grafana-dashboard-k8s-resources-namespace
      - configMap:
          name: grafana-dashboard-k8s-resources-node
        name: grafana-dashboard-k8s-resources-node
      - configMap:
          name: grafana-dashboard-k8s-resources-pod
        name: grafana-dashboard-k8s-resources-pod
      - configMap:
          name: grafana-dashboard-k8s-resources-workload
        name: grafana-dashboard-k8s-resources-workload
      - configMap:
          name: grafana-dashboard-k8s-resources-workloads-namespace
        name: grafana-dashboard-k8s-resources-workloads-namespace
      - configMap:
          name: grafana-dashboard-kubelet
        name: grafana-dashboard-kubelet
      - configMap:
          name: grafana-dashboard-namespace-by-pod
        name: grafana-dashboard-namespace-by-pod
      - configMap:
          name: grafana-dashboard-namespace-by-workload
        name: grafana-dashboard-namespace-by-workload
      - configMap:
          name: grafana-dashboard-node-cluster-rsrc-use
        name: grafana-dashboard-node-cluster-rsrc-use
      - configMap:
          name: grafana-dashboard-node-rsrc-use
        name: grafana-dashboard-node-rsrc-use
      - configMap:
          name: grafana-dashboard-nodes
        name: grafana-dashboard-nodes
      - configMap:
          name: grafana-dashboard-persistentvolumesusage
        name: grafana-dashboard-persistentvolumesusage
      - configMap:
          name: grafana-dashboard-pod-total
        name: grafana-dashboard-pod-total
      - configMap:
          name: grafana-dashboard-prometheus-remote-write
        name: grafana-dashboard-prometheus-remote-write
      - configMap:
          name: grafana-dashboard-prometheus
        name: grafana-dashboard-prometheus
      - configMap:
          name: grafana-dashboard-proxy
        name: grafana-dashboard-proxy
      - configMap:
          name: grafana-dashboard-scheduler
        name: grafana-dashboard-scheduler
      - configMap:
          name: grafana-dashboard-statefulset
        name: grafana-dashboard-statefulset
      - configMap:
          name: grafana-dashboard-workload-total
        name: grafana-dashboard-workload-total
```

第二个办法就是删除 Prometheus Operator 自带的 Grafana，自己通过 Helm 或者 manifest 部署不使用 Provisioning 的 Grafana。
