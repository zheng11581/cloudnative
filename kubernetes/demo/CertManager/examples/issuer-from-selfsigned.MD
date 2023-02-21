### Create SelfSigned `Issuer` or `ClusterIssuer`

```shell
cat > selfsigned-issuer.yaml << EOF
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: selfsigned-issuer
  namespace: cert-manager
spec:
  selfSigned: {}
EOF
```

### Bootstrapping CA `Issuers`

```shell
# 1. certmanager 生成 ca 证书
cat > selfsigned-ca-certificate.yaml << EOF
# Issue a root Certificate
# ca-example-com.certificate.cert-manager.yaml
# 参考：https://cert-manager.io/docs/usage/certificate/
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: selfsigned-ca
  namespace: cert-manager
spec:
  secretName: selfsigned-ca-tls ### Secret名字
  duration: 2160h
  renewBefore: 360h
  subject:
    organizations:
    - Example Inc. ###
  isCA: true ### 修改为true，isCA将将此证书标记为对证书签名有效。这会将cert sign自动添加到usages列表中。
  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 4096
  dnsNames:
  - ca.example.com
  # Issuer references are always required.
  issuerRef:
    name: selfsigned-issuer ### 指定为自签名发行人
    # We can reference ClusterIssuers by changing the kind here.
    # The default value is Issuer (i.e. a locally namespaced Issuer)
    kind: Issuer
    # This is optional since cert-manager will default to this value however
    # if you are using an external issuer, change this to that issuer group.
    group: cert-manager.io
EOF


cat > selfsigned-ca-issuer.yaml << EOF
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-ca-issuer
spec:
  ca:
    secretName: selfsigned-ca-tls
EOF
```

### Sign a site `Certificate` with ca issuer

```shell
cat > signed-site-tls.yaml << EOF
# 参考：https://cert-manager.io/docs/usage/certificate/
# api参考：https://cert-manager.io/docs/reference/api-docs/#cert-manager.io/v1alpha3.Certificate
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: gitlab-site
  namespace: gitlab-system
spec:
  # Secret names are always required.
  secretName: gitlab-site-tls ### Secret名字
  duration: 2160h # 90d
  renewBefore: 360h # 15d

  subject:
    organizations:
    - Example Inc.

  isCA: false
  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 2048

  # At least one of a DNS Name, URI, or IP address is required.
  dnsNames:
  - gitlab.gitops.com
  - kas.gitops.com
  - minio.gitops.com
  - registry.gitops.com
  - habor.gitops.com
  # uris: ### 注释了uris、ipAddresses
  # - spiffe://cluster.local/ns/sandbox/sa/example
  # ipAddresses:
  # - 192.168.0.5

  # Issuer references are always required.
  issuerRef:
    name: selfsigned-ca-issuer ### 使用CA Issuer
    kind: ClusterIssuer ### CA Issuer是ClusterIssuer
    # This is optional since cert-manager will default to this value however
    # if you are using an external issuer, change this to that issuer group.
    group: cert-manager.io
EOF
```

### Use in Ingress or Local

```shell
Ingress.spec:
  tls:
    - hosts:
        - site.example.com
      secretName: site-example-com-tls

# selfsigned-ca-tls.secret.cert-manager里的tls.crt文件，拷贝至开发机器上，windows直接打开安装证书至受信任的根证书颁发机构
```