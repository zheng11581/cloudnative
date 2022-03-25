### Deploy simple

```sh
kubectl create ns httpserver
kubectl create -f httpserver.yaml -n httpserver
kubectl create -f istio-specs.yaml -n httpserver
```
```text
istio-spec解析：
1. HTTP Request Header中匹配到Host: httpserver.cncamp.io
2. MATCH PORT 80端口
3. ROUTE 到SERVICE httpserver.httpserver.svc.cluster.local的8080端口
```

### Check ingress ip

```sh
k get svc -nistio-system

istio-ingressgateway   LoadBalancer   10.108.31.242
```

### Access the simple via ingress

```sh
export INGRESS_IP=10.108.31.242
curl -H "Host: httpserver.cncamp.io" $INGRESS_IP/hello -v
```

### istioctl pc
```shell
haichengz86@cloudshell:~$ kubectl get po -n istio-system -l istio=ingressgateway
NAME                                    READY   STATUS    RESTARTS   AGE
istio-ingressgateway-76b86f6b45-b5llb   1/1     Running   0          15h

haichengz86@cloudshell:~$ istioctl pc listener istio-ingressgateway-76b86f6b45-b5llb -n istio-system --port 8080
ADDRESS PORT MATCH DESTINATION
0.0.0.0 8080 ALL   Route: http.8080

haichengz86@cloudshell:~$ istioctl pc route istio-ingressgateway-76b86f6b45-b5llb -n istio-system |grep httpserver
http.8080     httpserver.cncamp.io     /*                     httpserver.httpserver

haichengz86@cloudshell:~$ istioctl pc cluster istio-ingressgateway-76b86f6b45-b5llb -n istio-system --port 8080
SERVICE FQDN                                PORT     SUBSET     DIRECTION     TYPE     DESTINATION RULE
httpserver.httpserver.svc.cluster.local     8080     -          outbound      EDS

haichengz86@cloudshell:~$ istioctl pc endpoint istio-ingressgateway-76b86f6b45-b5llb -n istio-system --port 8080 |grep httpserver
10.0.1.15:8080     HEALTHY     OK                outbound|8080||httpserver.httpserver.svc.cluster.local

haichengz86@cloudshell:~$ kubectl get po -n httpserver -owide
NAME                          READY   STATUS    RESTARTS   AGE     IP          NODE                                              NOMINATED NODE   READINESS GATES
httpserver-77fd57475c-q4tgl   1/1     Running   0          3h21m   10.0.1.15   gke-cluster-operator-default-pool-dcd5fb04-prgk   <none>           <none>
```
