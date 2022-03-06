### 练习9.1：测试对CPU的校验和准入行为
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: big-cpu
spec:
  nodeName: cn-node1
  containers: 
    - name: nginx-demo
      image: nginx
      resources:
        requests:
          cpu: 100
        limits:
          cpu: 100
```