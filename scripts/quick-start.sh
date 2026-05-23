#!/bin/bash

# BlueBell 快速部署脚本
# 使用方法：./scripts/quick-start.sh

set -e

echo "======================================"
echo "   BlueBell 快速部署脚本"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}错误: 未安装 Docker${NC}"
    echo "请访问 https://docs.docker.com/get-docker/ 安装 Docker"
    exit 1
fi

# 检查 Docker Compose
if ! command -v docker compose &> /dev/null; then
    echo -e "${RED}错误: 未安装 Docker Compose${NC}"
    echo "请访问 https://docs.docker.com/compose/install/ 安装 Docker Compose"
    exit 1
fi

echo -e "${GREEN}✓ Docker 已安装${NC}"
echo -e "${GREEN}✓ Docker Compose 已安装${NC}"
echo ""

# 检查前端文件
if [ ! -d "frontend/dist" ] || [ ! -f "frontend/dist/index.html" ]; then
    echo -e "${YELLOW}前端未构建，开始构建...${NC}"
    echo ""
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        echo -e "${RED}错误: 未安装 Node.js${NC}"
        echo "请访问 https://nodejs.org/ 安装 Node.js 18+"
        exit 1
    fi
    
    echo "安装前端依赖..."
    cd frontend
    npm install
    
    echo "构建前端..."
    npm run build
    
    cd ..
    echo ""
fi

echo -e "${GREEN}✓ 前端文件已准备${NC}"
echo ""

# 检查 .env 文件
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}未找到 .env 文件，使用默认配置${NC}"
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo -e "${GREEN}✓ 已从 .env.example 创建 .env${NC}"
    fi
fi

echo ""
echo "开始部署..."
echo ""

# 停止旧容器
if [ "$(docker compose ps -q 2>/dev/null)" ]; then
    echo "停止现有容器..."
    docker compose down
fi

# 启动服务
echo "启动服务..."
docker compose up -d

# 等待服务启动
echo ""
echo "等待服务启动..."
sleep 10

# 检查服务状态
echo ""
echo "服务状态:"
docker compose ps

echo ""
echo "======================================"
echo -e "${GREEN}部署完成！${NC}"
echo "======================================"
echo ""
echo "访问地址："
echo "  前端：http://localhost"
echo "  API：http://localhost:8080"
echo "  Swagger：http://localhost/swagger/index.html"
echo ""
echo "常用命令："
echo "  查看日志：docker compose logs -f"
echo "  停止服务：docker compose down"
echo "  重启服务：docker compose restart"
echo ""
