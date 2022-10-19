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
1. Use the Trivy open-source container scanner to detect images with severe vulnerabilities used by Pods in the namespace yavin.
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
# kubectl -n qa run backend --image=nginx --dry-run=client -oyaml > /cks/9/pod9.yaml
# vi /cks/9/pod9.yaml
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

# kubectl apply -f /cks/9/pod9.yaml

# 3
# kubecel -n qa delete sa <not-backend-sa>
```
### Note

6. Pod 安全策略-PodSecurityPolicy

7. NetworkPolicy 拒绝所有入口流量

8. NetworkPolicy

9.  RBAC/clusterrole

10. kube-apiserver 审计日志记录和采集

11. 创建 Secret

12. dockerfile 和 deployment 优化部分

13. 镜像安全 ImagePolicyWebhook

14. 删除非无状态或非不可变的 pod

15. gVsior/runtimeclass
