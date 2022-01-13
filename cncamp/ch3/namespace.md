### 在新 network namespace 执行 sleep 指令：

```sh
unshare -fn sleep 60

#### namespace操作：
#### clone（创建进程时，flags指定是否加入其它ns）
#### PROC.setns（调用进程PROC加入到已有namespace）
#### PROC.unshare（调用进程PROC加入新创建一个namespace）
```


### 查看进程信息

```sh
ps -ef|grep sleep
root       32882    4935  0 10:00 pts/0    00:00:00 unshare -fn sleep 60
root       32883   32882  0 10:00 pts/0    00:00:00 sleep 60
```

### 查看网络 Namespace

```sh
lsns -t net
4026532508 net       2 32882 root unassigned                                unshare

# -p, --task pid: Display only the namespaces held by the process with this pid.
# -r, --raw: Use the raw output format.
# -t, --type type: Display  the specified type of namespaces only.  The supported types are mnt, net, ipc, user, pid, uts and cgroup.  This option may be given more than once.

```

### 进入改进程所在 Namespace 查看网络配置，与主机不一致

```sh
nsenter -t 32882 -n ip a
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN group default qlen 1000
link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00

-t, --target pid: Specify a target process to get contexts from.  \
    The paths to the contexts specified by pid are:
# /proc/$pid/ns/mnt    the mount namespace
# /proc/$pid/ns/uts    the UTS namespace
# /proc/$pid/ns/ipc    the IPC namespace
# /proc/$pid/ns/net    the network namespace
# /proc/$pid/ns/pid    the PID namespace
# /proc/$pid/ns/user   the user namespace
# /proc/$pid/ns/cgroup the cgroup namespace
# /proc/$pid/root      the root directory
# /proc/$pid/cwd       the working directory respectively

-m, --mount[=file]: Enter the mount namespace.
-u, --uts[=file]: Enter  the UTS namespace.
-i, --ipc[=file]: Enter the IPC namespace.
-n, --net[=file]: Enter  the  network namespace. 
-p, --pid[=file]: Enter the PID namespace.
-U, --user[=file]: Enter  the  user  namespace.
-C, --cgroup[=file]: Enter the cgroup namespace.
```
