# 使用官方的 Golang 镜像作为构建阶段
FROM golang:1.16-alpine AS build

# 设置工作目录
WORKDIR /

# 复制 Go 模块文件并下载依赖项
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码到工作目录
COPY . .

# 编译 Go 应用程序
RUN go build -o holidays

# 运行时阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /

# 从构建阶段复制编译好的应用程序
COPY --from=build /holidays .

# 声明需要暴露的端口
EXPOSE 8081

# 启动应用程序
CMD ["./holidays"]
