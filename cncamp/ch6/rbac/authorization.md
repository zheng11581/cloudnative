### RBAC最佳实践
集群资源（cluster-resources）：Nodes PersistentVolumes Namespaces ResourceQuota等：ClusterRoleBinding ClusterRole
非资源URL（non-resources-url）: /healthz等url：ClusterRoleBinding ClusterRole
跨越 命名空间资源（namespace-resources）: 只需要多个命名空间创建独立的RoleBinding来绑定ClusterRole，这样只需要创建多个RoleBinding
指定 命令空间资源（namespace-resources）：需要独立的命名空间创建独立的RoleBinding来绑定Role，这样需要创建多个RoleBinding和Role
