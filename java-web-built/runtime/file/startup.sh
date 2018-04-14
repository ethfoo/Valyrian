#!/bin/bash

#最终运行tomcat的脚本
cd $CATALINA_HOME
$CATALINA_HOME/bin/catalina.sh run 2>&1 | tee /dev/stderr
