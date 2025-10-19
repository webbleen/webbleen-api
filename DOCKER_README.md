# webbleen-api Docker 部署指南

本文档说明如何使用 Docker 部署 webbleen-api 服务，包括 AI 聊天功能。

## 🐳 Docker 配置

### 基础镜像
- **构建阶段**: `node:20-bullseye-slim`
- **运行阶段**: `node:20-bullseye-slim`
- **Go 版本**: 1.23.4 (通过命令安装)

### 特性
- 多阶段构建，减小最终镜像大小
- 基于 Debian Bullseye，稳定可靠
- 包含 Node.js 20 和 Go 1.23.4 环境
- 非 root 用户运行，提高安全性
- 支持 PostgreSQL 和 SQLite 数据库

## 🚀 快速开始

### 1. 构建镜像

```bash
# 使用构建脚本
./build.sh webbleen-api:latest

# 或直接使用 Docker 命令
docker build -t webbleen-api:latest .
```

### 2. 运行容器

```bash
# 基本运行
docker run -p 8080:8000 \
  -e CURSOR_API_KEY=your_api_key_here \
  webbleen-api:latest

# 使用环境变量文件
docker run -p 8080:8000 \
  --env-file .env \
  webbleen-api:latest
```

### 3. 使用 Docker Compose

```bash
# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f webbleen-api

# 停止服务
docker-compose down
```

**注意**: 此服务使用外部 PostgreSQL 数据库，请确保数据库服务正在运行。所有配置都通过环境变量管理，无需配置文件。

## 🔧 环境变量配置

所有配置都通过环境变量管理，参考 `env.example` 了解所有可用的环境变量：

```bash
# Cursor AI API Key
CURSOR_API_KEY=your_cursor_api_key_here

# 数据库配置
DATABASE_URL=postgres://postgres:password@host.docker.internal:5432/webbleen_api?sslmode=disable

# 服务器配置
PORT=8000
GIN_MODE=debug

# CORS 配置
CORS_ALLOWED_ORIGINS=*
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=*
CORS_CREDENTIALS=true
```

## 📁 文件结构

```
webbleen-api/
├── Dockerfile              # Docker 镜像定义
├── docker-compose.yml      # Docker Compose 配置
├── .dockerignore           # Docker 忽略文件
├── build.sh               # 构建脚本
├── env.docker             # Docker 环境变量模板
├── env.example            # 开发环境变量模板
└── DOCKER_README.md       # 本文档
```

## 🛠️ 开发模式

### 本地开发
```bash
# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

### Docker 开发
```bash
# 构建开发镜像
docker build -t webbleen-api:dev .

# 运行开发容器
docker run -p 8080:8000 \
  -v $(pwd):/app \
  -e GIN_MODE=debug \
  webbleen-api:dev
```

## 🔍 调试和监控

### 查看容器日志
```bash
# 查看运行中的容器
docker ps

# 查看日志
docker logs -f <container_id>

# 进入容器
docker exec -it <container_id> /bin/bash
```

### 健康检查
```bash
# 检查服务状态
curl http://localhost:8080/healthz

# 检查项目状态
curl http://localhost:8080/api/project/status

# 测试 AI 聊天 API
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "你好", "type": "chat"}'
```

## 🚀 生产部署

### Railway 部署
1. 连接 GitHub 仓库
2. 设置环境变量：
   - `CURSOR_API_KEY`
   - `DATABASE_URL`
3. 部署自动开始

### 其他平台
```bash
# 构建生产镜像
docker build -t webbleen-api:prod .

# 推送到镜像仓库
docker tag webbleen-api:prod your-registry/webbleen-api:latest
docker push your-registry/webbleen-api:latest
```

## 🔒 安全注意事项

1. **API Key 安全**: 不要将 Cursor API Key 提交到版本控制
2. **数据库密码**: 使用强密码并定期更换
3. **网络安全**: 在生产环境中限制 CORS 来源
4. **容器安全**: 使用非 root 用户运行应用

## 🐛 故障排除

### 常见问题

1. **构建失败**
   ```bash
   # 清理 Docker 缓存
   docker system prune -a
   
   # 重新构建
   docker build --no-cache -t webbleen-api:latest .
   ```

2. **数据库连接失败**
   ```bash
   # 检查外部数据库连接
   # 确保 PostgreSQL 服务正在运行在 host.docker.internal:5432
   
   # 测试数据库连接
   ./test_db.sh
   ```

3. **AI 功能不可用**
   - 检查 `CURSOR_API_KEY` 是否正确设置
   - 确保网络可以访问 Cursor API

### 性能优化

1. **镜像大小优化**
   - 使用多阶段构建
   - 清理不必要的文件和缓存

2. **运行时优化**
   - 设置合适的资源限制
   - 使用健康检查

## 📚 相关文档

- [AI Chat 功能文档](./AI_CHAT_MIGRATION.md)
- [API 文档](./SWAGGER.md)
- [部署检查清单](./DEPLOYMENT_CHECKLIST.md)
