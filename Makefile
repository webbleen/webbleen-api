# webbleen-api Makefile
# 极简的开发和管理命令

# 变量定义
APP_NAME = webbleen-api
DOCKER_IMAGE = webbleen-api
DOCKER_TAG = latest
PORT = 8080

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "webbleen-api 可用命令:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""

# 核心开发命令
.PHONY: dev
dev: ## 启动开发环境
	@echo "🐳 启动开发环境..."
	@if ! docker compose version >/dev/null 2>&1; then \
		echo "❌ Docker Compose 未安装，请先安装 Docker Desktop"; \
		exit 1; \
	fi
	@docker compose up -d
	@echo "✅ 服务已启动！"
	@echo "🌐 API: http://localhost:$(PORT)"
	@echo "🤖 聊天: http://localhost:$(PORT)/chat"

.PHONY: stop
stop: ## 停止开发环境
	@echo "🛑 停止开发环境..."
	@if ! docker compose version >/dev/null 2>&1; then \
		echo "❌ Docker Compose 未安装"; \
		exit 1; \
	fi
	@docker compose down

.PHONY: logs
logs: ## 查看日志
	@if ! docker compose version >/dev/null 2>&1; then \
		echo "❌ Docker Compose 未安装"; \
		exit 1; \
	fi
	@docker compose logs -f

.PHONY: build
build: ## 构建镜像
	@echo "🔨 构建 Docker 镜像..."
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "❌ Docker 未安装，请先安装 Docker Desktop"; \
		exit 1; \
	fi
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# 测试命令
.PHONY: test
test: ## 运行测试
	@echo "🧪 运行测试..."
	@./test_env.sh
	@./test_db.sh
	@./test_ai_api.sh

# 代码质量
.PHONY: fmt
fmt: ## 格式化代码
	@if ! command -v go >/dev/null 2>&1; then \
		echo "❌ Go 未安装，请先安装 Go 或使用 Docker 环境"; \
		exit 1; \
	fi
	@go fmt ./...

.PHONY: tidy
tidy: ## 整理依赖
	@if ! command -v go >/dev/null 2>&1; then \
		echo "❌ Go 未安装，请先安装 Go 或使用 Docker 环境"; \
		exit 1; \
	fi
	@go mod tidy

# Railway 部署
.PHONY: deploy
deploy: ## 部署到 Railway
	@echo "🚀 部署到 Railway..."
	@echo "请确保已安装 Railway CLI: npm install -g @railway/cli"
	@echo "然后运行: railway login && railway up"

# 清理
.PHONY: clean
clean: ## 清理
	@if ! docker compose version >/dev/null 2>&1; then \
		echo "❌ Docker Compose 未安装"; \
		exit 1; \
	fi
	@docker compose down -v
	@echo "🧹 清理未使用的容器和网络..."
	@docker container prune -f
	@docker network prune -f

# 快速命令
.PHONY: start
start: dev ## 启动（别名）

.PHONY: restart
restart: stop dev ## 重启

.PHONY: quick
quick: build dev ## 快速启动
	@echo "✅ 快速启动完成！"

# 版本信息
.PHONY: version
version: ## 显示版本
	@echo "App: $(APP_NAME) | Image: $(DOCKER_IMAGE):$(DOCKER_TAG) | Port: $(PORT)"