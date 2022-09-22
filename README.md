# gocron - 定时任务管理系统
[![Downloads](https://img.shields.io/github/downloads/peng49/gocron/total.svg)](https://github.com/peng49/gocron/releases)
[![license](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/peng49/gocron/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/peng49/gocron.svg?label=Release)](https://github.com/peng49/gocron/releases)

# 项目简介

本项目基于 [ouqiang/gocron](https://github.com/ouqiang/gocron),在原有的定时任务管理的基础上,新增的进程管理模块,LDAP用户认证功能,新增了项目管理功能用来对主机和任务进行分组管理,并使用 Element Plus + Vue3 重构了前端页面

### 部分截图
![dashboard.jpg](https://s2.loli.net/2022/09/22/itsfQZVYM3BnwK5.jpg)
![tasks.jpg](https://s2.loli.net/2022/09/22/KkltuZFbRfrD1ic.jpg)
![process.jpg](https://s2.loli.net/2022/09/22/RzlSo3YVyQAgCuf.jpg)
![ldap.jpg](https://s2.loli.net/2022/09/22/k4ctQJr3nZ6YFxm.jpg)

### 测试地址
[测试地址](https://gocron-test.fly-develop.com) 

用户名：admin
密码: admin123


    
### 支持平台
> Windows、Linux、Mac OS

### 环境要求
>  MySQL


## 下载
[releases](https://github.com/peng49/gocron/releases)  

[版本升级](https://github.com/peng49/gocron/wiki/版本升级)

## 安装

###  二进制安装
1. 解压压缩包
2. `cd 解压目录`   
3. 启动        
* 调度器启动        
  * Windows: `gocron.exe web`   
  * Linux、Mac OS:  `./gocron web`
* 任务节点启动, 默认监听0.0.0.0:5921
  * Windows:  `gocron-node.exe`
  * Linux、Mac OS:  `./gocron-node`
4. 浏览器访问 http://localhost:5920

### 源码安装

- 安装Go 1.18+
- `go get -d github.com/peng49/gocron`
- `export GO111MODULE=on` 
- 编译 `make`
- 启动
    * gocron `./bin/gocron web`
    * gocron-node `./bin/gocron-node`


### docker

```shell
docker run --name gocron --link mysql:db -p 5920:5920 -d peng49/gocron
```

配置: /app/conf/app.ini

日志: /app/log/cron.log

镜像不包含gocron-node, gocron-node需要和具体业务一起构建


### 开发

1. 安装Go1.18+, Node.js, npm
2. 安装前端依赖 `make install-vue`
3. 启动gocron, gocron-node `make run`
4. 启动node server `make run-vue`, 访问地址 http://localhost:8080

访问http://localhost:8080, API请求会转发给gocron

`make` 编译

`make run` 编译并运行

`make package` 打包 
> 生成当前系统的压缩包 gocron-v2.x.x-darwin-amd64.tar.gz gocron-node-v2.x.x-darwin-amd64.tar.gz

`make package-all` 生成Windows、Linux、Mac的压缩包

### 命令

* gocron
    * -v 查看版本

* gocron web
    * --host 默认0.0.0.0
    * -p 端口, 指定端口, 默认5920
    * -e 指定运行环境, dev|test|prod, dev模式下可查看更多日志信息, 默认prod
    * -h 查看帮助
* gocron-node
    * -allow-root *nix平台允许以root用户运行
    * -s ip:port 监听地址  
    * -enable-tls 开启TLS    
    * -ca-file CA证书文件 
    * -cert-file 证书文件  
    * -key-file  私钥文件
    * -h 查看帮助
    * -v 查看版本

## 程序使用的组件
* Web框架 [Macaron](http://go-macaron.com/)
* 定时任务调度 [Cron](https://github.com/robfig/cron)
* ORM [Xorm](https://github.com/go-xorm/xorm)
* UI框架 [Element Plus](https://github.com/element-plus/element-plus)
* 依赖管理 [Govendor](https://github.com/kardianos/govendor)
* RPC框架 [gRPC](https://github.com/grpc/grpc)

## 反馈
提交[issue](https://github.com/peng49/gocron/issues/new)

## ChangeLog

v2.0.0
--------
+ LDAP用户认证
* 添加项目管理，项目和主机,任务关联
* 进程管理(队列消费程序)
* Vue3+ElementPlus 重构前端页面


[ouqiang/gocron ChangeLog](https://github.com/ouqiang/gocron#changelog)