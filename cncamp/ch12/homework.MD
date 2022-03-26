```text
把我们的 httpserver 服务以 Istio Ingress Gateway 的形式发布出来。以下是你需要考虑的几点：

1. 如何实现安全保证
2. 七层路由规则
3. 考虑 open tracing 的接入
```

### 1. 如何实现安全保证: TLS证书 && 2. 七层路由规则: HTTP Rewrite

#### [Deploy httpserver](../httpserver/deploy/httpserver-ingressgw.yaml)

```sh
kubectl create ns tlsserver
kubectl label ns tlsserver istio-injection=enabled
kubectl create -f httpserver-cm.yaml -n tlsserver
kubectl create -f httpserver-deploy.yaml -n tlsserver
```

```sh
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt
kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt
kubectl apply -f httpserver-ingresgw.yaml -n tlsserver
```

#### Check ingress ip

```sh
kubectl get svc -nistio-system

NAME                   TYPE           CLUSTER-IP    EXTERNAL-IP    PORT(S)                                                                      AGE
istio-egressgateway    ClusterIP      10.4.7.11     <none>         80/TCP,443/TCP                                                               17h
istio-ingressgateway   LoadBalancer   10.4.15.166   34.133.91.93   15021:30575/TCP,80:31842/TCP,443:32111/TCP,31400:31736/TCP,15443:30749/TCP   17h
istiod                 ClusterIP      10.4.9.29     <none>         15010/TCP,15012/TCP,443/TCP,15014/TCP                                        17h
```

#### Access the httpserver via ingress

```sh
export INGRESS_IP=34.133.91.93
curl --resolve httpsserver.cncamp.io:443:$INGRESS_IP https://httpsserver.cncamp.io/healthz -v -k

root@curl:/# curl --resolve httpsserver.cncamp.io:443:$INGRESS_IP https://httpsserver.cncamp.io/out/healthz -v -k
* Added httpsserver.cncamp.io:443:34.133.91.93 to DNS cache
* Hostname was found in DNS cache
*   Trying 34.133.91.93...
* Connected to httpsserver.cncamp.io (34.133.91.93) port 443 (#0)
* successfully set certificate verify locations:
*   CAfile: none
  CApath: /etc/ssl/certs
* SSLv3, TLS handshake, Client hello (1):
* SSLv3, TLS handshake, Server hello (2):
* SSLv3, TLS handshake, CERT (11):
* SSLv3, TLS handshake, Server key exchange (12):
* SSLv3, TLS handshake, Server finished (14):
* SSLv3, TLS handshake, Client key exchange (16):
* SSLv3, TLS change cipher, Client hello (1):
* SSLv3, TLS handshake, Finished (20):
* SSLv3, TLS change cipher, Client hello (1):
* SSLv3, TLS handshake, Finished (20):
* SSL connection using ECDHE-RSA-AES128-GCM-SHA256
* Server certificate:
*        subject: O=cncamp Inc.; CN=*.cncamp.io
*        start date: 2022-03-25 10:04:42 GMT
*        expire date: 2023-03-25 10:04:42 GMT
*        issuer: O=cncamp Inc.; CN=*.cncamp.io
*        SSL certificate verify result: self signed certificate (18), continuing anyway.
> GET /out/healthz HTTP/1.1
> User-Agent: curl/7.35.0
> Host: httpsserver.cncamp.io
> Accept: */*
>
< HTTP/1.1 200 OK
< date: Fri, 25 Mar 2022 11:08:07 GMT
< content-length: 3
< content-type: text/plain; charset=utf-8
< x-envoy-upstream-service-time: 2
* Server istio-envoy is not blacklisted
< server: istio-envoy
<
ok
* Connection #0 to host httpsserver.cncamp.io left intact

```
#### 3. Open Tracing

[https://gitee.com/zheng11581/httpserver](https://gitee.com/zheng11581/httpserver)