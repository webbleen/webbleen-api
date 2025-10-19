#!/bin/bash

# webbleen-api Docker 构建脚本
# 使用方法: ./build.sh [tag]
# 例如: ./build.sh webbleen-api:latest

TAG=${1:-"webbleen-api:latest"}

echo "🐳 开始构建 webbleen-api Docker 镜像..."
echo "标签: $TAG"
echo "================================"

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker"
    exit 1
fi

# 构建镜像
echo "📦 构建 Docker 镜像..."
docker build -t "$TAG" .

if [ $? -eq 0 ]; then
    echo "✅ 镜像构建成功！"
    echo ""
    echo "🚀 运行容器："
    echo "docker run -p 8080:8000 -e CURSOR_API_KEY=your_key_here $TAG"
    echo ""
    echo "🔍 查看镜像信息："
    echo "docker images $TAG"
else
    echo "❌ 镜像构建失败"
    exit 1
fi
