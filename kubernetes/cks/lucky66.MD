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


```


