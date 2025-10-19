FROM node:20-bullseye-slim AS builder

# 设置工作目录
WORKDIR /app

# 更新包列表并安装必要的工具
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    ca-certificates \
    tzdata \
    git \
    && rm -rf /var/lib/apt/lists/*

# 安装 Go 1.24.0 (支持最新依赖)
RUN wget -O go1.24.0.linux-amd64.tar.gz https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz \
    && rm go1.24.0.linux-amd64.tar.gz

# 设置 Go 环境变量
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV GOCACHE="/go/cache"

# 创建 Go 工作目录
RUN mkdir -p /go/src /go/bin /go/cache

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN /usr/local/go/bin/go mod download

# 复制源代码
COPY . .

# 整理依赖并构建应用
RUN /usr/local/go/bin/go mod tidy && /usr/local/go/bin/go build -o webbleen-api main.go

# 最终镜像
FROM node:20-bullseye-slim

# 更新包列表并安装运行时依赖
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

# 创建非root用户
RUN groupadd -r appuser && useradd -r -g appuser appuser

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