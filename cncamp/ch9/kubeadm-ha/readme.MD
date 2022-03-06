## First master

### Reset kubeadm

```sh
kubeadm reset
```

### Run the following command on 1st master

```sh
echo "10.95.102.110 master.cncamp.io" >> /etc/hosts
echo "10.95.102.111 cn-mater1" >> /etc/hosts
echo "10.95.102.112 cn-mater2" >> /etc/hosts
echo "10.95.102.113 cn-mater3" >> /etc/hosts
```

### Init on first master

```sh
kubeadm init --config kubeadm-conf.yaml --upload-certs=true
```

### Copy kubeconfig

```sh
$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

### Untaint master

```sh
$ kubectl taint nodes --all node-role.kubernetes.io/master-
```

### Install calico cni plugin

https://docs.projectcalico.org/getting-started/kubernetes/quickstart

```sh
$ kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml
$ kubectl create -f https://docs.projectcalico.org/manifests/custom-resources.yaml
```

### Second master

### Reset kubeadm

```sh
kubeadm reset
```

### Run the following command on 2nd master

```sh
echo "10.95.102.110 master.cncamp.io" >> /etc/hosts
echo "10.95.102.111 cn-mater1" >> /etc/hosts
echo "10.95.102.112 cn-mater2" >> /etc/hosts
echo "10.95.102.113 cn-mater3" >> /etc/hosts
```

### Join the 2nd master

```sh
kubeadm join master.cncamp.io:6443 --token 2m54gi.0v6s155lujz3ruyk \
--discovery-token-ca-cert-hash sha256:22aee251c9a2ea7442862555ed6fce301f44548069e46bfad9210aaffa4f1cef \
--control-plane --certificate-key 8a4d46120ad3eaa63b6e3395a59b4ef8c84eb90c2ef7aa375e707086959fa72b \
--apiserver-advertise-address=10.95.102.111
```

### Check etcd

```sh
ETCDCTL_API=3
etcdctl --endpoints https://10.95.102.111:2379 \
--cacert /etc/kubernetes/pki/etcd/ca.crt \
--cert /etc/kubernetes/pki/etcd/server.crt \
--key /etc/kubernetes/pki/etcd/server.key \
member list
```
