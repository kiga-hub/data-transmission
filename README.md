data transmission


github.com/kiga-hub/arc

# 项目文档

- 后台框架：echo web

## 1. 基本介绍

### 1.1 项目介绍

> deploy-tools 是自动化部署平台和应用的工具，部署工具使用Ansible，由后台直接调用

## 2. 使用说明

```
- golang版本 >= v1.22
- IDE推荐：Goland
```

### 2.1 server端

```bash
# 使用 go.mod

# 安装go依赖包
go list (go mod tidy)

# 编译
./build.sh deploy-tools deploy-tools
```

### 2.2 修改配置文件

```bash
vim deploy-tools.toml

# 单独运行，需要将inswarm配置为false
inswarm = false
```

### 2.3 启动

```bash
# 运行
./deploy-tools run
```

### 2.4 Docker安装

```bash
# 创建docker镜像
make

# 启动
docker run -d --restart=always --name deploy-tools -p 8081:80 -v "./data/task:/data/task"  -v "./data/images:/data/images" image.aithu.com/app/deploy-tools:v1.0.25-linux-amd64
```

## 3. 技术选型

- 后端：用`echo`快速搭建基础restful风格API。
- API文档：使用`Swagger`构建自动化文档。
- 配置文件：使用`deploy-tools.toml`格式的配置文件。
- 部署工具：使用`ansible`。

## 4. 项目架构

### 4.1 目录结构

```
    ├─ansible           （部署工具）
    ├─cmd               （CLI命令行工具）
        ├─server.go     （服务初始化入口）
    ├─ui                （前端页面）
    ├─pkg               （依赖包）
        ├─api           （接口）
        ├─component     （组件接口函数实现）
        ├─service       （服务）
        └─ansible       （调用ansible）
```

### 4.2 数据目录结构

```
     ├─task                                                 （部署任务目录）
        ├─20240829044247                                    （任务时间目录，作为任务唯一标识）
            ├─ansible.log                                   （ansible日志文件）
            ├─config.yml                                    （ansible任务配置文件）
            ├─task.yml                                      （部署任务信息）
        ...
     ├─rpm                                                  （rpm安装包目录）
        ├─docker                                            （docker相关离线安装包目录）
            ├─docker-20.10.3.tgz                            （docker二进制离线安装包）
            ├─docker-compose-v2.20.0                        （docker-compose二进制程序）
            ├─docker-rootless-extras-v20.10.3               （docker二进制工具）
        ├─ntp                                               （ntp离线安装包目录）
            ├─ntp-4.2.6p5-29.el7.centos.2.x86_64.rpm
        ├─tools                                             （其他工具安装包目录）
            ├─kubectl-1.28.2-0.x86_64.rpm
            ├─nvidia-container-toolkit-1.16.1-1.x86_64.rpm
            ...
     ├─images                                               （离线镜像目录，压缩包命名规则：[服务名]-[版本号].tar）
        ├─rancher-v2.5.12.tar
        ├─zk-3.6.3.tar
        ├─mysql-5.7.tar
        ...
     ├─deployments                                          （部署yaml文件目录, 命名规则：[服务名].yaml）
        ├─zk.yaml
        ├─mysql.yaml
        ...
     ├─configs                                              （配置映射yaml文件目录）
        ├─redis-configmap.yaml
        ├─servermysql.conf.yaml
        ...
     ├─data                                                 （初始化数据目录）
        ├─mysql-db.tar.gz
        ├─clickhouse-db.tar.gz
        ├─iotdb-db.tar.gz
     ├─vars                                                 （部署依赖变量）
        ├─ucd.yml
        ├─iot.yml
        ├─data.yml
        ├─algo.yml
        ├─voiceprint.yml
        ├─app.yml
     ├─scripts                                              （工具脚本目录）
        ├─requirements.txt
        ├─rancher_service_tools.py

```