### Install kind

### Create host cluster
```sh
./create_cluster.sh
```

### Compile clusterctl
```shell
mkdir $GOPATH/src/k8s.io/kubernetes-sigs
git clone git@github.com:kubernetes-sigs/cluster-api.git
cd cluster-api
make clusterctl
chmod 755 bin/clusterctl && mv bin/clusterctl /usr/local/bin/
```

### Generate cluster specs

```sh
cd cluster-api
./init_docker_provider.sh
./generate_workload_cluster.sh
```

### replace image repository
```shell
vim capi-quickstart.yaml
...
imageRepository: registry.aliyuncs.com/google_container
...
```

### apply cluster spec
```shell
kubectl apply -f capi-quickstart.yaml
```

### Check

```sh
$ docker ps|grep control-plane
b107b11771e5        kindest/haproxy:v20210715-a6da3463   "haproxy -sf 7 -W -d…"   4 minutes ago       Up 4 minutes        40295/tcp, 0.0.0.0:40295->6443/tcp     capi-quickstart-lb

$ clusterctl get kubeconfig capi-quickstart > capi-quickstart.kubeconfig

$ kubectl get no --kubeconfig capi-quickstart.kubeconfig --server https://127.0.0.1:40295
NAME                                    STATUS     ROLES                  AGE     VERSION
capi-quickstart-control-plane-6slwd     NotReady   control-plane,master   4m19s   v1.22.0
capi-quickstart-md-0-765cf784c5-6klwr   NotReady   <none>                 3m41s   v1.22.0
```

```sh
kubectl get cluster
kubectl get machineset
```
