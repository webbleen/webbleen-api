FROM golang:1.24-alpine AS builder

WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -o webbleen-api main.go

# 最终镜像
FROM alpine:latest

# 安装 ca-certificates 和 tzdata
RUN apk --no-cache add ca-certificates tzdata

# 创建非root用户
RUN adduser -D -s /bin/sh appuser

WORKDIR /app

# 复制构建的二进制文件
COPY --from=builder /app/webbleen-api .

# 复制配置文件和模板文件
COPY --from=builder /app/conf ./conf
COPY --from=builder /app/web ./web

# 更改文件所有者
RUN chown -R appuser:appuser /app

# 切换到非root用户
USER appuser

# 暴露端口（Railway 会动态分配端口，这里使用默认值）
EXPOSE 8000

# 运行应用
CMD ["./webbleen-api"]