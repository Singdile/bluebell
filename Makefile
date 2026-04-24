##
# Project Title
#
# @file
# @version 0.1

# 定义变量
DB_DSN := "mysql://root:123456@tcp(127.0.0.1:3306)/bluebell"
MIGRATE_DIR :=./migrations  #数据库迁移


# .PHONY 表示后面是一个动作
.PHONY: help
help:
		@echo "使用说明："
		@echo "  make migrate-create name=描述  - 创建新的迁移文件"
		@echo "  make migrate-up              - 执行所有向上迁移 (Up)"
		@echo "  make migrate-down            - 执行一步向下回滚 (Down)"
#定义操作
#创建迁移文件 migrate-create name=xxx
migrate-create:
	migrate create -ext sql -dir $(MIGRATE_DIR) -seq $(name)

#执行所有的up迁移
migrate-up:
	migrate -path $(MIGRATE_DIR) -database $(DB_DSN) up

#执行一步 down 回滚
migrate-down:
	migrate -path $(MIGRATE_DIR) -database $(DB_DSN) down 1

# 声明这些都是“动作”，不是文件
.PHONY: help migrate-create migrate-up migrate-down

# end
