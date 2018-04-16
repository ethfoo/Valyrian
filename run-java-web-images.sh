#!/bin/bash
docker run -v /var/log/valyrian:/var/log/valyrian -v /root/.ssh:/root/.ssh -v ~/.m2/repository:/root/.m2/repository -v /var/run/docker.sock:/var/run/docker.sock -v $(which docker):/usr/bin/docker java-web-build:v1
