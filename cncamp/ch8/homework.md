### 第一部分
#### 1. Make httpserver is a Graceful Shutdown web application
##### Open terminal A
```shell
macbookpro:httpserver zhenghc$ cd $GOPATH/cloudnative/cncamp/httpserver
macbookpro:httpserver zhenghc$ go run main.go
2022/02/21 23:31:50 Starting http server...
```
##### Open terminal B
```shell
macbookpro:~ zhenghc$ ps aux |grep main.go
zhenghc          21300   0.0  0.0  4268284    568 s001  R+   11:30PM   0:00.01 grep main.go
zhenghc          21265   0.0  0.0  4987740   4016 s000  Ss+  11:30PM   0:00.02 go run main.go
macbookpro:~ zhenghc$ kill 21265
```
##### Return to terminal A
```shell
macbookpro:httpserver zhenghc$ cd $GOPATH/cloudnative/cncamp/httpserver
macbookpro:httpserver zhenghc$ go run main.go
2022/02/21 23:31:50 Starting http server...
2022/02/21 23:31:52 Waiting 5 seconds for grace shutdown...
2022/02/21 23:31:57 Starting quit...

Process finished with the exit code 0
```

#### 2. Graceful start
[Add probe](../httpserver/deploy/httpserver-deploy.yaml)

#### 3. Graceful stop
[Use Tiny](../httpserver/Dockerfile)

[Use Chanel](../httpserver/main.go)

#### 4. Resource requests and limits (QOS)
[User resources.requests and limits](../httpserver/deploy/httpserver-deploy.yaml)

#### 5. Configuration Map
[User ConfigMap](../httpserver/deploy/httpserver-cm.yaml)

### 第二部分

#### 1. Service
[Use Service](../httpserver/deploy/httpserver-deploy.yaml)

#### 2. Ingress
[Use Ingress](../httpserver/deploy/httpserver-ingress.yaml)