# BlueBell Docker 部署指南

## 目录

- [前置要求](#前置要求)
- [快速开始](#快速开始)
- [详细部署步骤](#详细部署步骤)
- [配置说明](#配置说明)
- [常见问题](#常见问题)

---

## 前置要求

- Docker 20.10+
- Docker Compose v2+
- Git
- 2GB+ 可用内存
- 10GB+ 可用磁盘空间

### 安装 Docker（Ubuntu/Debian）

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | bash

# 启动 Docker
sudo systemctl start docker
sudo systemctl enable docker

# 添加当前用户到 docker 组（可选，免 sudo）
sudo usermod -aG docker $USER
newgrp docker
```

---

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/bluebell.git
cd bluebell
```

### 2. 一键部署

```bash
# 使用默认配置快速启动
docker compose up -d

# 查看服务状态
docker compose ps
```

### 3. 访问应用

打开浏览器访问：`http://localhost`

**默认端口**：
- 前端访问：`http://localhost` (80 端口)
- 后端 API：`http://localhost:8080` (仅内部访问，已通过 Nginx 代理)

---

## 详细部署步骤

### 步骤 1：克隆前端项目

```bash
# 克隆前端项目
cd ..
git clone https://github.com/yourusername/bluebell-web.git

# 进入前端目录
cd bluebell-web

# 安装依赖
npm install

# 构建生产版本
npm run build

# 返回后端项目目录
cd ../bluebell
```

### 步骤 2：准备前端文件

```bash
# 创建前端目录
mkdir -p frontend/dist

# 复制前端构建产物
cp -r ../bluebell-web/dist/* frontend/dist/

# 验证前端文件
ls -la frontend/dist/
```

### 步骤 3：配置环境变量（可选）

如果你需要自定义配置，创建 `.env` 文件：

```bash
cat > .env << 'EOF'
# MySQL 配置
MYSQL_ROOT_PASSWORD=your_secure_password
MYSQL_DATABASE=bluebell
MYSQL_USER=bluebell
MYSQL_PASSWORD=your_secure_password

# Redis 配置（可选，使用默认配置即可）

# 后端配置
BLUEBELL_MODE=release
EOF
```

### 步骤 4：启动服务

```bash
# 启动所有服务（后台运行）
docker compose up -d

# 查看启动日志
docker compose logs -f

# 按 Ctrl+C 退出日志查看
```

### 步骤 5：验证部署

```bash
# 检查所有容器状态
docker compose ps

# 应该看到以下容器运行中：
# - bluebell-mysql    (healthy)
# - bluebell-redis    (healthy)
# - bluebell-migrate  (Exited 0) - 正常，这是一次性任务
# - bluebell-backend  (Up)
# - bluebell-nginx    (Up)
```

### 步骤 6：测试 API

```bash
# 测试注册接口
curl -X POST http://localhost/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123","re_password":"test123"}'

# 测试登录接口
curl -X POST http://localhost/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}'

# 测试社区列表（需要认证）
curl http://localhost/v1/community \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 配置说明

### 目录结构

```
bluebell/
├── docker-compose.yml       # Docker Compose 配置
├── nginx.conf              # Nginx 配置文件
├── Dockerfile              # 后端镜像构建文件
├── .env.example            # 环境变量示例文件
├── migrations/             # 数据库迁移文件
│   ├── 20250101000001_init_schema.up.sql
│   ├── 20250101000001_init_schema.down.sql
│   └── ...
├── frontend/               # 前端文件目录
│   └── dist/
│       ├── index.html
│       ├── assets/
│       └── ...
└── docs/
    └── DOCKER_DEPLOYMENT.md
```

### 服务说明

| 服务 | 容器名 | 端口 | 说明 |
|------|--------|------|------|
| MySQL | bluebell-mysql | 127.0.0.1:13306 | 数据库，仅本地访问 |
| Redis | bluebell-redis | 127.0.0.1:6379 | 缓存，仅本地访问 |
| Backend | bluebell-backend | 127.0.0.1:8080 | 后端 API，通过 Nginx 代理 |
| Nginx | bluebell-nginx | 0.0.0.0:80 | 反向代理，对外暴露 |
| Migrate | bluebell-migrate | - | 数据库迁移（一次性任务） |

### 环境变量

#### MySQL 配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `MYSQL_ROOT_PASSWORD` | bluebell321 | MySQL root 密码 |
| `MYSQL_DATABASE` | bluebell | 数据库名称 |
| `MYSQL_USER` | bluebell | 数据库用户名 |
| `MYSQL_PASSWORD` | bluebell321 | 数据库用户密码 |

#### 后端配置

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `BLUEBELL_MODE` | release | 运行模式：debug/release |

---

## 常用命令

### 服务管理

```bash
# 启动所有服务
docker compose up -d

# 停止所有服务
docker compose down

# 重启所有服务
docker compose restart

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f [service_name]

# 进入容器
docker exec -it bluebell-backend sh
docker exec -it bluebell-mysql mysql -u root -p
```

### 数据管理

```bash
# 备份 MySQL 数据
docker exec bluebell-mysql mysqldump -u root -pbluebell321 bluebell > backup.sql

# 恢复 MySQL 数据
docker exec -i bluebell-mysql mysql -u root -pbluebell321 bluebell < backup.sql

# 清空所有数据（危险操作！）
docker compose down -v
```

### 更新部署

```bash
# 拉取最新代码
git pull

# 重新构建并启动
docker compose up -d --build

# 或者使用最新镜像
docker compose pull
docker compose up -d
```

---

## 常见问题

### Q1: 端口 80 被占用

**错误信息**：
```
Error: bind: address already in use
```

**解决方案**：

```bash
# 查看占用 80 端口的进程
sudo lsof -i :80

# 停止系统 Nginx（如果安装了）
sudo systemctl stop nginx
sudo systemctl disable nginx

# 或者修改 docker-compose.yml 中的端口映射
# 将 "80:80" 改为 "8080:80"
```

### Q2: 数据库连接失败

**错误信息**：
```
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

**解决方案**：

```bash
# 检查 MySQL 容器状态
docker compose ps mysql

# 查看 MySQL 日志
docker compose logs mysql

# 等待 MySQL 完全启动（健康检查通过）
docker compose up -d --wait
```

### Q3: 前端页面加载但 API 请求失败

**原因**：前端配置的 API 地址错误

**解决方案**：

确保前端 `.env.production` 配置正确：

```bash
# 在前端项目中
cat > .env.production << 'EOF'
VITE_APP_TITLE=BlueBell论坛
VITE_API_BASE_URL=
VITE_USE_MOCK=false
EOF

# 重新构建
npm run build

# 复制到后端项目
cp -r dist/* ../bluebell/frontend/dist/
```

### Q4: Nginx 返回 405 Not Allowed

**原因**：Nginx location 配置错误

**解决方案**：

确保 `nginx.conf` 中的 location 匹配正确：

```nginx
# 正确 ✅
location ~ ^/(v1|login|register|swagger) {
    proxy_pass http://backend:8080;
}

# 错误 ❌（末尾有多余的斜杠）
location ~ ^/(v1|login|register|swagger)/ {
    proxy_pass http://backend:8080;
}
```

### Q5: 数据库迁移失败

**解决方案**：

```bash
# 查看迁移日志
docker compose logs migrate

# 手动执行迁移
docker compose run --rm migrate up

# 回滚迁移
docker compose run --rm migrate down
```

### Q6: 如何修改默认密码

**解决方案**：

```bash
# 1. 修改 .env 文件中的密码
vim .env

# 2. 停止并删除所有容器和数据卷（会清空数据！）
docker compose down -v

# 3. 重新启动
docker compose up -d
```

---

## 生产环境建议

### 1. 安全加固

- 修改默认数据库密码
- 使用 HTTPS（配置 SSL 证书）
- 限制容器资源使用
- 定期更新镜像

### 2. 性能优化

- 调整 MySQL 缓冲池大小
- 配置 Redis 内存限制
- 启用 Nginx Gzip 压缩
- 使用 CDN 加速静态资源

### 3. 备份策略

```bash
# 定期备份脚本
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec bluebell-mysql mysqldump -u root -pbluebell321 bluebell > backup_${DATE}.sql
# 保留最近 7 天的备份
find . -name "backup_*.sql" -mtime +7 -delete
```

### 4. 监控告警

- 使用 Prometheus + Grafana 监控
- 配置日志收集（ELK Stack）
- 设置健康检查告警

---

## 技术栈

- **后端**：Go 1.26 + Gin + GORM + Redis
- **前端**：Vue 3 + Vite + Element Plus
- **数据库**：MySQL 8.0
- **缓存**：Redis 7
- **反向代理**：Nginx 1.31
- **容器化**：Docker + Docker Compose

---

## 许可证

MIT License

---

## 联系方式

- 作者：Your Name
- Email：your.email@example.com
- GitHub：https://github.com/yourusername/bluebell

---

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request
