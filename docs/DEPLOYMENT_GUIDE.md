# BlueBell Docker 部署流程图

## 部署架构

```
┌─────────────────────────────────────────────────────────────┐
│                        用户浏览器                             │
│                    http://localhost                          │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│                   Nginx (bluebell-nginx)                     │
│                        端口 80                                │
│                                                              │
│  ┌──────────────────┐    ┌──────────────────┐              │
│  │  前端静态文件     │    │  后端 API 代理    │              │
│  │  /               │    │  /v1, /login      │              │
│  │  ↓               │    │  ↓                │              │
│  │  frontend/dist/  │    │  backend:8080     │              │
│  └──────────────────┘    └──────────────────┘              │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│                 Backend (bluebell-backend)                   │
│                      端口 8080                                │
│                                                              │
│  - Go + Gin 框架                                             │
│  - JWT 认证                                                  │
│  - RESTful API                                               │
└─────────────────────────────────────────────────────────────┘
                    ↓           ↓
    ┌───────────────────┐    ┌───────────────────┐
    │  MySQL (3306)      │    │  Redis (6379)     │
    │  bluebell-mysql    │    │  bluebell-redis   │
    │                    │    │                   │
    │  - 用户数据         │    │  - 会话缓存        │
    │  - 帖子数据         │    │  - 投票队列        │
    │  - 社区数据         │    │  - 热点缓存        │
    └───────────────────┘    └───────────────────┘
```

---

## 部署流程

### 方式一：一键部署（推荐）

```
┌─────────────────────────────────────────────────────────────┐
│  步骤 1: 克隆项目                                            │
├─────────────────────────────────────────────────────────────┤
│  git clone https://github.com/Singdile/bluebell.git         │
│  cd bluebell                                                │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  步骤 2: 构建前端                                            │
├─────────────────────────────────────────────────────────────┤
│  cd frontend                                                │
│  npm install                                                │
│  npm run build                                              │
│  cd ..                                                      │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  步骤 3: 启动服务                                            │
├─────────────────────────────────────────────────────────────┤
│  docker compose up -d                                       │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│  步骤 4: 访问应用                                            │
├─────────────────────────────────────────────────────────────┤
│  浏览器打开: http://localhost                                │
└─────────────────────────────────────────────────────────────┘
```

### 方式二：使用快速脚本

```
┌─────────────────────────────────────────────────────────────┐
│  git clone https://github.com/Singdile/bluebell.git         │
│  cd bluebell                                                │
│  ./scripts/quick-start.sh                                   │
└─────────────────────────────────────────────────────────────┘
                          ↓
              自动执行所有步骤
                          ↓
              访问 http://localhost
```

---

## 服务启动顺序

```
1. MySQL 启动
   ↓
2. MySQL 健康检查通过
   ↓
3. Redis 启动
   ↓
4. Redis 健康检查通过
   ↓
5. 数据库迁移执行（migrate）
   ↓
6. 迁移完成
   ↓
7. 后端服务启动（backend）
   ↓
8. Nginx 启动
   ↓
9. 服务就绪 ✅
```

---

## 数据流向

### 用户注册流程

```
用户输入
  ↓
前端表单验证
  ↓
POST /register
  ↓
Nginx 代理
  ↓
Backend 处理
  ├─→ 密码加密
  ├─→ 写入 MySQL
  └─→ 返回成功
  ↓
前端跳转登录页
```

### 帖子浏览流程

```
用户访问帖子列表
  ↓
GET /v1/posts
  ↓
Nginx 代理
  ↓
Backend 处理
  ├─→ 检查 Redis 缓存
  │   ├─→ 缓存命中 → 直接返回
  │   └─→ 缓存未命中 ↓
  ├─→ 查询 MySQL
  ├─→ 写入 Redis 缓存
  └─→ 返回数据
  ↓
前端渲染列表
```

### 投票流程（异步）

```
用户点击投票
  ↓
POST /v1/vote
  ↓
Backend 处理
  ├─→ 验证用户
  ├─→ 写入 Redis 队列
  └─→ 立即返回成功
  ↓
前端更新 UI（乐观更新）
  ↓
后台异步任务（每 5 秒）
  ├─→ 读取 Redis 队列
  ├─→ 批量更新 MySQL
  └─→ 清空队列
```

---

## 网络隔离

```
外部网络（可访问）
├─→ Nginx (0.0.0.0:80)
└─→ Nginx (0.0.0.0:443)

内部网络（仅容器间通信）
├─→ Backend (backend:8080)
├─→ MySQL (mysql:3306)
└─→ Redis (redis:6379)

本地访问（仅宿主机）
├─→ Backend (127.0.0.1:8080)
├─→ MySQL (127.0.0.1:13306)
└─→ Redis (127.0.0.1:6379)
```

---

## 文件结构

```
bluebell/
├── docker-compose.yml       # Docker Compose 配置
├── nginx.conf              # Nginx 配置
├── Dockerfile              # 后端镜像构建
├── .env.example            # 环境变量模板
│
├── frontend/               # 前端源码
│   ├── src/               # Vue 源代码
│   ├── public/            # 静态资源
│   ├── package.json       # 依赖配置
│   └── vite.config.js     # Vite 配置
│
├── migrations/             # 数据库迁移
│   ├── *.up.sql           # 升级脚本
│   └── *.down.sql         # 回滚脚本
│
├── scripts/               # 部署脚本
│   └── quick-start.sh     # 快速部署
│
└── docs/                  # 文档
    └── DOCKER_DEPLOYMENT.md
```

---

## 常见问题排查

### 问题：服务无法启动

```
检查步骤:
1. docker compose ps          # 查看容器状态
2. docker compose logs [服务] # 查看日志
3. docker compose down        # 停止所有服务
4. docker compose up -d       # 重新启动
```

### 问题：前端无法访问后端

```
检查步骤:
1. curl http://localhost/v1/community  # 测试 Nginx 代理
2. curl http://localhost:8080/v1/community  # 测试后端直接访问
3. docker logs bluebell-nginx  # 查看 Nginx 日志
4. 检查 frontend/.env.production 配置
```

### 问题：数据库连接失败

```
检查步骤:
1. docker compose ps mysql  # 查看 MySQL 状态
2. docker logs bluebell-mysql  # 查看日志
3. docker exec -it bluebell-mysql mysql -u root -p  # 测试连接
4. 检查 .env 中的数据库密码
```

---

## 性能优化建议

### 数据库优化

```sql
-- 调整缓冲池大小
SET GLOBAL innodb_buffer_pool_size = 1073741824;  -- 1GB

-- 添加索引
CREATE INDEX idx_posts_community ON posts(community_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

### Redis 优化

```bash
# 设置内存限制
docker exec bluebell-redis redis-cli CONFIG SET maxmemory 256mb

# 设置淘汰策略
docker exec bluebell-redis redis-cli CONFIG SET maxmemory-policy allkeys-lru
```

### Nginx 优化

```nginx
# 启用 Gzip 压缩
gzip on;
gzip_types text/plain text/css application/json application/javascript;

# 启用缓存
location ~* \.(js|css|png|jpg)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

---

## 监控指标

### 容器状态

```bash
# CPU 和内存使用
docker stats

# 磁盘使用
docker system df

# 容器健康状态
docker compose ps
```

### 应用指标

```bash
# 后端响应时间
curl -w "@curl-format.txt" -o /dev/null -s http://localhost/v1/posts

# 数据库连接数
docker exec bluebell-mysql mysql -e "SHOW STATUS LIKE 'Threads_connected'"

# Redis 内存使用
docker exec bluebell-redis redis-cli INFO memory
```

---

## 备份与恢复

### 备份

```bash
# 备份 MySQL
docker exec bluebell-mysql mysqldump -u root -pbluebell321 bluebell > backup_$(date +%Y%m%d).sql

# 备份 Redis
docker exec bluebell-redis redis-cli BGSAVE
docker cp bluebell-redis:/data/dump.rdb redis_backup_$(date +%Y%m%d).rdb
```

### 恢复

```bash
# 恢复 MySQL
docker exec -i bluebell-mysql mysql -u root -pbluebell321 bluebell < backup_20250522.sql

# 恢复 Redis
docker cp redis_backup_20250522.rdb bluebell-redis:/data/dump.rdb
docker restart bluebell-redis
```

---

## 更新部署

### 更新前端

```bash
cd frontend
git pull
npm install
npm run build
cd ..
docker restart bluebell-nginx
```

### 更新后端

```bash
# 拉取最新镜像
docker pull singdile/bluebell:v1.0.0

# 重启服务
docker compose up -d --no-deps backend
```

### 更新全部

```bash
git pull
docker compose pull
docker compose up -d
```
