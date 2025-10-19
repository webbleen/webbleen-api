#!/bin/bash

# API 测试脚本
# 使用方法: ./test_ai_api.sh [base_url]
# 例如: ./test_ai_api.sh http://localhost:8080

BASE_URL=${1:-"http://localhost:8080"}

echo "🔍 测试 API 功能"
echo "基础 URL: $BASE_URL"
echo "================================"

# 测试健康检查
echo "1. 测试健康检查..."
curl -s "$BASE_URL/healthz" | jq '.' || echo "健康检查失败"
echo ""

# 测试就绪检查
echo "2. 测试就绪检查..."
curl -s "$BASE_URL/readyz" | jq '.' || echo "就绪检查失败"
echo ""

# 测试统计 API
echo "3. 测试统计 API..."
curl -s "$BASE_URL/stats/visits" | jq '.' || echo "统计 API 失败"
echo ""

# 测试代理 API
echo "4. 测试代理 API..."
curl -s "$BASE_URL/proxy/ip" | jq '.' || echo "代理 API 失败"
echo ""

echo "================================"
echo "✅ 测试完成！"
echo ""
echo "💡 提示："
echo "- 如果看到正确的响应，说明 API 功能正常工作"
echo "- 如果看到错误，请检查服务器是否正在运行"
echo "- 访问 $BASE_URL/dashboard 查看仪表板"
