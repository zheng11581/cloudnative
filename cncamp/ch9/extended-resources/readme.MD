### 1. Start kubeclt proxy
```shell
kubectl proxy 
```
### 2. Deliver patch to api-server
```sh
curl -XPATCH --header "Content-Type: application/json-patch+json" \
--data '[{"op": "add", "path": "/status/capacity/cncamp.com~1reclaimed-cpu", "value": "2"}]' \
http://localhost:8001/api/v1/nodes/cn-node2/status
```

### 3. Apply a pod request the extended resource

```shell
$ k apply -f nginx-deploy.yaml 
$ k get po nginx-deployment-7d95bd597d-bbnww -owide
NAME                                READY   STATUS    RESTARTS   AGE   IP       NODE     NOMINATED NODE   READINESS GATES
nginx-deployment-7d95bd597d-bbnww   0/1     Pending   0          42s   <none>   <none>   <none>           <none>

```

### 4. Adjust cncamp.com/reclaim-cpu=3 on cn-node2

```shell
$ curl -XPATCH --header "Content-Type: application/json-patch+json" \
--data '[{"op": "add", "path": "/status/capacity/cncamp.com~1reclaimed-cpu", "value": "3"}]' \
http://localhost:8001/api/v1/nodes/cn-node2/status

$ k get po nginx-deployment-7d95bd597d-bbnww -owide
NAME                                READY   STATUS              RESTARTS   AGE     IP       NODE       NOMINATED NODE   READINESS GATES
nginx-deployment-7d95bd597d-bbnww   0/1     ContainerCreating   0          5m26s   <none>   cn-node2   <none>           <none>
```