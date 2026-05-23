#!/bin/sh
set -e

echo "=> 等待 MySQL 启动..."
until mysql -h mysql -P 3306 -u root -pbluebell321 --skip-ssl -e "SELECT 1" &> /dev/null; do
  echo "   MySQL 未就绪，等待中..."
  sleep 2
done

echo "=> MySQL 已就绪，开始执行数据库迁移..."
migrate -path /app/migrations -database "mysql://root:bluebell321@tcp(mysql:3306)/bluebell" up

echo "=> 迁移完成，启动应用..."
exec ./bluebell