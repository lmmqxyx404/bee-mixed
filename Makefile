SHELL = /bin/bash

#SCRIPT_DIR         = $(shell pwd)/etc/script
#请选择golang版本
BUILD_IMAGE_SERVER  = golang:1.22
#请选择node版本
BUILD_IMAGE_WEB     = node:18
#项目名称
PROJECT_NAME        = github.com/flipped-aurora/gin-vue-admin/server
#配置文件目录
CONFIG_FILE         = config.yaml
#镜像仓库命名空间
IMAGE_NAME          = gva
#镜像地址
REPOSITORY          = registry.cn-hangzhou.aliyuncs.com/${IMAGE_NAME}

# 如果用户没有在命令行通过 make TAGS_OPT=xxx 指定版本标签，则默认 TAGS_OPT=latest。
ifeq ($(TAGS_OPT),)
TAGS_OPT            = latest
else
endif

#容器环境前后端共同打包
# todo: 补充前端 build-web
build:  build-server

#容器环境打包前端
build-web:
	docker run --name build-web-local --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE_WEB} make build-web-local

#容器环境打包后端
build-server:
	docker run --name build-server-local --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE_SERVER} make build-server-local

#本地环境打包前端
build-web-local:
	@cd web/ && if [ -d "dist" ];then rm -rf dist; else echo "OK!"; fi \
	&& yarn install && yarn build
# 不设置国内代理
# && yarn config set registry http://mirrors.cloud.tencent.com/npm/ 


#本地环境打包后端
build-server-local:
	@cd server/ && if [ -f "server" ];then rm -rf server; else echo "OK build-server-local!"; fi \
	&& go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct \
	&& go env -w CGO_ENABLED=0 && go env  && go mod tidy \
	&& go build -buildvcs=false -ldflags "-B 0x$(shell head -c20 /dev/urandom|od -An -tx1|tr -d ' \n') -X main.Version=${TAGS_OPT}" -v