## Question 1 | Contexts (1%)

### Describe

```text
You have access to multiple clusters from your main terminal through kubectl contexts. 
1. Write all those context names into /opt/course/1/contexts.
2. Next write a command to display the current context into /opt/course/1/context_default_kubectl.sh, the command should use kubectl.
3. Finally write a second command doing the same thing into /opt/course/1/context_default_no_kubectl.sh, but without the use of kubectl.
```

### Result

```shell
# 1
# k config get-contexts -o name > /opt/course/1/contexts

# 2
# k config current-context > /opt/course/1/context_default_kubectl.sh

# 3
# cat ~/.kube/config |grep current |sed -e "s/current-context: //" > /opt/course/1/context_default_no_kubectl.sh
```

## Question 2 | Schedule Pod on Master Node (3%)

### Describe
```text
1. Create a single Pod of image httpd:2.4.41-alpine in Namespace default. 
2. The Pod should be named pod1 and the container should be named pod1-container. 
3. This Pod should only be scheduled on a master node, do not add new labels any nodes.
```

### Result

```shell
# k get node # find master node

# k describe node cluster1-master1 | grep Taint -A1 # get master node taints

# k get node cluster1-master1 --show-labels # get master node labels

# kubectl run pod1 --image=httpd:2.4.41-alpine --namespace=default $do > 2.yaml

```
Final yaml is:

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: pod1
  name: pod1
  namespace: default
spec:
  tolerantions:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
  - effect: NoSchedule
    key: node-role.kubernetes.io/control-plane
  nodeSelectors:
    node-role.kubernetes.io/control-plane: ""
  # Or you can just specify a nodeName
  # nodeName: <masterNodeName>
  containers:
  - image: httpd:2.4.41-alpine
    name: pod1
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}
```

```shell
# k -f 2.yaml create
```

## Question 3 | Scale down StatefulSet (3%)

### Describe

```text
1. There are two Pods named o3db-* in Namespace project-c13. 
3. C13 management asked you to scale the Pods down to one replica to save resources.
```

### Result

```shell
# k get deploy,sts,ds -n project-c13 |grep o3db 
# It is a statfulset

# k get edit -n project-c13 o3db
# Replace replicas OR
# k -n project-c13 scale sts o3db --replicas 1


```

## Question 4 | Pod Ready if Service is reachable

### Describe

```text
Do the following in Namespace default. 
1. Create a single Pod named ready-if-service-ready of image nginx:1.16.1-alpine. 
2. Configure a LivenessProbe which simply runs true. 
3. Also configure a ReadinessProbe which does check if the url http://service-am-i-ready:80 is reachable, you can use wget -T2 -O- http://service-am-i-ready:80 for this. 
4. Start the Pod and confirm it isn't ready because of the ReadinessProbe.

5. Create a second Pod named am-i-ready of image nginx:1.16.1-alpine with label id: cross-server-ready. 
6. The already existing Service service-am-i-ready should now have that second Pod as endpoint.

Now the first Pod should be in ready state, confirm that.
```

### Result

```shell
# k run ready-if-service-ready --image=nginx:1.16.1-alpine $do > 4-pod1.yaml


# k run am-i-ready --image=nginx:1.16.1-alpine --labels="id=cross-server-ready" $do > 4-pod2.yaml
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: ready-if-service-ready
  name: ready-if-service-ready
spec:
  containers:
  - image: nginx:1.16.1-alpine
    name: ready-if-service-ready
    livenessProbe:
      exec:
        command: 
        - "true"
    readinessProbe:
      exec:
        command:
        - sh
        - -c
        - "wget -T2 -O- http://service-am-i-ready:80"
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    id: cross-server-ready
  name: am-i-ready
spec:
  containers:
  - image: nginx:1.16.1-alpine
    name: am-i-ready
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}
```

```shell
# k create -f 4-pod1.yaml
# k get pod ready-if-service-ready
# k create -f 4-pod2.yaml
# k get ep service-am-i-ready 
# k get pod ready-if-service-ready
```

## Question 5 | Kubectl sorting (1%)

### Describe

```text
There are various Pods in all namespaces. 
1. Write a command into /opt/course/5/find_pods.sh which lists all Pods sorted by their AGE (metadata.creationTimestamp).

2. Write a second command into /opt/course/5/find_pods_uid.sh which lists all Pods sorted by field metadata.uid. Use kubectl sorting for both commands.
```

### Result

```shell
# kubectl get pod -A --sort-by=.metadata.creationTimestamp >  /opt/course/5/find_pods.sh
# kubectl get pod -A --sort-by=.metadata.uid > /opt/course/5/find_pods_uid.sh
```

## Question 6 | Storage, PV, PVC, Pod volume (8%)

### Describe

```text
1. Create a new PersistentVolume named safari-pv. It should have a capacity of 2Gi, accessMode ReadWriteOnce, hostPath /Volumes/Data and no storageClassName defined.

2. Next create a new PersistentVolumeClaim in Namespace project-tiger named safari-pvc . It should request 2Gi storage, accessMode ReadWriteOnce and should not define a storageClassName. The PVC should bound to the PV correctly.

3. Finally create a new Deployment safari in Namespace project-tiger which mounts that volume at /tmp/safari-data. The Pods of that Deployment should be of image httpd:2.4.41-alpine.
```

### Result

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: safari-pv
spec:
  capacity:
    storage: 2Gi
  accessModes:
  - ReadWriteOnce
  hostPath:
    path: "/Volumes/Data"
```

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: safari-pvc
spec:
  accessModes:
- ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    app: safari
  name: safari
  namespace: project-tiger
spec:
  replicas: 1
  selector:
    matchLabels:
      app: safari
  strategy: {}
  template:
    metadata:
      labels:
        app: safari
    spec:
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: safari-pvc
      containers:
      - name: website
        image: httpd:2.4.41-alpine
        ports:
        - containerPort: 80
          name: "http-server"
        volumeMounts:
        - mountPath: "/tmp/safari-data"
          name: data
```

## Question 7 | Node and Pod Resource Usage (1%)

### Describe

```text
The metrics-server has been installed in the cluster. Your college would like to know the kubectl commands to:
1. show Nodes resource usage
2. show Pods and their containers resource usage
Please write the commands into /opt/course/7/node.sh and /opt/course/7/pod.sh.
```

### Result

```shell
# k top node > /opt/course/7/node.sh
# k top pod --containers=true > /opt/course/7/pod.sh
```



## Question 8 | Get Master Information (2%)

### Describe

```text
Ssh into the master node with ssh cluster1-master1. Check how the master components kubelet, kube-apiserver, kube-scheduler, kube-controller-manager and etcd are started/installed on the master node. Also find out the name of the DNS application and how it's started/installed on the master node.

Write your findings into file /opt/course/8/master-components.txt. The file should be structured like:

# /opt/course/8/master-components.txt
kubelet: [TYPE]
kube-apiserver: [TYPE]
kube-scheduler: [TYPE]
kube-controller-manager: [TYPE]
etcd: [TYPE]
dns: [TYPE] [NAME]
Choices of [TYPE] are: not-installed, process, static-pod, pod
```

### Result

```shell
# ssh cluster1-master1

# ps aux | grep kubelet

# find /etc/systemd/system/ | grep kube

# find /etc/systemd/system/ |grep etcd

# find /etc/kubernetes/manifests/

# kubectl -n kube-system get pod -o wide | grep master1

# vim /opt/course/8/master-components.txt
kubelet: [process]
kube-apiserver: [static-pod]
kube-scheduler: [static-pod]
kube-controller-manager: [static-pod]
etcd: [static-pod]
dns: [pod] [coredns]
```


## Question 9 | Kill Scheduler, Manual Scheduling (5%)

### Describe

1. Ssh into the master node with `ssh cluster2-master1`. Temporarily stop the kube-scheduler, this means in a way that you can start it again afterwards.

2. Create a single Pod named `manual-schedule` of image `httpd:2.4-alpine`, confirm its created but not scheduled on any node.

3. Now you're the scheduler and have all its power, manually schedule that Pod on node cluster2-master1. Make sure it's running.

4. Start the kube-scheduler again and confirm its running correctly by creating a second Pod named `manual-schedule2` of image `httpd:2.4-alpine` and check if it's running on cluster2-worker1.

### Result

```shell
# 1
# ssh cluster2-master1
# cd /etc/kunernetes/manifests/
# mv kube-scheduler.yaml ..
# kubectl -n kube-system get pod | grep schedule

# 2 
# k run manual-schedule --image=httpd:2.4-alpine $do > 9.yaml
# k get pod -o |grep 
```

## Question 1 | Contexts (1%)

### Describe

### Result