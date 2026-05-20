# 构建阶段：编译应用
FROM golang:1.26-alpine AS builder

# 为docker 内部的 go 配置代理
ENV GOPROXY=https://goproxy.cn,direct

# 创建并切换到，用于编译的目录
WORKDIR /build

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码到builder 环境中
COPY . .

# 编译应用 bluebell
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bluebell .


##########################################
# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要的工具
RUN apk --no-cache add ca-certificates tzdata

# 复制上一阶段的编译的可执行文件和配置文件
COPY --from=builder /build/bluebell .
COPY --from=builder /build/settings ./settings

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 执行启动命令
CMD ["./bluebell"]
