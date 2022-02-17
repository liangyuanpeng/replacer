# replacer

Replace container image of kubernetes mutate webhook

# 说明 

这个库是一个k8s的mutate webhook参考实现,目前唯一的作用是将k8s.gcr.io的镜像修改为lank8s.cn,后续会支持配置.

例如将xxxalixx/google_container修改为lank8s.cn, 以及自定义配置将aaa镜像修改为bbb.  

# 开发webhook

看其他博客说开发webhook时只能是写代码然后编译部署到K8S,如果有问题只能写完代码再编译放到K8S里面去验证.  

其实我们只需要部署一个webhook的转发就可以了,例如将envoy作为webhook部署在K8S里面,然后将所有请求都转发到你自己的开发环境中.  

# 部署  

目前可以直接使用仓库中deploy文件夹的内容.  

```shell
git clone git@github.com:liangyuanpeng/replacer.git
kubectl apply -f replacer/deploy
```  

马上会推出helm版本的部署方式.
