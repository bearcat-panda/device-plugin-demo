# device-plugin
k8s device-plugin demo.

* i-device-plugin 将会新增资源 `test.com/device`。
* 将会扫描获取 `/etc/demo` 目录下的文件作为对应的设备。
* 将设备分配给 Pod 后，会在 Pod 中新增环境变量`Demo=$deviceId`


## 构建镜像
```bash
make build-image
```

## 部署
使用 DaemonSet 来部署 `device-plugin`，以便其能运行到集群中的所有节点上。

```bash
kubectl apply -f deploy/daemonset.yaml
```
检测 Pod 运行情况

```bash
[root@master240 device-plugin]# kubectl get po -n kube-system|grep device
device-plugin-56qtt                        1/1     Running     0          23m
```

## 测试

新增设备,在该 Demo 中，把 /etc/demo 目录下的文件作为设备，因此我们只需要到 /etc/demo 目录下创建文件，模拟有新的设备接入即可。
```bash
mkdir /etc/demo

touch /etc/demo/g1
```
查看 device plugin pod 日志,可以正常感知到设备
```bash
I0814 02:41:14.564573       1 api.go:29] waiting for device update
I0814 02:41:49.497473       1 device_monitor.go:78] fsnotify device event: /etc/demo/g1 CREATE
I0814 02:41:49.497609       1 device_monitor.go:87] find new device [g1]
I0814 02:41:49.497618       1 device_monitor.go:78] fsnotify device event: /etc/demo/g1 CHMOD
I0814 02:41:49.497681       1 api.go:32] device update, new device list [g1]

```
查看 node capacity 信息，能够看到新增的资源
```bash
[root@master240 device-plugin]# kubectl get node master240.dongwu -oyaml|grep  capacity -A 7
  capacity:
    cpu: "16"
    ephemeral-storage: 85138832Ki
    hugepages-1Gi: "0"
    hugepages-2Mi: "0"
    memory: 31620396Ki
    pods: "110"
    test.com/device: "1"
```

创建 Pod 申请该资源
```bash
kubectl apply -f deploy/test-pod.yaml
```
Pod 启动成功

```bash
[root@master240 device-plugin]# kubectl get po
NAME   READY   STATUS    RESTARTS   AGE
demo   1/1     Running   0          19m
```

之前分配设备是添加 Demo=xxx  这个环境变量，现在看下是否正常分配

```bash
[root@master240 device-plugin]# kubectl exec -it demo -- env|grep Demo
Demo=g1
```

ok,环境变量存在，可以看到分配给该 Pod 的设备是 g1。