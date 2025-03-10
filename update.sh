#!/bin/bash

# 锁文件路径
LOCK_FILE="/var/run/hrss_update.lock"

# 创建锁文件并获取文件描述符 200
exec 200>$LOCK_FILE

# 尝试锁定文件描述符 200，若失败则输出错误并退出脚本
flock -n 200 || { echo "Another instance is running, exiting script"; exit 1; }

# 项目目录
cd /opt/1panel/apps/hrss/hrss || exit

# 运行更新程序
./hrss_update

# 获取最新版本号和当前版本号，并去除多余的空格和回车符
latest_version=$(tr -d '[:space:]' < latest_version.txt)
current_version=$(tr -d '[:space:]' < current_version.txt)

# 如果版本号不同，那么执行 Docker 容器的删除和更新操作
if [ "$latest_version" != "$current_version" ]; then
  # 构建镜像并启动容器
  docker-compose build

  # 检查 docker-compose 命令是否成功
  # shellcheck disable=SC2181
  if [ $? -ne 0 ]; then
    echo "docker-compose build failed, exiting script"
    exit 1
  fi

  docker-compose down
  docker-compose up -d

  # 前端
  website_dir="/opt/1panel/apps/openresty/openresty/www/sites/172.16.101.164/index/"
  # 设置文件权限(1panel)
  chown -R root:root "$website_dir"
  # 删除原有文件
  rm -rf "${website_dir:?}"*

  # 解压 dist.zip 到指定目录
  unzip ./download/dist.zip -d ./download/unzip/

  # 使用 rsync 代替 mv，移动 dist 目录下的所有文件到上一级目录，并覆盖同名文件
  rsync -a --remove-source-files ./download/unzip/dist/ "$website_dir"
  # 更新当前版本号
  echo "$latest_version" > current_version.txt
  cp current_version.txt "$website_dir"
  # 设置文件权限(1panel)
  chown -R 1000:1000 "$website_dir"
fi

# 释放文件锁
flock -u 200
