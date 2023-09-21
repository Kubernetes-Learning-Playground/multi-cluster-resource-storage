## kubernetes简易多集群方案

### 项目思路与功能
项目背景：在目前云原生中，会常有需要同时操作"多集群"的场景，不论是多集群"查询"或是"分发资源"等操作，本项目采用 **informer** + **operator** 进行扩展封装，
实现**多集群**且**多资源**方案。

支持功能：
1. 支持"多集群"配置
2. 支持"多资源"配置
3. 支持跳过 restconfig tls 认证
4. 实现 http server 支持查询接口
5. 支持查询多集群命令行插件(list,describe)
6. 支持多集群下发资源

### 配置文件
- **重要** 配置文件可参考config.yaml中配置，调用方只需要关注配置文件中的内容即可。
```yaml
clusters:                     # 集群列表
  - metadata:
      clusterName: cluster1   # 自定义集群名
      insecure: false          # 是否开启跳过tls证书认证
      configPath: /Users/zhenyu.jiang/.kube/config # kube config配置文件地址
      # 资源类型
      resources:
        - rType: apps/v1/deployments
        - rType: core/v1/pods
        - rType: core/v1/configmaps
  - metadata:
      clusterName: cluster2   # 自定义集群名
      insecure: true          # 是否开启跳过tls证书认证
      configPath: /Users/zhenyu.jiang/go/src/golanglearning/new_project/multi_resource/resource/config1 # kube config配置文件地址
      resources:
        - rType: apps/v1/deployments
        - rType: core/v1/pods
        - rType: core/v1/configmaps
```
![](https://github.com/Kubernetes-Learning-Playground/multi-cluster-resource-storage/blob/main/image/%E6%97%A0%E6%A0%87%E9%A2%98-2023-08-10-2343.png?raw=true)

### 多集群命令行查询(也支持 http server查询)
目前支持查询资源
- pods
- configmaps
- deployments

后缀参数：
- --namespace：按命名空间查询，不填默认所有命名空间
- --clusterName：按集群名查询，不填默认所有集群
- --name: 按名称查询，不填默认所有名称
```bash
➜  cmd git:(main) ✗ go run ctl_plugin/main.go list configmaps --clusterName=cluster2      
集群名称          NAME                                   NAMESPACE               DATA 
cluster2        test-scheduling-config                  kube-system             1       
cluster2        loki-loki-stack-test                    loki-stack              1       
cluster2        kube-root-ca.crt                        loki-stack              1       
cluster2        loki-loki-stack                         loki-stack              1       
cluster2        kube-root-ca.crt                        etcd01                  1       
cluster2        kube-root-ca.crt                        mycsi                   1  

➜  cmd git:(main) ✗ go run ctl_plugin/main.go configmaps --clusterName=cluster2 --name=coredns --namespace=kube-system       
集群名称        CONFIGMAP       NAMESPACE       DATA 
cluster2        coredns         kube-system     1       
```
查询多集群 pods 资源
```bash
➜  cmd git:(main) ✗ go run ctl_plugin/main.go list pods --clusterName=cluster2                                   
集群名称         NAME                                                    NAMESPACE               POD IP          状态             容器名                           容器静像                                                                        
cluster2        virtual-kubelet-pod-test-bash                           default                                 Running         ngx1                            nginx:1.18-alpine                                                                    
cluster2        testpod1                                                default                                 Running         mytest                          nginx:1.18-alpine                                                                    
cluster2        loki-promtail-zxpvg                                     loki-stack                              Running         promtail                        docker.io/grafana/promtail:2.7.4                                                     
cluster2        node-exporter-srqk4                                     prometheus                              Running         node-exporter                   bitnami/node-exporter:1.4.0                                                          
cluster2        node-exporter-m5whb                                     prometheus                              Running         node-exporter                   bitnami/node-exporter:1.4.0                                                          
cluster2        loki-promtail-fcpsb                                     loki-stack                              Running         promtail                        docker.io/grafana/promtail:2.7.4                                                     
cluster2        testpod                                                 default                                 Pending         mytest                          nginx:1.18-alpine                                                                    
cluster2        nginx-kubelet                                           default                                 Running         nginx                           nginx:1.18-alpine                                                                    
cluster2        dep-test-8b4fcc97-pzbqd                                 default                 10.244.0.124    Running         dep-test-container              nginx:1.18-alpine                                                                    
cluster2        dep-test-8b4fcc97-jkkx7                                 default                 10.244.0.127    Running         dep-test-container              nginx:1.18-alpine                                                                    
cluster2        dep-test-8b4fcc97-wl6td                                 default                 10.244.0.128    Running         dep-test-container              nginx:1.18-alpine                                               

# 不指定clusterName，默认查询所有集群
➜  multi_resource git:(main) ✗ go run cmd/ctl_plugin/main.go list pods                           
集群名称         NAME                                                            NAMESPACE                               NODE                    POD IP          状态             容器名                        容器静像                                                                            
cluster1        patch-deployment-7877dfff-975bn                                 default                                 minikube                10.244.1.40     Running         nginx                        nginx:1.15.2                                                                            
cluster1        patch-deployment-7877dfff-dwpxj                                 default                                 minikube                10.244.1.39     Running         nginx                        nginx:1.15.2                                                                            
cluster2        virtual-kubelet-pod-test-bash                                   default                                 mynode                                  Running         ngx1                         nginx:1.18-alpine                                                                       
cluster1        kueue-controller-manager-56987d8f8c-69gr7                       kueue-system                            minikube                10.244.1.16     Running         manager                      registry.k8s.io/kueue/kueue:v0.4.1                                                      
cluster2        testpod1                                                        default                                 my-sample-kubelet                       Running         mytest                       nginx:1.18-alpine                                                                       
cluster2        loki-promtail-zxpvg                                             loki-stack                              my-sample-kubelet                       Running         promtail                     docker.io/grafana/promtail:2.7.4                                                        
cluster2        node-exporter-srqk4                                             prometheus                              my-sample-kubelet                       Running         node-exporter                bitnami/node-exporter:1.4.0                                                             
cluster2        node-exporter-m5whb                                             prometheus                              myk8s                                   Running         node-exporter                bitnami/node-exporter:1.4.0                                                             
cluster2        loki-promtail-fcpsb                                             loki-stack                              myk8s                                   Running         promtail                     docker.io/grafana/promtail:2.7.4                                                        
cluster2        testpod                                                         default                                 myjtthink                               Pending         mytest                       nginx:1.18-alpine                                                                       
cluster2        nginx-kubelet                                                   default                                 myjtthink                               Running         nginx                        nginx:1.18-alpine                                                                       
cluster2        dep-test-8b4fcc97-pzbqd                                         default                                 vm-0-16-centos          10.244.0.124    Running         dep-test-container           nginx:1.18-alpine                                                                       
cluster2        dep-test-8b4fcc97-jkkx7                                         default                                 vm-0-16-centos          10.244.0.127    Running         dep-test-container           nginx:1.18-alpine                                                                       
cluster2        dep-test-8b4fcc97-wl6td                                         default                                 vm-0-16-centos          10.244.0.128    Running         dep-test-container           nginx:1.18-alpine                                                                       
cluster2        dep-test-8b4fcc97-znlp5                                         default                                 vm-0-16-centos          10.244.0.125    Running         dep-test-container           nginx:1.18-alpine                                                                       
cluster2        dep-test-8b4fcc97-vxf55                                         default                                 vm-0-16-centos          10.244.0.126    Running         dep-test-container           nginx:1.18-alpine                                                                       
cluster2        inspect-script-task-task3--1-fjxm9                              default                                 vm-0-16-centos          10.244.0.94     Pending         default                      inspect-operator/script-engine:v1
```
查询多集群 deployments 资源
```bash
➜  cmd git:(main) ✗ go run ctl_plugin/main.go list deployments --clusterName=cluster2
集群名称         NAME                                    NAMESPACE               TOTAL   AVAILABLE       READY 
cluster2        dep-test                                default                 5       5               5       
cluster2        testngx                                 default                 10      10              10      
cluster2        test-pod-maxnum-scheduler               kube-system             1       1               1       
cluster2        myingress-controller                    default                 1       1               1       
cluster2        myapi                                   default                 1       1               1       
```
查询多集群 pods 资源详细
```bash
➜  multi_resource git:(main) ✗ go run cmd/ctl_plugin/main.go describe pods --clusterName=cluster2 --namespace=default --name=myredis-0
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2023-01-18T15:14:48Z"
  managedFields:
  - apiVersion: v1
```

### 多集群下发资源
- crd 资源对象如下
```yaml
apiVersion: mulitcluster.practice.com/v1alpha1
kind: MultiClusterResource
metadata:
  name: mypod.pod
  namespace: default
spec:
   # 资源模版，内部填写需要的 k8s 原始资源
   template:
     apiVersion: v1
     kind: Pod
     metadata:
       name: multicluster-pod
       namespace: default
     spec:
       containers:
         - image: busybox
           command:
             - sleep
             - "3600"
           imagePullPolicy: IfNotPresent
           name: busybox
       restartPolicy: Always
   # 可以自行选择不同集群下发，如果修改后，
   # 也会相应的新增或删除特定集群的资源    
   placement:
     clusters:
       - name: cluster1
       - name: cluster2
       - name: cluster3
```

- 使用
可以看出，当在主集群创建 CRD 后，会自动下发到其他集群。
```bash
# apply
➜  multi_resource git:(main) ✗ kubectl apply -f yaml/test.yaml    
multiclusterresource.mulitcluster.practice.com/mypod.pod created
# 查询
➜  multi_resource git:(main) ✗ kubectl get multiclusterresources.mulitcluster.practice.com    
NAME        AGE
mypod.pod   45m
➜  multi_resource git:(main) ✗ go run cmd/ctl_plugin/main.go list pods  --namespace=default --name=multicluster-pod
集群名称        NAME                    NAMESPACE       NODE            POD IP          状态    容器名  容器静像 
cluster3        multicluster-pod        default         vm-0-17-centos  10.244.167.193  Running busybox busybox         
cluster1        multicluster-pod        default         minikube        10.244.1.48     Running busybox busybox         
cluster2        multicluster-pod        default         vm-0-16-centos  10.244.0.142    Running busybox busybox         
```
