### 对于non-namespace和non-resource，需要使用ClusterRole和ClusterRoleBinding进行授权
### Create ClusterRole <pv-reader> has <get,list,watch> </api/v1/persistentvolumes> permissions
```shell
k create clusterrole pv-reader --verb=get,list,watch --resource=persistentvolumes
root@cn-master1:~# k get clusterrole pv-reader -oyaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  creationTimestamp: "2022-02-09T07:54:33Z"
  name: pv-reader
  resourceVersion: "5840129"
  uid: d7f6013c-143a-4aef-8d41-2f6d6a007378
rules:
- apiGroups:
  - ""
  resources:
  - persistentvolumes
  verbs:
  - get
  - list
  - watch
```

### Use RoleBinding bind ClusterRole
```shell
k create rolebinding pv-test --clusterrole=pv-reader --serviceaccount=qa:default --namespace=qa
```

### Try to get persistent volumes in pod qa/test
```shell
root@cn-master1:~# k exec -it test -n qa -- sh
/ # curl localhost:8001/api/v1/persistentvolumes
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "persistentvolumes is forbidden: User \"system:serviceaccount:qa:default\" cannot list resource \"persistentvolumes\" in API group \"\" at the cluster scope",
  "reason": "Forbidden",
  "details": {
    "kind": "persistentvolumes"
  },
  "code": 403
}
```
Can't get persistent volumes in pod qa/test

### 必须使用ClusterRoleBinding来对集群级别资源进行授权访问
```shell
k delete rolebinding pv-test -n qa
k create clusterrolebinding pv-test --clusterrole=pv-reader --serviceaccount=qa:default
k get clusterrolebinding pv-test -oyaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  creationTimestamp: "2022-02-09T08:18:49Z"
  name: pv-test
  resourceVersion: "5845009"
  uid: f92a5be5-9766-4e9f-9b69-054b06b0717a
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: pv-reader
subjects:
- kind: ServiceAccount
  name: default
  namespace: qa
k exec -it test -n qa -- sh
/ # curl localhost:8001/api/v1/persistentvolumes
{
  "kind": "PersistentVolumeList",
  "apiVersion": "v1",
  "metadata": {
    "resourceVersion": "5845177"
  },
  "items": []
}
```

### 对于访问非资源型的URL授权，通常通过system:discovery ClusterRole和同名的ClusterRoleBinding自动完成
```shell
k get clusterrole system:discovery -oyaml

kind: ClusterRole
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  creationTimestamp: "2022-01-20T01:40:11Z"
  labels:
    kubernetes.io/bootstrapping: rbac-defaults
  name: system:discovery
  resourceVersion: "74"
  uid: 4562dcec-8be1-424b-ad83-b67b8e4b7f83
rules:
- nonResourceURLs:
  - /api
  - /api/*
  - /apis
  - /apis/*
  - /healthz
  - /livez
  - /openapi
  - /openapi/*
  - /readyz
  - /version
  - /version/
  verbs:
  - get

k get clusterrolebinding system:discovery -oyaml

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    rbac.authorization.kubernetes.io/autoupdate: "true"
  creationTimestamp: "2022-01-20T01:40:12Z"
  labels:
    kubernetes.io/bootstrapping: rbac-defaults
  name: system:discovery
  resourceVersion: "138"
  uid: 81972035-dce0-45d7-8e92-ec52538ed2d0
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: system:discovery
subjects:
- apiGroup: rbac.authorization.k8s.io
  kind: Group
  name: system:authenticated
```
表明所有经过认证的用户（组）都可以访问到ClusterRole system:discovery定义的url