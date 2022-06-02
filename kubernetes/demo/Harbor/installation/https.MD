### 1. 配置HTTPS访问Harbor

#### 生成CA证书
```shell
# Generate a CA certificate private key
openssl genrsa -out ca.key 4096

# Generate the CA certificate
openssl req -x509 -new -nodes -sha512 -days 3650 \
 -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=yourdomain.com" \
 -key ca.key \
 -out ca.crt
```

#### 生成服务器证书
```shell
# Generate a private key
openssl genrsa -out yourdomain.com.key 4096

# Generate a certificate signing request (CSR)
openssl req -sha512 -new \
    -subj "/C=CN/ST=Beijing/L=Beijing/O=example/OU=Personal/CN=yourdomain.com" \
    -key yourdomain.com.key \
    -out yourdomain.com.csr

# Generate an x509 v3 extension file
cat > v3.ext <<-EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1=yourdomain.com
DNS.2=yourdomain
DNS.3=hostname
EOF

# Use the v3.ext file to generate a certificate for your Harbor host
openssl x509 -req -sha512 -days 3650 \
    -extfile v3.ext \
    -CA ca.crt -CAkey ca.key -CAcreateserial \
    -in yourdomain.com.csr \
    -out yourdomain.com.crt
```

#### 重新Harbor配置

```shell
# Copy certificate and key to right place
cp yourdomain.com.crt /data/cert/
cp yourdomain.com.key /data/cert/

# Reconfigure Harbor
vim /glzt/harbor/harbor.yml

hostname: yourdomain.com
https:
  # https port for harbor, default is 443
  port: 443
  # The path of cert and key files for nginx
  certificate: /data/cert/yourdomain.com.crt
  private_key: /data/cert/yourdomain.com.key
  
```

#### 重新配置Docker（Harbor host）

```shell
# Convert yourdomain.com.crt to yourdomain.com.cert, for use by Docker
openssl x509 -inform PEM -in yourdomain.com.crt -out yourdomain.com.cert

# Copy the server certificate, key and CA files into the Docker certificates folder on the Harbor host. You must create the appropriate folders first.
cp yourdomain.com.cert /etc/docker/certs.d/yourdomain.com/
cp yourdomain.com.key /etc/docker/certs.d/yourdomain.com/
cp ca.crt /etc/docker/certs.d/yourdomain.com/

# Restart Docker Engine
systemctl restart docker
```

#### 重启Harbor
```shell
./prepare
docker-compose down -v
docker-compose up -d
```

#### 配置客户端Docker 

Harbor新的版本可以下载ca.crt，旧的不可以只能手动分发CA证书

```shell
mkdir -p /etc/docker/certs.d/yourdomain.com/
cp ca.crt /etc/docker/certs.d/yourdomain.com/
systemctl restart docker
```
