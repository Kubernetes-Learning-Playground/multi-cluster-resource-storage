## kubernetes simple multi-cluster solution
<a href="./README.md">English</a> | <a href="./README-zh.md">简体中文</a>
### Introduction
Project background: In the current cloud native architecture, there are often scenarios where "multiple clusters" need to be operated at the same time, whether it is multi-cluster "query" or "resource distribution" operations. This project uses **informer** + **operator** performs extended encapsulation,
Implement **multi-cluster** and **multi-resource** solutions.

Supported functions:
1. Support "multi-cluster" configuration
2. Support "multi-resource" configuration
3. Support skipping restconfig tls authentication
4. Implement http server support query interface
5. Support querying multi-cluster command line plug-ins (list, describe)
6. Support multiple clusters to deliver resources
7. Support multi-cluster **differentiated configuration**

### Configuration file
- **Important** The configuration file can refer to the configuration in config.yaml. [here](./config.yaml)
- The caller only needs to pay attention to the content in the configuration file.
```yaml
clusters:                     # Cluster list
  - metadata:
      clusterName: tencent1   # Custom cluster name
      insecure: false         # Whether to enable skipping tls certificate authentication
      configPath: /Users/zhenyu.jiang/.kube/config # kubeconfig configuration file address
      # Resource Type
      resources:
        - rType: apps/v1/deployments
        - rType: core/v1/pods
        - rType: core/v1/configmaps
  - metadata:
      clusterName: tencent2
      insecure: true
      configPath: /Users/zhenyu.jiang/go/src/golanglearning/new_project/multi_resource/multiclusterresource/config1 # kube config配置文件地址
      resources:
        - rType: apps/v1/deployments
        - rType: core/v1/pods
        - rType: core/v1/configmaps
```
![](https://github.com/Kubernetes-Learning-Playground/multi-cluster-resource-storage/blob/main/image/%E6%97%A0%E6%A0%87%E9%A2%98-2023-08-10-2343.png?raw=true)

### Multi-cluster ctl query (also supports http server query)
Support query resources:
- pods
- configmaps
- deployments

Suffix parameter:
- --namespace：Query by namespace. If not filled in, all namespaces will be queried by default.
- --clusterName：Query by cluster name. If not filled in, all clusters will be queried by default.
- --name: Query by name, if not filled in, query all by default
```bash
➜  cmd git:(main) ✗ go run ctl_plugin/main.go list configmaps --clusterName=tencent2      
集群名称          NAME                                   NAMESPACE               DATA 
tencent2        test-scheduling-config                  kube-system             1       
tencent2        loki-loki-stack-test                    loki-stack              1       
tencent2        kube-root-ca.crt                        loki-stack              1       
tencent2        loki-loki-stack                         loki-stack              1       
tencent2        kube-root-ca.crt                        etcd01                  1       
tencent2        kube-root-ca.crt                        mycsi                   1  

➜  cmd git:(main) ✗ go run ctl_plugin/main.go configmaps --clusterName=tencent2 --name=coredns --namespace=kube-system       
集群名称        CONFIGMAP       NAMESPACE       DATA 
tencent2        coredns         kube-system     1       
```
Query multi-cluster pods resources
```bash
➜  cmd git:(main) ✗ go run ctl_plugin/main.go list pods --clusterName=tencent2                                   
集群名称         NAME                                                    NAMESPACE               POD IP          状态             容器名                           容器静像                                                                        
tencent2        virtual-kubelet-pod-test-bash                           default                                 Running         ngx1                            nginx:1.18-alpine                                                                    
tencent2        testpod1                                                default                                 Running         mytest                          nginx:1.18-alpine                                                                    
tencent2        loki-promtail-zxpvg                                     loki-stack                              Running         promtail                        docker.io/grafana/promtail:2.7.4                                                     
tencent2        node-exporter-srqk4                                     prometheus                              Running         node-exporter                   bitnami/node-exporter:1.4.0                                                          
tencent2        node-exporter-m5whb                                     prometheus                              Running         node-exporter                   bitnami/node-exporter:1.4.0                                                          
tencent2        loki-promtail-fcpsb                                     loki-stack                              Running         promtail                        docker.io/grafana/promtail:2.7.4                                                     
tencent2        testpod                                                 default                                 Pending         mytest                          nginx:1.18-alpine                                                                    
tencent2        nginx-kubelet                                           default                                 Running         nginx                           nginx:1.18-alpine                                                                    
tencent2        dep-test-8b4fcc97-pzbqd                                 default                 10.244.0.124    Running         dep-test-container              nginx:1.18-alpine                                                                    
tencent2        dep-test-8b4fcc97-jkkx7                                 default                 10.244.0.127    Running         dep-test-container              nginx:1.18-alpine                                                                    
tencent2        dep-test-8b4fcc97-wl6td                                 default                 10.244.0.128    Running         dep-test-container              nginx:1.18-alpine                                               

# If clusterName is not specified, all clusters will be queried by default.
➜  multi_resource git:(main) ✗ go run cmd/ctl_plugin/main.go list pods                           
集群名称         NAME                                                            NAMESPACE                               NODE                    POD IP          状态             容器名                        容器静像                                                                            
tencent1        patch-deployment-7877dfff-975bn                                 default                                 minikube                10.244.1.40     Running         nginx                        nginx:1.15.2                                                                            
tencent1        patch-deployment-7877dfff-dwpxj                                 default                                 minikube                10.244.1.39     Running         nginx                        nginx:1.15.2                                                                            
tencent2        virtual-kubelet-pod-test-bash                                   default                                 mynode                                  Running         ngx1                         nginx:1.18-alpine                                                                       
tencent1        kueue-controller-manager-56987d8f8c-69gr7                       kueue-system                            minikube                10.244.1.16     Running         manager                      registry.k8s.io/kueue/kueue:v0.4.1                                                      
tencent2        testpod1                                                        default                                 my-sample-kubelet                       Running         mytest                       nginx:1.18-alpine                                                                       
tencent2        loki-promtail-zxpvg                                             loki-stack                              my-sample-kubelet                       Running         promtail                     docker.io/grafana/promtail:2.7.4                                                        
tencent2        node-exporter-srqk4                                             prometheus                              my-sample-kubelet                       Running         node-exporter                bitnami/node-exporter:1.4.0                                                             
tencent2        node-exporter-m5whb                                             prometheus                              myk8s                                   Running         node-exporter                bitnami/node-exporter:1.4.0                                                             
tencent2        loki-promtail-fcpsb                                             loki-stack                              myk8s                                   Running         promtail                     docker.io/grafana/promtail:2.7.4                                                        
tencent2        testpod                                                         default                                 myjtthink                               Pending         mytest                       nginx:1.18-alpine                                                                       
tencent2        nginx-kubelet                                                   default                                 myjtthink                               Running         nginx                        nginx:1.18-alpine                                                                       
tencent2        dep-test-8b4fcc97-pzbqd                                         default                                 vm-0-16-centos          10.244.0.124    Running         dep-test-container           nginx:1.18-alpine                                                                       
tencent2        dep-test-8b4fcc97-jkkx7                                         default                                 vm-0-16-centos          10.244.0.127    Running         dep-test-container           nginx:1.18-alpine                                                                       
tencent2        dep-test-8b4fcc97-wl6td                                         default                                 vm-0-16-centos          10.244.0.128    Running         dep-test-container           nginx:1.18-alpine                                                                       
tencent2        dep-test-8b4fcc97-znlp5                                         default                                 vm-0-16-centos          10.244.0.125    Running         dep-test-container           nginx:1.18-alpine                                                                       
tencent2        dep-test-8b4fcc97-vxf55                                         default                                 vm-0-16-centos          10.244.0.126    Running         dep-test-container           nginx:1.18-alpine                                                                       
tencent2        inspect-script-task-task3--1-fjxm9                              default                                 vm-0-16-centos          10.244.0.94     Pending         default                      inspect-operator/script-engine:v1
```
Query multi-cluster deployments resources
```bash
➜  cmd git:(main) ✗ go run ctl_plugin/main.go list deployments --clusterName=cluster2
集群名称         NAME                                    NAMESPACE               TOTAL   AVAILABLE       READY 
tencent1        dep-test                                default                 5       5               5       
tencent1        testngx                                 default                 10      10              10      
tencent1        test-pod-maxnum-scheduler               kube-system             1       1               1       
tencent1        myingress-controller                    default                 1       1               1       
tencent1        myapi                                   default                 1       1               1       
```
Query multi-cluster pods resource details
```bash
➜  multi_resource git:(main) ✗ go run cmd/ctl_plugin/main.go describe pods --clusterName=tencent2 --namespace=default --name=myredis-0
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: "2023-01-18T15:14:48Z"
  managedFields:
  - apiVersion: v1
```

#### Deploy Ctl tool
```bash
➜  ctl_plugin git:(main) ✗ pwd
/xxxxx/multi_resource/cmd/ctl_plugin
➜  ctl_plugin git:(main) ✗ go build -o kubectl-multicluster .
➜  ctl_plugin git:(main) ✗ chmod 777 kubectl-multicluster                     
➜  ctl_plugin git:(main) ✗ mv kubectl-multicluster ~/go/bin/ 
➜  ~ kubectl-multicluster list pods --clusterName=tencent1 --name=multiclusterresource-deployment-75d98bb7bd-xj5z5
集群名称        NAME                                                    NAMESPACE       NODE            POD IP          状态    容器名                  容器镜像            
tencent1        multiclusterresource-deployment-75d98bb7bd-xj5z5        default         vm-0-12-centos  10.244.29.32    Running example-container       nginx:1.19.0-alpine   
```

Note: The default command line of this project will read the ~/.multi-cluster-operator/config configuration, which will be automatically created when the project starts.
If there are related errors when using the command line, you can handle the relevant parts yourself.
```bash
➜  .multi-cluster-operator pwd
/Users/zhenyu.jiang/.multi-cluster-operator
➜  .multi-cluster-operator cat config                
server: 31888
```

### Delivering resources to multiple clusters
- The crd resource object is as follows. For more information, please refer to [Reference](./yaml)
```yaml
apiVersion: mulitcluster.practice.com/v1alpha1
kind: MultiClusterResource
metadata:
  name: mypod.pod
  namespace: default
spec:
   # Resource template, fill in the required k8s original resources
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
     # You can choose different clusters for delivery. If modified,
     # Resources for specific clusters will also be added or deleted accordingly.
     # Note: If placement is deleted, the corresponding cluster of customize also needs to be deleted! !
   placement:
     clusters:
       - name: tencent1
       - name: tencent2
       - name: tencent4
   # Can be blank
   customize:
     clusters:
       - name: tencent1
         action:
           # Replace image
           - path: "/spec/containers/0/image"
             op: "replace"
             value:
               - "nginx:1.19.0-alpine"
       - name: tencent2
         action:
           # add annotations
           - path: "/metadata/annotations/example"
             value:
               - "example"
             op: "add"
```
- usage

It can be seen that after a CRD is created in the main cluster, it will be automatically distributed to other clusters.
```bash
# apply
➜  multi_resource git:(main) ✗ kubectl apply -f yaml/test.yaml    
multiclusterresource.mulitcluster.practice.com/mypod.pod created
# query
➜  multi_resource git:(main) ✗ kubectl get multiclusterresources.mulitcluster.practice.com    
NAME        AGE
mypod.pod   45m
➜  multi_resource git:(main) ✗ go run cmd/ctl_plugin/main.go list pods  --namespace=default --name=multicluster-pod
集群名称        NAME                    NAMESPACE       NODE            POD IP          状态    容器名  容器静像 
tencent4        multicluster-pod        default         vm-0-17-centos  10.244.167.193  Running busybox busybox         
tencent1        multicluster-pod        default         minikube        10.244.1.48     Running busybox busybox         
tencent2        multicluster-pod        default         vm-0-16-centos  10.244.0.142    Running busybox busybox   
# Use the customize field to achieve differentiated deployment
➜  cmd git:(main) ✗ go run ctl_plugin/main.go list deployments --name=multiclusterresource-deployment
集群名称        NAME                            NAMESPACE       TOTAL   AVAILABLE       READY 
tencent4        multiclusterresource-deployment default         3       3               3       
tencent1        multiclusterresource-deployment default         2       2               2       
tencent2        multiclusterresource-deployment default         1       1               1       
```
