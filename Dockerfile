# 使用官方 Go 镜像作为构建环境
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY geoserver/go.mod geoserver/go.sum ./

# 下载所有依赖
RUN go mod download

# 复制整个 geoserver 目录
COPY geoserver/ ./

# 确认文件已正确复制
RUN ls -la

# 编译应用程序。确保 CGO 被禁用，生成的二进制文件为 Linux 平台
# 调整路径以指向包含 main 函数的 go 文件所在目录
RUN CGO_ENABLED=0 GOOS=linux go build -o geoserver api/geoserver.go

# 使用 alpine 作为最小运行环境，因为它包含了必要的基本命令行工具
FROM alpine

# 创建所需的目录结构
RUN mkdir -p /geoserver/api/image/
RUN mkdir -p /geoserver/api/etc/
# 设置工作目录
WORKDIR /geoserver/api/

# 复制构建好的应用程序和配置文件
COPY --from=builder /app/geoserver /geoserver/api
COPY --from=builder /app/api/etc/geoserver-api.yaml /geoserver/api/etc/
COPY --from=builder /app/api/image/geoserverImage.tar /geoserver/api/image/

# 确认文件已正确复制
RUN ls -la

# 暴露端口
EXPOSE 8888

# 运行编译好的二进制文件
CMD ["/geoserver/api/geoserver"]
