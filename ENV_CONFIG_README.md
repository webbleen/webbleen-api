# 环境变量配置说明

webbleen-api 现在完全使用环境变量进行配置，无需任何配置文件。

## 🔧 配置方式

### 1. 环境变量设置

所有配置都通过环境变量管理，参考 `env.example` 了解所有可用的环境变量：

```bash
# 服务器配置
PORT=8000
GIN_MODE=debug
READ_TIMEOUT=60
WRITE_TIMEOUT=60

# 数据库配置
DATABASE_URL=postgres://postgres:password@host.docker.internal:5432/webbleen_api?sslmode=disable
# 或使用 SQLite（可选）
# DATABASE_URL=sqlite://webbleen.db


# CORS 配置（可选，默认空值）
# CORS_ALLOWED_ORIGINS=http://127.0.0.1:8080,http://localhost:1313,https://webbleen.com
# CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
# CORS_ALLOWED_HEADERS=Content-Type,Authorization,Accept-Encoding,Content-Length
# CORS_CREDENTIALS=true

# 应用配置
JWT_SECRET=your_jwt_secret_here
PAGE_SIZE=10
```

### 2. 默认值

如果环境变量未设置，系统将使用开发环境的默认值：

| 环境变量 | 默认值 | 说明 |
|---------|--------|------|
| `PORT` | `8000` | HTTP 服务端口 |
| `GIN_MODE` | `debug` | Gin 运行模式 |
| `READ_TIMEOUT` | `60` | 读取超时（秒） |
| `WRITE_TIMEOUT` | `60` | 写入超时（秒） |
| `DATABASE_URL` | `postgres://postgres:password@host.docker.internal:5432/webbleen_api?sslmode=disable` | 数据库连接字符串 |
| `CORS_ALLOWED_ORIGINS` | 空值 | CORS 允许的来源 |
| `CORS_ALLOWED_METHODS` | 空值 | CORS 允许的方法 |
| `CORS_ALLOWED_HEADERS` | 空值 | CORS 允许的头部 |
| `CORS_CREDENTIALS` | `false` | 是否允许携带凭据 |
| `JWT_SECRET` | 默认密钥 | JWT 签名密钥 |
| `PAGE_SIZE` | `10` | 分页大小 |

## 🚀 使用方法

### 1. 本地开发

```bash
# 设置环境变量
export DATABASE_URL="postgres://postgres:password@host.docker.internal:5432/webbleen_api?sslmode=disable"

# 启动服务
./start.sh -l
```

### 2. Docker 开发

```bash
# 使用 Docker Compose（已预配置环境变量）
./start.sh -d

# 或手动启动
docker-compose up -d
```

### 3. 生产环境

```bash
# 设置生产环境变量
export PORT=8080
export GIN_MODE=release
export DATABASE_URL="postgres://user:pass@host:5432/db"
export JWT_SECRET="production_jwt_secret"

# 启动服务
go run main.go
```

## 🔍 配置验证

### 1. 测试环境变量

```bash
# 运行环境变量测试
./test_env.sh
```

### 2. 查看当前配置

启动服务时会自动打印当前配置信息：

```
=== 当前配置 ===
运行模式: debug
HTTP 端口: 8000
读取超时: 1m0s
写入超时: 1m0s
分页大小: 10
数据库 URL: sqlite://webbleen.db
Cursor API Key: key_***e9d
CORS 允许来源: [http://127.0.0.1:8080 http://localhost:1313 https://webbleen.com]
CORS 允许方法: [GET POST PUT DELETE OPTIONS]
CORS 允许头部: [Content-Type Authorization Accept-Encoding Content-Length]
CORS 允许凭据: true
================
```

## 📁 文件结构

```
webbleen-api/
├── pkg/setting/setting.go          # 配置管理（环境变量）
├── env.example                     # 环境变量示例
├── docker-compose.yml              # Docker 环境变量配置
├── test_env.sh                     # 环境变量测试脚本
├── start.sh                        # 启动脚本
└── ENV_CONFIG_README.md            # 本文档
```

## ⚠️ 注意事项

1. **安全性**: 生产环境请设置强密码和安全的 API Key
2. **环境隔离**: 不同环境使用不同的环境变量值
3. **敏感信息**: 不要将包含敏感信息的 `.env` 文件提交到版本控制
4. **默认值**: 开发环境默认值仅用于开发，生产环境必须设置所有必要的环境变量

## 🐛 故障排除

### 常见问题

1. **配置未生效**
   ```bash
   # 检查环境变量
   env | grep -E "(PORT|DATABASE_URL)"
   
   # 重新启动服务
   ./start.sh -l
   ```

2. **数据库连接失败**
   ```bash
   # 检查数据库 URL 格式
   echo $DATABASE_URL
   
   # 测试数据库连接
   # PostgreSQL: psql $DATABASE_URL
   # SQLite: sqlite3 webbleen.db
   ```

3. **CORS 问题**
   ```bash
   # 检查 CORS 配置
   echo $CORS_ALLOWED_ORIGINS
   
   # 确保包含正确的域名
   ```

## 📚 相关文档

- [AI Chat 功能文档](./AI_CHAT_MIGRATION.md)
- [Docker 部署文档](./DOCKER_README.md)
- [API 测试文档](./test_ai_api.sh)

现在 webbleen-api 完全基于环境变量配置，更加灵活和安全！
