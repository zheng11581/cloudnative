### Create two namespaces and two pods
```shell
k create ns prod
k create ns qa
k run test --image=luksa/kubectl-proxy -n prod
k run test --image=luksa/kubectl-proxy -n qa
```
Note: kubectl proxy expose 8001 port for pod connecting to api-server

### Try to list services from pod qa/test
```shell
k exec -it test -n qa -- sh
/ # curl localhost:8001/api/v1/namespaces/qa/services
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "services is forbidden: User \"system:serviceaccount:qa:default\" cannot list resource \"services\" in API group \"\" in the namespace \"qa\"",
  "reason": "Forbidden",
  "details": {
    "kind": "services"
  },
  "code": 403
}
```
The pod qa/test has no permission to list service in the qa namespace

### Create role <service-reader> which has <get,list,watch> </api/v1/namespaces/qa/services> permissions in qa:default
### Bind <service-reader> to qa:default service account
```shell
k create role service-reader --verb=get,list,watch --resource=services --namespace=qa
k create rolebinding test --role=service-reader --serviceaccount=qa:default --namespace=qa
```

### Try to list services from pod qa/test again
```shell
k exec -it test -n qa -- sh
/ # curl localhost:8001/api/v1/namespaces/qa/services
{
  "kind": "ServiceList",
  "apiVersion": "v1",
  "metadata": {
    "resourceVersion": "5827902"
  },
  "items": []
}
```
The pod qa/test can list services in the qa namespace

### Try to list other pods from pod qa/test
```shell
/ # curl localhost:8001/api/v1/namespaces/qa/pods
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "pods is forbidden: User \"system:serviceaccount:qa:default\" cannot list resource \"pods\" in API group \"\" in the namespace \"qa\"",
  "reason": "Forbidden",
  "details": {
    "kind": "pods"
  },
  "code": 403
}
```
The pod qa/test can't list pods in the qa namespace

### Try to list services from prod namespace in pod qa/test
```shell
/ # curl localhost:8001/api/v1/namespaces/prod/services
{
  "kind": "Status",
  "apiVersion": "v1",
  "metadata": {},
  "status": "Failure",
  "message": "services is forbidden: User \"system:serviceaccount:qa:default\" cannot list resource \"services\" in API group \"\" in the namespace \"prod\"",
  "reason": "Forbidden",
  "details": {
    "kind": "services"
  },
  "code": 403
}
```
The pod qa/test can't list services in the prod namespace

### Bind <service-reader> to prod:default service account
```shell
k create role service-reader --verb=get,list,watch --resource=services --namespace=prod
k create rolebinding test --role=service-reader --serviceaccount=qa:default --serviceaccount=prod:default --namespace=prod
```

### Try to list services from prod namespace in pod qa/test
```shell
/ # curl localhost:8001/api/v1/namespaces/prod/services
{
  "kind": "ServiceList",
  "apiVersion": "v1",
  "metadata": {
    "resourceVersion": "5835338"
  },
  "items": []
}
```
You can list services in prod namespace from pod qa/test