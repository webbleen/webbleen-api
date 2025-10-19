#!/bin/bash

# webbleen-api 快速启动脚本
# 使用方法: ./start.sh [选项]
# 选项:
#   -d, --docker    使用 Docker Compose 启动
#   -l, --local     本地启动（需要 Go 环境）
#   -t, --test      启动后运行测试
#   -h, --help      显示帮助信息

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 默认选项
USE_DOCKER=false
RUN_TEST=false
LOCAL_START=false

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--docker)
            USE_DOCKER=true
            shift
            ;;
        -l|--local)
            LOCAL_START=true
            shift
            ;;
        -t|--test)
            RUN_TEST=true
            shift
            ;;
        -h|--help)
            echo "webbleen-api 快速启动脚本"
            echo ""
            echo "使用方法: $0 [选项]"
            echo ""
            echo "选项:"
            echo "  -d, --docker    使用 Docker Compose 启动"
            echo "  -l, --local     本地启动（需要 Go 环境）"
            echo "  -t, --test      启动后运行测试"
            echo "  -h, --help      显示帮助信息"
            echo ""
            echo "示例:"
            echo "  $0 -d           # 使用 Docker 启动"
            echo "  $0 -l -t        # 本地启动并运行测试"
            echo "  $0 --docker --test  # Docker 启动并运行测试"
            exit 0
            ;;
        *)
            echo "未知选项: $1"
            echo "使用 -h 或 --help 查看帮助信息"
            exit 1
            ;;
    esac
done

# 如果没有指定启动方式，默认使用 Docker
if [ "$USE_DOCKER" = false ] && [ "$LOCAL_START" = false ]; then
    USE_DOCKER=true
fi

echo -e "${BLUE}🚀 启动 webbleen-api 服务...${NC}"
echo "================================"

# Docker 启动
if [ "$USE_DOCKER" = true ]; then
    echo -e "${YELLOW}🐳 使用 Docker Compose 启动...${NC}"
    
    # 检查 Docker 是否运行
    if ! docker info > /dev/null 2>&1; then
        echo -e "${RED}❌ Docker 未运行，请先启动 Docker${NC}"
        exit 1
    fi
    
    # 启动服务
    docker-compose up -d
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Docker 服务启动成功！${NC}"
        echo ""
        echo "📊 服务状态:"
        docker-compose ps
        echo ""
        echo "🌐 访问地址:"
        echo "  - API 服务: http://localhost:8080"
        echo "  - AI 聊天界面: http://localhost:8080/chat"
        echo "  - AI 聊天 API: http://localhost:8080/api/ai/chat"
        echo ""
        echo "📝 查看日志:"
        echo "  docker-compose logs -f webbleen-api"
        echo ""
        echo "🛑 停止服务:"
        echo "  docker-compose down"
    else
        echo -e "${RED}❌ Docker 服务启动失败${NC}"
        exit 1
    fi
fi

# 本地启动
if [ "$LOCAL_START" = true ]; then
    echo -e "${YELLOW}💻 本地启动...${NC}"
    
    # 检查 Go 是否安装
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Go 未安装，请先安装 Go 环境${NC}"
        exit 1
    fi
    
    # 安装依赖
    echo "📦 安装依赖..."
    go mod tidy
    
    # 设置开发环境变量（如果未设置）
    if [ -z "$DATABASE_URL" ]; then
        export DATABASE_URL="postgres://postgres:password@host.docker.internal:5432/webbleen_api?sslmode=disable"
        echo "🔧 设置默认数据库: PostgreSQL"
        echo "⚠️  请确保外部 PostgreSQL 服务正在运行在 host.docker.internal:5432"
    fi
    
    # 启动服务
    echo "🚀 启动服务..."
    go run main.go &
    SERVER_PID=$!
    
    # 等待服务启动
    echo "⏳ 等待服务启动..."
    sleep 3
    
    # 检查服务是否启动成功
    if curl -s http://localhost:8080/healthz > /dev/null; then
        echo -e "${GREEN}✅ 本地服务启动成功！${NC}"
        echo ""
        echo "🌐 访问地址:"
        echo "  - API 服务: http://localhost:8080"
        echo "  - AI 聊天界面: http://localhost:8080/chat"
        echo "  - AI 聊天 API: http://localhost:8080/api/ai/chat"
        echo ""
        echo "🛑 停止服务: kill $SERVER_PID"
    else
        echo -e "${RED}❌ 本地服务启动失败${NC}"
        kill $SERVER_PID 2>/dev/null
        exit 1
    fi
fi

# 运行测试
if [ "$RUN_TEST" = true ]; then
    echo ""
    echo -e "${YELLOW}🧪 运行测试...${NC}"
    
    # 等待服务完全启动
    sleep 5
    
    # 运行 API 测试
    if [ -f "./test_ai_api.sh" ]; then
        ./test_ai_api.sh
    else
        echo -e "${YELLOW}⚠️  测试脚本不存在，跳过测试${NC}"
    fi
fi

echo ""
echo -e "${GREEN}🎉 启动完成！${NC}"
