### Error: x509: certificate signed by unknown authority
正常情况下，Linux docker 配置完ca.crt后，就可以登录Harbor了。
但是如果还是报上面的错，按照一下方式修改 (docker desktop也可以） 

#### Mac
```shell
echo -n | openssl s_client -showcerts -connect 192.168.110.72:443 2>/dev/null | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > /Users/zhenghc/goharbor.com.cert
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain /Users/zhenghc/goharbor.com.cert
# Restart docker desktop
```

#### Win10
```txt
控制面板->管理计算机证书
授信的根证书颁发机构->右击->所有任务->导入
选择ca.crt->一直Next完成
Restart Docker for Windows
```

#### Ubuntu
```shell
cp yourdomain.com.crt /usr/local/share/ca-certificates/yourdomain.com.crt 
update-ca-certificates
systemclt restart docker
```

#### CentOS
```shell
cp yourdomain.com.crt /etc/pki/ca-trust/source/anchors/yourdomain.com.crt
update-ca-trust
systemclt restart docker
```
