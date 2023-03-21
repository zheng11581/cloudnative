## 1. Apparmor

### Task

On the cluster's worker node, enforce the prepared AppArmor profile located at `/etc/apparmor.d/nginx_apparmor`. 
Edit the prepared manifest file located at `/cks/1/pod1.yaml` to apply the AppArmor profile.
Finally, apply the manifest file and create the pod specified in it.

### Spec

```shell
# cat <<EOF >>/etc/apparmor.d/nginx_apparmor
#include <tunables/global>
profile nginx-profile flags=(attach_disconnected) {
 #include <abstractions/base>
 file,
 # Deny all file writes.
 deny /** w,
}

# cat <<EOF >>/cks/1/pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello-apparmor
  annotations:
spec:
  containers:
  - name: hello
    image: busybox
    command: [ "sh", "-c", "echo 'Hello AppArmor!' && sleep 1h" ]

```

### Answer

```shell
# enforce the prepared AppArmor profile located at `/etc/apparmor.d/nginx_apparmor
# ssh <node>
# apparmor_parser -q /etc/apparmor.d/nginx_apparmor
# apparmor_status --pretty-json |grep nginx
"nginx-profile": "enforce",

# k create -f /cks/1/pod.yaml
# k edit pod hello-apparmor 
  annotations:
    # 告知 Kubernetes 去应用 AppArmor 配置 "nginx-profile"
    # container.apparmor.security.beta.kubernetes.io/<container_name>: <profile_ref>
    container.apparmor.security.beta.kubernetes.io/hello: localhost/nginx-profile

# k get events | grep hello-apparmor
60m         Normal   Killing     pod/hello-apparmor   Stopping container hello
58m         Normal   Scheduled   pod/hello-apparmor   Successfully assigned default/hello-apparmor to cks-node
58m         Normal   Pulling     pod/hello-apparmor   Pulling image "busybox"
58m         Normal   Pulled      pod/hello-apparmor   Successfully pulled image "busybox" in 2.424343562s
58m         Normal   Created     pod/hello-apparmor   Created container hello
58m         Normal   Started     pod/hello-apparmor   Started container hello


# k exec hello-apparmor -- cat /proc/1/attr/current
nginx-profile (enforce)

# k exec hello-apparmor -- touch /tmp/test
touch: /tmp/test: Permission denied
command terminated with exit code 1
```

### Note




## 2. kube-bench 对 k8s 进行基准测试

### Task
ACIS Benchmark tool was run against the kubeadm-created cluster and found multiple issues that must be addressed immediately.

Fix all issues via configuration and restart theaffected components to ensure the new settings take effect.
1. Fix all of the following violations that were found against the API server:
- Ensure that the 1.2.7 --authorization-mode FAIL argument is not set to AlwaysAllow
- Ensure that the 1.2.8 --authorization-mode FAIL argument includes Node
- Ensure that the 1.2.9 --authorization-mode FAIL argument includes RBAC
- Ensure that the 1.2.18 --insecure-bind-address FAIL argument is not set
- Ensure that the 1.2.19 --insecure-port FAIL argument is set to 0

2. Fix all of the following violations that were found against the kubelet:
- Ensure that the 4.2.1 anonymous-auth FAIL argument is set to false
- Ensure that the 4.2.2 --authorization-mode FAIL argument is not set to AlwaysAllow. Use webhook authn/authz where possible.

3. Fix all of the following violations that were found against etcd:
- Ensure that the 4.2.1 --client-cert-auth FAIL argument is set to true

### Answer

```shell
# docker run --rm -v `pwd`:/host aquasec/kube-bench:latest install 
# ./kube-bench master

# 1 
# cat /etc/kubenetes/manifest/kube-apiserver.yaml |egrep "authorization-mode|insecure-bind-address|insecure-port"
# vim /etc/kubenetes/manifest/kube-apiserver.yaml
- --authorization-mode=Node,RBAC
- --insecure-port
# delete --insecure-bind-address

# 2 Master and Node
# vim /var/lib/kubelet/config.yaml
authentication:
  anonymous:
    enabled: false
authorization:
  mode: Webhook
  webhook:
    cacheAuthorizedTTL: 0s
    cacheUnauthorizedTTL: 0s

# 3
# vim /etc/kubernetes/manifests/etcd.yaml
- --client-cert-auth=true

```

### Note


## 3. Trivy 进行镜像扫描

### Task
1. Use the Trivy open-source container scanner to detect images with severe vulnerabilities used by Pods in the namespace `yavin`.
2. Look for images with High or Critical severity vulnerabilities,and delete the Pods that use those images.

Trivy is pre-installed on the cluster's master node only; it is not available on the base system or the worker nodes. You'll have to connect to the cluster's master node to use Trivy.

### Answer

```shell
# 1
$ ssh root@kssc00401-master
$ kubectl get pod -n mall-swarm -o jsonpath='{.items[*].spec.containers[0].image}'


# 2 
$ trivy image --skip-update -–severity HIGH,CRITICAL '刚才搜到的镜像' 
$ 
```

### Note



## 4. sysdig & falco

### Task
You may use your browser to open one additional tab to access sysdig's documentation or Falco's documentation.

Use runtime detection tools to detect anomalous processes spawning and executing frequently in the single container belorging to Pod `redis`.
Two tools are available to use:
- sysdig
- falco

The tools are pre-installed on the cluster's worker node only; they are not available on the base system or the master node.

1. Using the tool of your choice (including any non pre-installed tool), analyse the container's behaviour for at least 30 seconds, using filters that detect newly spawning and executing processes.
2. Store an incident file at /opt/2/report , containing the detected incidents, one per line, in the following format:
[timestamp],[uid],[processName]

Keep the tool's original timestamp-format as-is.
Make sure to store the incident fileon the cluster's worker node.

### Answer
```shell
# 1
$ ssh root@vms62.rhce.cc
$ docker ps | grep redis
$ sysdig -l | grep time
$ sysdig -l | grep uid
$ sysdig -l | grep proc

# 2
$ sysdig -M 30 -p "*%evt.time,%user.uid,%proc.name" container.id=846a386c1e34 > /opt/2/report
```
### Note

## 5. serviceAccount

### Task

1. Create a new ServiceAccount named `backend-sa` in the existing namespace `qa` ,which must not have access to any secrets. 
2. Inspect the Pod named `backend` running in the namespace `qa`. Edit the Pod to use the newly created serviceAccount `backend-sa`. You can 
find the Pod's manifest file at `/cks/9/pod9.yaml`. Ensure that the modified specification is applied and the Pod is running. 
3. Finally, clean-up and delete the now unused serviceAccount that the Pod used initially.

### Answer
```shell
# 1
# kubectl -n qa create sa backend-sa --dry-run=client -oyaml > 5-sa.yaml
# vi 5-sa.yaml
# automountServiceAccountToken: false

# 2
# kubectl -n qa run backend --image=nginx --dry-run=client -oyaml > /cks/5/5-pod.yaml
# vi /cks/5/5-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: backend
  name: backend
  namespace: qa
spec:
  serviceAccountName: backend-sa # New add
  containers:
  - image: nginx
    name: backend
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always
status: {}

# kubectl apply -f /cks/5/5-pod.yaml

# 3
# kubecel -n qa delete sa <not-backend-sa>
```
### Note

6. Pod 安全策略-PodSecurityPolicy

### Task

1. Create a new PodSecurityPolicy named prevent-psp-policy, which prevents the creation of privileged Pods.
2. Create a new ClusterRole named restrict-access-role, which uses the newly created PodSecurityPolicy prevent-psp-policy
3. Create a new serviceAccount named psp-denial-sa in the existing namespace development.
4. Finally, create a new clusterRoleBinding named dany-access-bind, which binds the newly created ClusterRole restrict-access-role to the newly created serviceAccount psp-denial-sa.

### Answer
```shell
# 0
# vim /etc/kubernetes/manifests/kube-apiserver.yaml
- --enable-admission-plugins=NodeRestriction,PodSecurityPolicy
# systemctl restart kubelet

# 1
# vim /cks/6/6-psp.yaml
apiVersion: policy/v1beta1
kind: PodSecurityPolicy
metadata: 
  name: prevent-psp-policy
spec:
  privileged: false #禁止容器的特权模式
  seLinux:
    rule: RunAsAny #设置 selinux 参数，不会校验
  supplementalGroups:
    rule: RunAsAny
  runAsUser:
    rule: RunAsAny #设置运行容器的用户 ID
  fsGroup:
    rule: RunAsAny #设置组 ID
  volumes:
  - '*' #可以用任何 volume 卷
# kubectl -f apply -f /cks/6/6-psp.yamlm bnmjjj
```

7. NetworkPolicy 拒绝所有入口流量

### Task

1. Create a new default-deny NetworkPolicy named `denynetwork` in the namespace `development` for all traffic of type Ingress.
2. The new NetworkPolicy must deny all lngress traffic in the namespace `development`.
3. Apply the newly created default-deny NetworkPolicy to all Pods running in namespace development.

### Answer
```shell
cat /cks/7/7-networkpolicy.yaml <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: denynetwork
  namespace: development
spec:
  podSelector: {}
  policyTypes:
  - Ingres
EOF
```

8. NetworkPolicy

### Task

1. create a NetworkPolicy named `pod-access` to restrict access to Pod `products-service` running in namespace `development`.
2. only allow the following Pods to connect to Pod `products-service`:
    - Pods in the namespace `testing`
    - Pods with label `environment: staging`, in any namespace
3. Make sure to apply the NetworkPolicy.
You can find a skelet on manifest file at /cks/8/8-networkpolicy.yaml

### Answer
```shell
cat > /cks/8/8-networkpolicy.yaml <<EOF
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: pod-access
  namespace: development
spec:
  podSelector:
    matchLabels:
      run: products-service
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          name: testing
    - podSelector:
        matchLabels:
          environment: staging
EOF
```


9.  RBAC/clusterrole

### Task
A Role bound to a Pod's serviceAccount grants overly permissive permissions. Complete the following tasks to reduce the set of permissions.

1. Given an existing Pod named `web-pod` running in the namespace `monitoring`. 
2. Edit the existing Role bound to the Pod's serviceAccount `sa-1` to only allow performing `list` operations, only on resources of type `Endpoints`.
3. Create a new Role named `role-2` in the namespace `monitoring`, which only allows performing，`update` operations, only on resources of type `persistentvolumeclaims`.
4. Create a new RoleBinding named `role-2-binding` binding the newly created Role to the Pod's serviceAccount.
Don't delete the existing RoleBinding.

```shell
# 1 & 2 find the RoleBinding with ServiceAccount sa-1 and edit the Role with list operations on Endpoints
kubectl get rolebindings.rbac.authorization.k8s.io -n monitoring |grep sa-1

# 3 
cat > /cks/9/9-role.yaml <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  namespace: monitoring
  name: role-2
rules:
- apiGroups: [""] 
  resources: ["persistentvolumeclaims"]
  verbs: ["update"]
EOF

# 4 
cat > /cks/9/9-rolebinding.yaml <<EOF
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: role-2-binding
  namespace: monitoring
subjects:
- kind: ServiceAccount
  name: monitoring:sa-1
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role 
  name: role-2
  apiGroup: rbac.authorization.k8s.io
EOF
```

10. kube-apiserver 审计日志记录和采集

### Task
1. Enable audit logs in the cluster.
To do so, enable the log backend, and ensure that:
- logs are stored at /var/log/kubernetes/audit-logs.txt
- log files are retained for 5 days
- at maximum, a number of 10 auditlog files are retained

2. A basic policy is provided at /etc/kubernetes/logpolicy/sample-policy.yaml. it only specifies what not to do.
The base policy is located on the cluster's master node.
Edit and extend the basic policy to log:
- `namespaces` changes at RequestResponse level
- the request body of `pods` changes in the namespace `front-apps`
- `configMap` and `secret` changes in all namespaces at the Metadata level
- Also, add a catch all rule to log all other requests at the Metadata level.
Don't forget to apply the modifiedpolicy

```shell
# 1 Edit kube-apiserver.yaml
--audit-policy-file=/etc/kubernetes/logpolicy/sample-policy.yaml
--audit-log-path=/var/log/kubernetes/audit-logs.txt
--audit-log-maxage=5
--audit-log-maxbackup=10

# 2 Edit /etc/kubernetes/logpolicy/sample-policy.yaml

apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: RequestResponse
  resources:
  - group: ""
    resources: ["namespaces"]

- level: Request
  resources:
  - group: "" 
    resources: ["pods"]
    namespaces: ["front-apps"]

- level: Metadata
  resources:
  - group: ""
    resources: ["secrets", "configmaps"]

- level: Metadata
  omitStages:
  - "RequestReceived"

```


11. 创建 Secret

### Task
1. Retrieve the content of the existing secret named `db1-secret-test` in the `istio-system` namespace. store the `username` field in a file named `/cks/11/old-username.txt` , and the `password` field in a file named `/cks/11/old-pass.txt`. You must create both files; they don't exist yet.
2. Do not modify the created files in the following steps, create new temporary files if needed.
- Create a new secret named `test` in the `istio-system` namespace, with the following
content:
 username : thanos
 password : hahahaha
3. Finally, create a new Pod that has access to the secret `test` via a volume:
pod name | dev-pod
namespace | istio-system
container name | dev-container
image | nginx:1.9
volume name | dev-volume
mount path | /etc/test-secret

```shell
# 1 
# kubectl create secret generic db1-secret-test --from-file=username=/cks/11/old-username.txt --from-file=password=/cks/11/old-pass.txt -n istio-system
kubectl get secrets -n istio-system db1-secret-test -ojsonpath={.data.password} |base64 -d > /cks/11/old-pass.txt
kubectl get secrets -n istio-system db1-secret-test -ojsonpath={.data.username} |base64 -d > /cks/11/old-username.txt

# 2
kubectl create secret generic test --from-literal=username=thanos --from-literal=password=hahahaha -n istio-system --dry-run=client -oyaml > /cks/11/11-secret.yaml
kubect; apply -f /cks/11/11-secret.yaml

# 3 
kubectl run dev-pod -n istio-system --image=nginx:1.9 --dry-run=client -oyaml > /cks/11/11-pod.yaml
# Edit /cks/11/11-pod.yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: dev-pod
  name: dev-pod
  namespace: istio-system
spec:
  volumes:
  - name: dev-volume
    secret:
      secretName: test
  containers:
  - image: nginx:1.9
    name: dev-container
    volumeMounts:
    - name: dev-volume
      readOnly: true
      mountPath: "/etc/test-secret"
  dnsPolicy: ClusterFirst
  restartPolicy: Always
# Apply /cks/11/11-pod.yaml
kubectl apply -f /cks/11/11-pod.yaml
```


12. dockerfile 和 deployment 优化部分

### Task
1. Analyze and edit the given Dockerfile (based on the ubuntu:16.04 image) `/cks/12/Dockerfile` fixing `two instructions` present in the file being prominent security/best-practice issues.
2. Analyze and edit the given manifest file `/cks/12/deployment.yaml` fixing `two fields` present in the file being prominent security/best-practice issues.

```shell
# 1 
# vim /cks/12/Dockerfile
# FROM ubuntu:latest
FROM ubuntu:16.04 #把 latest 变成 16.04
# USER root 
USER nobody #把 USER ROOT 变成 USER nobody

# 2  
# vim /cks/12/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
labels:
app: dev
name: dev
spec:
replicas: 1
selector:
matchLabels:
app: dev
template:
 metadata:
 labels:
 app: dev
 spec:
 containers:
 - image: mysql
 name: mysql
 securityContext: {'capabilities':{'add':['NET_ADMIN'],'drop':['all']},'privileged': False,'readOnlyRootFilesystem': True, 'runAsUser': 65535}
# 'privileged'变成 false，readOnlyRootFilesystem’变成 True
```

13. 镜像安全 ImagePolicyWebhook

### Task
You have to complete the entire task on the cluster's master node, where all services and files have been prepared and placed.
Given an incomplete configuration in directory `/etc/kubernetes/config` and a functional container image scanner with HTTPS endpoint `https://localhost:8081/image_policy`

- Enable the necessary plugins to create an image policy
- validate the control configuration and change it to an `implicit deny`
- Edit the configuration to point the provided HTTPS endpoint correctly

Finally , test if the configuration is working by trying to deploy the vulnerable resource /cks/1/web1.yaml
You can find the container image scanner's log file at /var/loglimagepolicyiacme.log

```shell
--enable-admission-plugins
--admission-control-config-file
/etc/kubernetes/config/imagepolicyconfig.yaml
kubeConfigFile
```


14. 删除非无状态或非不可变的 pod



15. gVsior/runtimeclass

### Task
1. Create a RuntimeClass named `untrusted` using the prepared runtime handler named `runsc`.
2. Update all Pods in the namespace `client` to run on gvisor, unless they are already running on anon-default runtime handler.
You can find a skeleton manifest file at `/cks/15/rc.yaml`

```shell
# 1
cat > /cks/15/rc.yaml <<EOF
apiVersion: node.k8s.io/v1
kind: RuntimeClass
metadata:
  name: untrusted
handler: runsc
EOF

# 2
kubectl get pod -n client -oyaml > /cks/15/pods.yaml
cp /cks/15/pods.yaml /cks/15/runsc-pods.yaml
# Edit /cks/15/runsc-pods.yaml 
runtimeClassName: untrusted
kubectl delete -f /cks/15/runsc-pods.yaml
kubectl create -f /cks/15/runsc-pods.yaml
```
