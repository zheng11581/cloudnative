## Scale-up ceph osd

### Check OSD status

```shell
# k -n rook-ceph exec -it rook-ceph-tools-78f6fdf966-v9tm2 -- ceph osd status
ID  HOST    USED  AVAIL  WR OPS  WR DATA  RD OPS  RD DATA  STATE      
 0  node1  21.1M  99.9G      0        0       0        0   exists,up  
 1  node2  22.3M  99.9G      0        0       0        0   exists,up 
```

### Check OSD tree

```shell
# k -n rook-ceph exec -it rook-ceph-tools-78f6fdf966-v9tm2 -- ceph osd tree
ID  CLASS  WEIGHT   TYPE NAME       STATUS  REWEIGHT  PRI-AFF
-1         0.03897  root default                             
-3         0.01949      host node1                           
 0    hdd  0.01949          osd.0       up   1.00000  1.00000
-5         0.01949      host node2                           
 1    hdd  0.01949          osd.1       up   1.00000  1.00000

```

### Add Disk to the node

```shell
# partprobe

```

### Extend the LV of ceph 

```shell
# pvresize -t /dev/sdb
# pvresize /dev/sdb
# pvs
# lvdisplay 
# lvextend -l 100%VG /dev/ceph-70ab501d-33fb-41de-b43c-9abc0cbcbab0/osd-block-3f196f8d-26d1-42e6-b367-0de3c791e30d
# lvdisplay 

```

### Make rook-ceph reload the OSD

```shell
# k -n rook-ceph exec -it rook-ceph-tools-78f6fdf966-v9tm2 -- ceph osd df
# k -n rook-ceph get deployments.apps rook-ceph-osd-0 -oyaml > rook-ceph-osd-0-deploy.yaml
# k -n rook-ceph edit deployments.apps rook-ceph-osd-0
# k -n rook-ceph get pod

      containers:
      - args:
        - sleep 3600
        command:
        - sh
        - -c

# k -n rook-ceph exec -it rook-ceph-osd-0-67fdd7d9fb-mh7vt -- sh

# ceph-bluestore-tool bluefs-bdev-sizes --path /var/lib/ceph/osd/ceph-0/
# ceph-bluestore-tool bluefs-bdev-expand --path /var/lib/ceph/osd/ceph-0
# ceph-bluestore-tool bluefs-bdev-sizes --path /var/lib/ceph/osd/ceph-0/

# vim rook-ceph-osd-0-deploy.yaml 
delete resourceVersion and uid

# k apply -f rook-ceph-osd-0-deploy.yaml -n rook-ceph 
```

