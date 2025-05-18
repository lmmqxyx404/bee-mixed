# todo
add the latest gva

# note
## 后端是四个项目
bee-api
这个项目要单独部署
data-sdk
server
yunlaba-sdk

# tutorial
0. 本人是使用 `docker` 进行配置与开发的，推荐在 `linux` 环境下使用
1. `make run-dev` 开始配置开发环境
2. 注意第一次进入后台后，要设置数据库 `host` 字段时。要根据实际情况填写。
3. (deprecated) (没用了) bee-api 要使用 `config.yml.demo` 去配置
不需要专门的的去启动 api 服务，server 插件支持对应的服务, 需要暴露端口

# todo
## (done)维护地址信息的sql表 要正确导入
需要在对应的 `mysql` 容器内导入

## 部署生产代码步骤
### 前端
使用容器部署，要注意对应的打包步骤。而且真正使用的 nginx 配置是 `my.conf`