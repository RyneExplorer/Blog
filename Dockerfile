# 1. 使用 Go 官方镜像作为构建阶段，避免把编译工具带入最终运行镜像。
FROM golang:1.25-alpine AS builder

# 2. 安装 git，部分 Go 依赖在下载时可能需要它。
RUN apk add --no-cache git

WORKDIR /src

# 3. 先复制依赖清单并下载依赖，利用 Docker 缓存减少重复构建时间。
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download

# 4. 复制业务代码并编译 Linux 静态可执行文件。
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/blog-api ./cmd/server/main.go

# 运行阶段
# 5. 使用轻量 Alpine 镜像作为运行阶段，只保留运行所需文件。
FROM alpine:3.22

# 6. 安装时区和证书，保证 HTTPS 请求、邮件服务和日志时间正常。
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# 7. 复制后端可执行文件。
COPY --from=builder /out/blog-api ./blog-api

# 8. 复制 Docker 专用配置为应用默认读取的 configs/config.yaml。
COPY configs/config.docker.yaml ./configs/config.yaml

# 9. 创建日志和上传目录，后续由 docker-compose 挂载 volume 持久化。
RUN mkdir -p logs uploads

EXPOSE 8083

# 10. 启动 Go 后端服务。
CMD ["./blog-api"]
