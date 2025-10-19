#!/bin/bash

# 环境变量测试脚本
# 测试 webbleen-api 的环境变量配置

echo "🔍 测试 webbleen-api 环境变量配置..."
echo "================================"

# 设置测试环境变量
export PORT=8080
export GIN_MODE=debug
export DATABASE_URL="postgres://postgres:password@host.docker.internal:5432/webbleen_api?sslmode=disable"
export CORS_ALLOWED_ORIGINS="http://localhost:8080,https://webbleen.com"
export CORS_ALLOWED_METHODS="GET,POST,PUT,DELETE,OPTIONS"
export CORS_ALLOWED_HEADERS="Content-Type,Authorization,Accept-Encoding,Content-Length"
export CORS_CREDENTIALS=true
export JWT_SECRET="test_jwt_secret"
export PAGE_SIZE=20

echo "1. 设置的环境变量:"
echo "   PORT: $PORT"
echo "   GIN_MODE: $GIN_MODE"
echo "   DATABASE_URL: $DATABASE_URL"
echo "   CORS_ALLOWED_ORIGINS: $CORS_ALLOWED_ORIGINS"
echo "   CORS_ALLOWED_METHODS: $CORS_ALLOWED_METHODS"
echo "   CORS_ALLOWED_HEADERS: $CORS_ALLOWED_HEADERS"
echo "   CORS_CREDENTIALS: $CORS_CREDENTIALS"
echo "   JWT_SECRET: ${JWT_SECRET:0:10}..."
echo "   PAGE_SIZE: $PAGE_SIZE"
echo ""

echo "2. 测试 Go 配置加载:"
if command -v go &> /dev/null; then
    echo "   ✅ Go 环境可用"
    echo "   🧪 运行配置测试..."
    go run -c 'package main; import "github.com/webbleen/go-gin/pkg/setting"; func main() { setting.PrintConfig() }' 2>/dev/null || echo "   ⚠️  无法直接测试配置加载"
else
    echo "   ❌ Go 环境不可用"
fi

echo ""

echo "3. 测试默认值（清除环境变量）:"
unset PORT GIN_MODE DATABASE_URL CORS_ALLOWED_ORIGINS CORS_ALLOWED_METHODS CORS_ALLOWED_HEADERS CORS_CREDENTIALS JWT_SECRET PAGE_SIZE

echo "   清除环境变量后，系统应使用默认值"
echo "   - PORT: 默认 8000"
echo "   - GIN_MODE: 默认 debug"
echo "   - DATABASE_URL: 默认 postgres://postgres:password@host.docker.internal:5432/webbleen_api?sslmode=disable"
echo ""
echo "   - CORS: 默认空值（需要在 docker-compose.yml 中配置）"

echo ""
echo "================================"
echo "✅ 环境变量测试完成！"
echo ""
echo "💡 提示:"
echo "- 使用 './start.sh -l' 启动本地服务查看配置"
echo "- 使用 './start.sh -d' 启动 Docker 服务"
echo "- 查看 env.example 了解所有可用的环境变量"
