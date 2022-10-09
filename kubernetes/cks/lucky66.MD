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


## 2. Kube-bench

### Task
ACIS Benchmark tool was run against the kubeadm-created cluster and found multiple issues that must be addressed immediately.

Fix all issues via configuration and restart theaffected components to ensure the new settings take effect.
Fix all of the following violations that were found against the API server:
- Ensure that the 1.2.7 --authorization-mode FAIL argument is not set to AlwaysAllow
- Ensure that the 1.2.8 --authorization-mode FAIL argument includes Node
- Ensure that the 1.2.9 --authorization-mode FAIL argument includes RBAC
- Ensure that the 1.2.18 --insecure-bind-address FAIL argument is not set
- Ensure that the 1.2.19 --insecure-port FAIL argument is set to 0


### Answer

```shell
# ctr i pull docker.io/aquasec/kube-bench:latest
# ctr c create --mount type=bind,src=`pwd`,dst=/host,options=rbind:ro docker.io/aquasec/kube-bench:latest install ./kube-bench cks-master
# docker run --rm -v `pwd`:/host aquasec/kube-bench:latest install ./kube-bench cks-master

# cat /etc/kubenetes/manifest/kube-apiserver.yaml |egrep "authorization-mode|insecure-bind-address|insecure-port"
# vim /etc/kubenetes/manifest/kube-apiserver.yaml
- --authorization-mode=Node,RBAC
- --insecure-port
# delete --insecure-bind-address


```