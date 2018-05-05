#!/bin/bash
docker run -v ~/.m2/repository:/root/.m2/repository -v /var/run/docker.sock:/var/run/docker.sock -v $(which docker):/usr/bin/docker valyrian:buildv1 -r git@github.com:ethfoo/base-java-web.git -b master -s -i ethoojava
#-v /var/log/valyrian:/var/log/valyrian
#-v /root/.ssh:/root/.ssh 
