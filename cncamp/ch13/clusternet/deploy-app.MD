### create deployment in host cluster
```shell
kubectl clusternet apply -f deployment.yaml
```
### where the object being stored, as manifest object in clusternet-reserved
```shell
k get manifest -A
NAMESPACE             NAME                       AGE
clusternet-reserved   deployments-foo-my-nginx   9h
clusternet-reserved   namespaces-foo             9h
```
### link the deployment and target cluster
```shell
kubectl apply -f localization.yaml
kubectl apply -f subscription.yaml

```
### check generated obj
```
k get base -A
NAMESPACE          NAME       AGE
clusternet-bmqv4   app-demo   9h

k get description -A
NAMESPACE          NAME               DEPLOYER   STATUS    AGE
clusternet-bmqv4   app-demo-generic   Generic    Success   9h
```

```
etcdctl --endpoints https://127.0.0.1:2379 \
--cacert /etc/kubernetes/pki/etcd/ca.crt \
--cert /etc/kubernetes/pki/etcd/server.crt \
--key /etc/kubernetes/pki/etcd/server.key \
```