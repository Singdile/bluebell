# 构建阶段
FROM golang:1.26-alpine AS builder

# 代理
ENV GOPROXY=https://goproxy.cn,direct

# 创建并切换目录
WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

#复制源代码
COPY . .



#编译二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bluebell .


#######################################################
# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要的工具
# 使用 sed 命令将默认的 dl-cdn.alpinelinux.org 替换为阿里云镜像站
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata mysql-client

# 从构建阶段复制文件
COPY --from=builder /build/bluebell .
COPY --from=builder /build/settings/config.docker.yaml ./settings/config.docker.yaml

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8080

# 直接启动应用
CMD ["./bluebell"]
