# 1. 变量定义 (放在最顶层，方便以后修改)
APP_NAME    := bluebell
DB_DSN      := mysql://root:123456@tcp(127.0.0.1:3306)/bluebell
MIGRATE_DIR := ./migrations

# 2. 【核心】默认目标 (必须放在第一个目标位置)
# 这样你按 C-c c + 回车，直接就开始整理依赖并编译了
all: tidy build

# 3. 基础开发指令 (Build, Run, Test)
build:
	@echo "=> 🚀 正在构建二进制文件 [$(APP_NAME)]..."
	go build -o $(APP_NAME) main.go

run:
	@echo "=> ⚡ 正在启动应用..."
	go run main.go

test:
	@echo "=> 🧪 正在执行全量单元测试..."
	go test -v -race ./...

tidy:
	@echo "=> 📦 正在整理 go.mod 依赖..."
	go mod tidy
	go fmt ./...

# 4. 数据库迁移指令 (特殊操作)
migrate-create:
	@if [ -z "$(name)" ]; then echo "错误: 请提供迁移名称，例如 make migrate-create name=add_user_table"; exit 1; fi
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)

migrate-up:
	migrate -path $(MIGRATE_DIR) -database "$(DB_DSN)" up

migrate-down:
	migrate -path $(MIGRATE_DIR) -database "$(DB_DSN)" down 1

# 5. 清理与帮助 (工具类)
redis:
	docker run -d --name redis-blue -p 6379:6379 redis:latest


clean:
	@echo "=> 🧹 正在清理构建缓存..."
	rm -f $(APP_NAME)
	go clean

help:
	@echo "使用说明："
	@echo "  make (all)             - 整理依赖并编译项目 (默认)"
	@echo "  make run               - 直接运行项目"
	@echo "  make migrate-create name=xxx - 创建迁移文件"
	@echo "  make migrate-up        - 执行所有 Up 迁移"
	@echo "  make migrate-down      - 执行一步 Down 回滚"

# 6. 【统一声明】.PHONY (放在最后或目标上方，防止冲突)
.PHONY: all build run test tidy clean help migrate-create migrate-up migrate-down
