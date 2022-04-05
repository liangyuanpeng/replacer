# replacer

Replace container image of kubernetes mutate webhook

# 说明 

这个库是一个k8s的mutate webhook参考实现,目前唯一的作用是将`k8s.gcr.io`的镜像修改为`lank8s.cn`,将`gcr.io`替换成`gcr.lank8s.cn`,后续会支持配置.  

[lank8s.cn](https://github.com/lank8s)是目前国内唯一一个可以访问`k8s.gcr.io`和`gcr.io`的一个镜像代理.

例如将xxxalixx/google_container修改为lank8s.cn, 以及自定义配置将aaa镜像修改为bbb.  

# 开发webhook

看其他博客说开发webhook时只能是写代码然后编译部署到K8S,如果有问题只能写完代码再编译放到K8S里面去验证.  

其实我们只需要部署一个webhook的转发就可以了,例如将envoy作为webhook部署在K8S里面,然后将所有请求都转发到你自己的开发环境中.  

# 支持的工作负载

- [x] Deployment
- [x] StatefulSet
- [x] DaemonSet
- [x] Job
- [ ] CronJob  


# 部署  

## kubectl apply

目前可以直接使用仓库中deploy文件夹的内容.  

```shell
git clone git@github.com:liangyuanpeng/replacer.git
cd replacer
kubectl create namespace replacer
kubectl apply -f deploy -n replacer
```  

## Helm 

```
helm repo add lyp https://liangyuanpeng.github.io/charts
helm install replacer lyp/replacer -n replacer --create-namespace
```

# 查看部署情况 

```
kubectl get po -n replacer
```

# 测试镜像替换效果  

## 拉取代码仓库 

如果你在前面操作已经拉取过了那么不需要再次拉取
```
git clone git@github.com:liangyuanpeng/replacer.git
cd replacer
kubectl apply -f deploy/test/kube-proxy-deployment.yaml
```

测试文件中的镜像为`k8s.gcr.io/kube-proxy:v1.10.1`,如果pod都够正常启动并且你的网络无法访问`k8s.gcr.io`那么说明webhook已经在正常工作了,接下来无需为任何`k8s.gcr.io`或`gcr.io`镜像拉取问题而烦恼了!

祝你使用愉快!

# RoadMap

- [ ] 替换镜像支持可配置,(配置将仓库A的镜像替换成仓库B的镜像)
- [x] 支持Job和CronJob的工作负载.

