#### 1. Create namespace and RBAC 
```shell
k create ns devops
k create sa jenkins -n devops
k create clusterrolebinding jenkins-cluster-admin --clusterrole=cluster-admin --serviceaccount=devops:jenkins
```

#### 2. Create StorageClass if not present

```shell
k apply -f jenkins-storageclass.yaml
```

#### 3. Create PVC for pod
```shell
k apply -f jenkins-pvc.yaml
```

#### 4. Create jenkins deployment
```shell
k apply -f jenkins-deployment.yaml
```

#### 5. Wait for pod is Running 
```shell
k get po -n devops --watch
```

#### 6. Unlock jenkins use initialAdminPassword
```shell
k exec -it -n devops jenkinszh-cf86b456-5gdch -- cat /var/jenkins_home/secrets/initialAdminPassword
```
