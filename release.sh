#!/usr/bin/env bash

# 超简易部署脚本：支持三种参数
# 用法：
#   deploy.sh server [HOST]
#   deploy.sh dist [HOST] [LOCAL_DIR]
# 默认 HOST 为 aliyun，LOCAL_DIR 默认为 ./web/dist
# warning: 注意每个变量的用法
set -euo pipefail

MODE="${1:-}"
HOST="${2:-ali}"
# 修改成你自己的远程前端目录
REMOTE_DIR="/var/www/html"

if [[ -z "$MODE" ]]; then
  echo "Usage: $0 {server|dist} [HOST] [LOCAL_DIR]"
  exit 1
fi

case "$MODE" in
  server)
    # 本地打包
    make build-server
    # 之后上传到远程服务器
    echo "🔌 停止服务 bee-server.service on $HOST..."
    ssh root@"$HOST" "systemctl stop bee-server.service"

    echo "🚚 上传二进制 ./server/server -> /root"
    scp ./server/server root@"$HOST":/root/

    echo "🔄 启动服务 bee-server.service on $HOST..."
    ssh root@"$HOST" "systemctl start bee-server.service"
    ;;

  dist)
    # 本地打包
    make build-web
    # 之后上传到远程服务器
    LOCAL_DIR="${3:-./web/dist}"
    echo "📁 上传目录 $LOCAL_DIR -> \$REMOTE_DIR on $HOST..."
    scp -r "$LOCAL_DIR"/* root@"$HOST":"\$REMOTE_DIR/"
    ;;

  *)
    echo "Usage: $0 {server|dist} [HOST] [LOCAL_DIR]"
    exit 1
    ;;

esac

echo "✅ 完成模式 '$MODE' 的部署！"
