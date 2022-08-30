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

1. Create a single Pod of image `httpd:2.4.41-alpine` in Namespace default. 
2. The Pod should be named `pod1` and the container should be named `pod1-container`. 
3. This Pod should only be scheduled on a master node, do not add new labels any nodes.


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


1. There are two Pods named `o3db-*` in Namespace `project-c13`. 
3. C13 management asked you to scale the Pods down to one replica to save resources.

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


Do the following in Namespace `default`. 
1. Create a single Pod named `ready-if-service-ready` of image `nginx:1.16.1-alpine`. 
2. Configure a LivenessProbe which simply runs true. 
3. Also configure a ReadinessProbe which does check if the url http://service-am-i-ready:80 is reachable, you can use `wget -T2 -O- http://service-am-i-ready:80` for this. 
4. Start the Pod and confirm it isn't ready because of the ReadinessProbe.
5. Create a second Pod named `am-i-ready` of image `nginx:1.16.1-alpine` with label `id: cross-server-ready`. 
6. The already existing Service `service-am-i-ready` should now have that second Pod as endpoint.

Now the first Pod should be in ready state, confirm that.


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


There are various Pods in all namespaces. 
1. Write a command into `/opt/course/5/find_pods.sh` which lists all Pods sorted by their AGE (`metadata.creationTimestamp`).
2. Write a second command into `/opt/course/5/find_pods_uid.sh` which lists all Pods sorted by field `metadata.uid`. Use kubectl sorting for both commands.

### Result

```shell
# kubectl get pod -A --sort-by=.metadata.creationTimestamp >  /opt/course/5/find_pods.sh
# kubectl get pod -A --sort-by=.metadata.uid > /opt/course/5/find_pods_uid.sh
```

## Question 6 | Storage, PV, PVC, Pod volume (8%)

### Describe

1. Create a new PersistentVolume named `safari-pv`. It should have a capacity of 2Gi, accessMode ReadWriteOnce, hostPath `/Volumes/Data` and no storageClassName defined.

2. Next create a new PersistentVolumeClaim in Namespace `project-tiger` named `safari-pvc` . It should request 2Gi storage, accessMode ReadWriteOnce and should not define a storageClassName. The PVC should bound to the PV correctly.

3. Finally create a new Deployment `safari` in Namespace `project-tiger` which mounts that volume at `/tmp/safari-data`. The Pods of that Deployment should be of image `httpd:2.4.41-alpine`.


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


The metrics-server has been installed in the cluster. Your college would like to know the kubectl commands to:
1. show Nodes resource usage
2. show Pods and their containers resource usage
Please write the commands into `/opt/course/7/node.sh` and `/`opt/course/7/pod.sh`.


### Result

```shell
# k top node > /opt/course/7/node.sh
# k top pod --containers=true > /opt/course/7/pod.sh
```



## Question 8 | Get Master Information (2%)

### Describe

Ssh into the master node with `ssh cluster1-master1`. Check how the master components kubelet, kube-apiserver, kube-scheduler, kube-controller-manager and etcd are started/installed on the master node. Also find out the name of the DNS application and how it's started/installed on the master node.

Write your findings into file `/opt/course/8/master-components.txt`. The file should be structured like:

# /opt/course/8/master-components.txt
kubelet: [TYPE]
kube-apiserver: [TYPE]
kube-scheduler: [TYPE]
kube-controller-manager: [TYPE]
etcd: [TYPE]
dns: [TYPE] [NAME]
Choices of [TYPE] are: not-installed, process, static-pod, pod

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

4. Start the kube-scheduler again and confirm its running correctly by 

5. creating a second Pod named `manual-schedule2` of image `httpd:2.4-alpine` and check if it's running on cluster2-worker1.

### Result

```shell
# 1
# ssh cluster2-master1
# cd /etc/kunernetes/manifests/
# mv kube-scheduler.yaml ..
# kubectl -n kube-system get pod | grep schedule

# 2 
# k run manual-schedule --image=httpd:2.4-alpine $do > 9.yaml
# k create -f 9.yaml
# k get pod manual-schedule -w

# 3 
# k get pod manual-schedule -o yaml > 9.yaml
# k replace --force -f 9.yaml
# k get pod manual-schedule -w

# 4
# cd /etc/kunernetes/manifests/
# mv ../kube-scheduler.yaml .
# k get pod -n kube-system -w

# 5
# k run manual-schedule2 --image=httpd:2.4-alpine
# k get pod manual-schedule2 -w
```

## Question 10 | RBAC ServiceAccount Role RoleBinding (6%)

### Describe

1. Create a new ServiceAccount `processor` in Namespace `project-hamster`. 
2. Create a Role and RoleBinding, both named `processor` as well. 
3. hese should allow the new SA to only create Secrets and ConfigMaps in that Namespace.

### Result

```shell
# 1
# k create sa processor -n project-hamster

# 2 Searching RBAC on k8s doc
# cat <<EOF >>10.yaml
> apiVersion: rbac.authorization.k8s.io/v1
> kind: Role
> metadata:
>   namespace: default
>   name: pod-reader
> rules:
> - apiGroups: [""] # "" indicates the core API group
>   resources: ["pods"]
>   verbs: ["get", "watch", "list"]
> EOF

# vi 10.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: project-hamster    # change
  name: processor   # change
rules:
- apiGroups: [""] 
  resources: ["configmaps", "secrets"]    #change
  verbs: ["create"]   #change
---

# cat <<EOF >>10.yaml 
> apiVersion: rbac.authorization.k8s.io/v1
> # This role binding allows "jane" to read pods in the "default" namespace.
> # You need to already have a Role named "pod-reader" in that namespace.
> kind: RoleBinding
> metadata:
>   name: read-pods
>   namespace: default
> subjects:
> # You can specify more than one "subject"
> - kind: User
>   name: jane # "name" is case sensitive
>   apiGroup: rbac.authorization.k8s.io
> roleRef:
>   # "roleRef" specifies the binding to a Role / ClusterRole
>   kind: Role #this must be Role or ClusterRole
>   name: pod-reader # this must match the name of the Role or ClusterRole you wish to bind to
>   apiGroup: rbac.authorization.k8s.io
> EOF

# vi 10.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: project-hamster    # change
  name: processor   # change
rules:
- apiGroups: [""] 
  resources: ["configmaps", "secrets"]    #change
  verbs: ["create"]   #change
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: processor   # change
  namespace: project-hamster    # change
subjects:
- kind: ServiceAccount    # change
  name: processor   # change
  namespace: project-hamster    # change
roleRef:
  kind: Role
  name: processor   # change
  apiGroup: rbac.authorization.k8s.io

# k create -f 10.yaml

# 3 
# k -n project-hamster auth can-i create secret --as system:serviceaccount:project-hamster:processor
# k -n project-hamster auth can-i create configmap --as system:serviceaccount:project-hamster:processor
# k -n project-hamster auth can-i create pod --as system:serviceaccount:project-hamster:processor
```

## Question 11 | DaemonSet on all Nodes (4%)

### Describe

Use Namespace `project-tiger` for the following. 
1. Create a DaemonSet named `ds-important` with image `httpd:2.4-alpine` and labels `id=ds-important` and `uuid=18426a0b-5f59-4e10-923f-c0e078e82462`. 
2. The Pods it creates should request 10 millicore cpu and 10 mebibyte memory. 
3. The Pods of that DaemonSet should run on all nodes, master and worker.

### Result

```shell
# 1
# k -n project-tiger create deployment ds-important --image=httpd:2.4-alpine $do > 11.yaml
# vim 11.yaml
apiVersion: apps/v1
kind: DaemonSet   # change
metadata:
  creationTimestamp: null
  labels:
    id: ds-important   # add 
    uuid: 18426a0b-5f59-4e10-923f-c0e078e82462    # add
  name: ds-important
  namespace: project-tiger
spec:
  selector:
    matchLabels:
      app: ds-important
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ds-important
    spec:
      containers:
      - image: httpd:2.4-alpine
        name: httpd
        resources: {}
# status: {}    remove

# 2
# vim 11.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  creationTimestamp: null
  labels:
    id: ds-important
    uuid: 18426a0b-5f59-4e10-923f-c0e078e82462
  name: ds-important
  namespace: project-tiger
spec:
  selector:
    matchLabels:
      app: ds-important
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ds-important
    spec:
      containers:
      - image: httpd:2.4-alpine
        name: httpd
        resources:    # add
          requests:   # add
            cpu: "10m"    # add
            memory: 10Mi    # add
# k create -f 11.yaml

# 3
# k describe nodes cka-master |grep Taint -A1
# vi 11.yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  creationTimestamp: null
  labels:
    id: ds-important
    uuid: 18426a0b-5f59-4e10-923f-c0e078e82462
  name: ds-important
  namespace: project-tiger
spec:
  selector:
    matchLabels:
      app: ds-important
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ds-important
    spec:
      containers:
      - image: httpd:2.4-alpine
        name: httpd
        resources:
          requests:
            cpu: "10m"
            memory: 10Mi
      tolerations:    # add
      - effect: NoSchedule    # add                               
        key: node-role.kubernetes.io/master   # add
      - effect: NoSchedule    # add
        key: node-role.kubernetes.io/control-plane  # add
# k replace --force -f 11.yaml
```

## Question 12 | Deployment on all Nodes (6%)

### Describe

Use Namespace `project-tiger` for the following. 
1. Create a Deployment named `deploy-important` with label `id=very-important` (the Pods should also have this label) and 3 replicas. It should contain two containers, the first named `container1` with image `nginx:1.17.6-alpine` and the second one named `container2` with image `kubernetes/pause`.

There should be only ever one Pod of that Deployment running on one worker node. We have two worker nodes: `cluster1-worker1` and `cluster1-worker2`. 
2. Because the Deployment has three replicas the result should be that on both nodes one Pod is running. The third Pod won't be scheduled, unless a new worker node will be added.
3. In a way we kind of simulate the behaviour of a DaemonSet here, but using a Deployment and a fixed number of replicas.

 

### Result

```shell
# 1 
# k -n project-tiger create deployment deploy-important --image=nginx:1.17.6-alpine --replicas=3 $do > 12.yaml
# vim 12.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    id: very-important    # add
  name: deploy-important
  namespace: project-tiger
spec:
  replicas: 3
  selector:
    matchLabels:
      id: very-important    # change
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        id: very-important    # change
    spec:
      containers:
      - image: nginx:1.17.6-alpine
        name: container1    # change
        resources: {}
      containers:
      - image: kubernetes/pause   # change
        name: container2    # change
        resources: {}
status: {}

# 2 
# vim 12.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  creationTimestamp: null
  labels:
    id: very-important    # add
  name: deploy-important
  namespace: project-tiger
spec:
  replicas: 3
  selector:
    matchLabels:
      id: very-important    # change
  strategy: {}
  template:
    metadata:
      creationTimestamp: null
      labels:
        id: very-important    # change
    spec:
      containers:
      - image: nginx:1.17.6-alpine
        name: container1    # change
        resources: {}
      - image: kubernetes/pause   # change
        name: container2    # change
        resources: {}
      affinity:   # add
        podAntiAffinity:    # add
          requiredDuringSchedulingIgnoredDuringExecution:   # add
          - labelSelector:    # add
              matchExpressions:    # add
              - key: id    # add
                operator: In    # add
                values:   # add
                - very-important    # add
            topologyKey: kubernetes.io/hostname   # add
      tolerations:    # add
      - effect: NoSchedule    # add
        key: node-role.kubernetes.io/master   # add
      - effect: NoSchedule    # add
        key: node-role.kubernetes.io/control-plane    # add
status: {}

# 3 
# k -n project-tiger get pod -o wide -l id=very-important

```

## Question 13 | Multi Containers and Pod shared Volume (4%)

### Describe

### Result


## Question 1 | Contexts (1%)

### Describe

### Result