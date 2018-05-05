#!/bin/bash
#构建 构建使用的镜像
cd output/nce/build-shell/java-web-built/
docker build -t valyrian:buildv1 .