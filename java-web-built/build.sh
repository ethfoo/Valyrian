#!/bin/bash

#该脚本
# 1. 拉取代码
# 2. maven构建war包
# 3. 修改运行镜像的startup.sh中设置的环境变量，包括JVM参数，不同环境配置等
# 4. 版本控制，生成特定格式的镜像TAG，然后构建镜像

#####全局变量######
#使用的基础镜像中已有的环境变量：
#$JAVA_HOME=/docker-java-home
#$MAVEN_HOME=/usr/share/maven

#远程GIT代码仓库
REMOTE_GIT_ADDR=
#代码分支
BRANCH=
#构建的MAVEN子模块
TARGET_MODULE=
#镜像TAG前缀
IMAGE_PRE_NAME=
###################

######### git拉取代码 ############
cd /root
git clone --branch=$(BRANCH) $(REMOTE_GIT_ADDR) workapp
cd /root/workapp
_COMMIT_HASH=`git rev-parse HEAD`
echo 'git HEAD commit hash:' $(REV_HASH)

######### maven构建war包 #########
mvn -T 1.5C clean install -pl ${TARGET_MODULE} -am -Dmaven.test.skip=true
cd /root/workapp/${TARGET}/target
#将war包复制到file目录里，COPY到运行时镜像中
mv ./*.war builder/files/ROOT.war

####### 修改startup.sh脚本中的环境变量（JVM参数等均以环境变量的方式注入）#############
# 不同环境配置在startup.sh中以环境变量的方式设置，用户可以在启动容器的时候以传环境变量CONF=XX的方式动态配置，在此设置会置一个默认值

######### 镜像TAG版本控制 ########
_REV=${COMMIT_HASH:0:8}
_NOW=`date +%Y%m%d-%H%M%S`

IMAGE_FULL_NAME=$(IMAGE_PRE_NAME)$(TARGET_MODULE):$(_NOW)$(_REV)
######### 构建镜像 ###############
cd /root/builder
docker build -t $(IMAGE_FULL_NAME) .
