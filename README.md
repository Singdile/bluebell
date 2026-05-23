# BlueBell 社区论坛

一个基于 Go + Vue 的现代化社区论坛系统。

## 技术栈

**后端**
- Go 1.26 + Gin
- MySQL 8.0 + Redis 7
- JWT 认证 + Swagger 文档

**前端**
- Vue 3 + Vite


## 功能特性

- 用户注册/登录（JWT 认证）
- 社区创建与管理
- 帖子发布与浏览
- 投票系统（点赞/踩）
- 按热度/时间排序
- API 限流保护

## 快速部署

```bash
# 1. 克隆项目
git clone https://github.com/Singdile/bluebell.git
cd bluebell

# 2. 启动服务
docker compose up -d

# 3. 浏览器访问
http://localhost
```

## API 文档

启动后访问：http://localhost/swagger/index.html

## 默认配置

- MySQL 密码：`bluebell321`
- Redis：无密码
- 后端端口：8080（内部）
- 前端端口：80（Nginx）

## 修改密码

如需修改 MySQL 密码，编辑 `.env` 文件：

```bash
cp .env.example .env
# 编辑 .env 中的 MYSQL_ROOT_PASSWORD
docker compose up -d
```

注意：修改密码后需要重新构建后端镜像。

## 项目结构

```
bluebell/
├── controllers/    # 控制器
├── dao/           # 数据访问层
├── logic/         # 业务逻辑
├── models/        # 数据模型
├── router/        # 路由
├── frontend/dist/ # 前端构建产物
├── docs/          # Swagger 文档
└── docker-compose.yml
```

## License

MIT
