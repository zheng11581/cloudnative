#### 1. Create namespace and RBAC 
```shell
kubectl create ns devops
kubectl create sa jenkins -n devops
kubectl create clusterrolebinding jenkins-cluster-admin --clusterrole=cluster-admin --serviceaccount=devops:jenkins
```

#### 2. Create StorageClass if not present

```shell
kubectl apply -f jenkins-storageclass.yaml
```

#### 3. Create PVC for pod
```shell
kubectl apply -f jenkins-pvc.yaml
```

#### 4. Create jenkins deployment
```shell
kubectl apply -f jenkins-deployment.yaml
```

#### 5. Wait for pod is Running 
```shell
kubectl get po -n devops --watch
```

#### 6. Unlock jenkins use initialAdminPassword
```shell
kubectl exec -it -n devops jenkinszh-cf86b456-5gdch -- cat /var/jenkins_home/secrets/initialAdminPassword
```
