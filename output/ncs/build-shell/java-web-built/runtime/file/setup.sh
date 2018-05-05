#!/bin/bash

#为最终生成的容器镜像做初始化工作
# 
set -e
rm -rf $CATALINA_HOME/webapps/*
mv /ROOT.war $CATALINA_HOME/webapps/
mv /setenv.sh $CATALINA_HOME/bin/

rm /setup.sh
