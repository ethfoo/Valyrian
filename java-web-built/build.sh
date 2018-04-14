#!/bin/bash

#该脚本
# 1. 拉取代码
# 2. maven构建war包

#####全局变量######
#远程GIT代码仓库
REMOTE_GIT_ADDR=
#代码分支
BRANCH=
#构建的MAVEN子模块
TARGET_MODULE=


# git拉取代码


# maven构建war包


# 修改startup.sh脚本中的环境变量（JVM参数等均以环境变量的方式注入）
# 不同环境配置在startup.sh中以环境变量的方式设置，用户可以在启动容器的时候以传环境变量CONF=XX的方式动态配置，在此设置会置一个默认值


# 构建镜像

