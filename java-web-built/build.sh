#!/bin/bash
set -e
set -x

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
#REMOTE_GIT_ADDR='git@github.com:ethfoo/base-java-web.git'
#REMOTE_GIT_ADDR='https://github.com/ethfoo/base-java-web.git'
#代码分支
#BRANCH='master'
#MAVEN项目类型：有子模块还是没有
#STANDALONE='-s'
#构建的MAVEN子模块
#TARGET_MODULE='base-springmvc-spring Webapp'
#镜像TAG前缀
#IMAGE_PRE_NAME='valyrian'
###################
#usage() { echo "Usage: $0 -r <remote git addr> -b <branch> -t <target module> -s <standalone> -i <image pre name>" 1>2&; exit 1; }
while getopts "r:b:t:i:s" o; do
	case "${o}" in
		r)
			REMOTE_GIT_ADDR=${OPTARG}
			;;
		b)
			BRANCH=${OPTARG}
			;;
		t)
			TARGET_MODULE=${OPTARG}
			;;
		s)
			STANDALONE="-s"
			;;
		i)
			IMAGE_PRE_NAME=${OPTARG}
			;;
		*)
			usage
			;;
	esac
done
shift $((OPTIND-1))


LOG='/var/log/valyrian'

######### git拉取代码 ############
cd /root
git clone --branch=${BRANCH} ${REMOTE_GIT_ADDR} workapp
cd /root/workapp
_COMMIT_HASH=`git rev-parse HEAD`

######### maven构建war包 #########
if [ "$STANDALONE" = "-s" ]; then
	mvn clean install -Dmaven.teset.skip=true
	cd /root/workapp/target
else
	mvn clean install -pl ${TARGET_MODULE} -am -Dmaven.test.skip=true
	cd /root/workapp/${TARGET_MODULE}/target
fi

#将war包复制到file目录里，COPY到运行时镜像中
mv ./*.war /root/builder/file/ROOT.war

####### 修改startup.sh脚本中的环境变量（JVM参数等均以环境变量的方式注入）#############
# 不同环境配置在startup.sh中以环境变量的方式设置，用户可以在启动容器的时候以传环境变量CONF=XX的方式动态配置，在此设置会置一个默认值

######### 镜像TAG版本控制 ########
_REV=${COMMIT_HASH:0:8}
_NOW=`date +%Y%m%d-%H%M%S`

_IMAGE_FULL_NAME="${IMAGE_PRE_NAME}:${_NOW}${_REV}"
######### 构建镜像 ###############
cd /root/builder
docker build -t ${_IMAGE_FULL_NAME} .
