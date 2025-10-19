# Makefile 使用说明

webbleen-api 项目提供了极简的 Makefile，只包含最核心的开发和管理命令。

## 🚀 快速开始

### 查看所有命令
```bash
make help
```

### 快速启动
```bash
make quick
```

## 📋 命令列表

### 核心开发命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `make dev` | 启动开发环境 | `make dev` |
| `make start` | 启动（别名） | `make start` |
| `make stop` | 停止开发环境 | `make stop` |
| `make restart` | 重启服务 | `make restart` |
| `make logs` | 查看日志 | `make logs` |

### 构建和测试

| 命令 | 说明 | 示例 |
|------|------|------|
| `make build` | 构建镜像 | `make build` |
| `make test` | 运行测试 | `make test` |

### 代码质量

| 命令 | 说明 | 示例 |
|------|------|------|
| `make fmt` | 格式化代码 | `make fmt` |
| `make tidy` | 整理依赖 | `make tidy` |

### 部署和清理

| 命令 | 说明 | 示例 |
|------|------|------|
| `make deploy` | 部署到 Railway | `make deploy` |
| `make clean` | 清理 | `make clean` |

### 其他

| 命令 | 说明 | 示例 |
|------|------|------|
| `make quick` | 快速启动 | `make quick` |
| `make version` | 显示版本 | `make version` |

## 🔧 常用工作流

### 1. 新开发者快速开始
```bash
make quick
```

### 2. 日常开发
```bash
# 启动开发环境
make dev

# 查看日志
make logs

# 运行测试
make test
```

### 3. 代码提交前
```bash
make fmt
make tidy
make test
```

### 4. Railway 部署
```bash
# 安装 Railway CLI
npm install -g @railway/cli

# 登录并部署
make deploy
```

### 5. 问题排查
```bash
# 查看日志
make logs

# 重启服务
make restart

# 清理并重新开始
make clean
make quick
```

## 🚀 Railway 部署

### 首次部署
1. 安装 Railway CLI
   ```bash
   npm install -g @railway/cli
   ```

2. 部署项目
   ```bash
   make deploy
   ```

3. 按照提示完成登录和部署

### 环境变量设置
在 Railway 控制台中设置：
- `DATABASE_URL` - PostgreSQL 数据库连接字符串
- `CURSOR_API_KEY` - Cursor AI API 密钥

## 💡 提示

- 使用 `make help` 查看所有命令
- 本地开发使用 Docker，生产部署使用 Railway
- 大部分命令都很简单，易于记忆和使用