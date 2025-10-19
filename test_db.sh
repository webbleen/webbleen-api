#!/bin/bash

# 数据库连接测试脚本
# 测试 PostgreSQL 数据库连接

echo "🔍 测试 PostgreSQL 数据库连接..."
echo "================================"

# 数据库连接信息
# 检测是否在 Docker 容器内运行
if [ -f /.dockerenv ]; then
    # 在 Docker 容器内，使用 host.docker.internal
    DB_HOST="host.docker.internal"
else
    # 在宿主机器上，使用 localhost
    DB_HOST="localhost"
fi

DB_PORT="5432"
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="webbleen_api"

echo "1. 数据库连接信息:"
echo "   主机: $DB_HOST"
echo "   端口: $DB_PORT"
echo "   用户: $DB_USER"
echo "   数据库: $DB_NAME"
echo ""

# 检查 PostgreSQL 客户端是否可用
if command -v psql &> /dev/null; then
    echo "2. 测试数据库连接:"
    
    # 设置环境变量
    export PGPASSWORD="$DB_PASSWORD"
    
    # 测试连接
    if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT version();" &> /dev/null; then
        echo "   ✅ 数据库连接成功"
        
        # 显示数据库版本
        echo "   📊 数据库版本:"
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT version();" | head -3
        
        # 检查表是否存在
        echo "   📋 检查表结构:"
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "\dt" 2>/dev/null || echo "   ⚠️  表结构检查失败（可能数据库为空）"
        
    else
        echo "   ❌ 数据库连接失败"
        echo "   💡 请检查："
        echo "      - PostgreSQL 服务是否正在运行"
        echo "      - 数据库地址是否正确: $DB_HOST:$DB_PORT"
        echo "      - 用户名和密码是否正确"
        echo "      - 数据库 $DB_NAME 是否存在"
    fi
    
    # 清理环境变量
    unset PGPASSWORD
    
else
    echo "2. 测试数据库连接:"
    echo "   ❌ psql 客户端不可用"
    echo "   💡 请安装 PostgreSQL 客户端或使用 Docker"
fi

echo ""

# 测试网络连接
echo "3. 测试网络连接:"
if ping -c 1 "$DB_HOST" &> /dev/null; then
    echo "   ✅ 主机 $DB_HOST 可达"
else
    echo "   ❌ 主机 $DB_HOST 不可达"
    echo "   💡 请检查网络连接或主机地址"
fi

# 测试端口连接
echo "4. 测试端口连接:"
if nc -z "$DB_HOST" "$DB_PORT" 2>/dev/null; then
    echo "   ✅ 端口 $DB_PORT 可访问"
else
    echo "   ❌ 端口 $DB_PORT 不可访问"
    echo "   💡 请检查 PostgreSQL 服务是否在端口 $DB_PORT 上运行"
fi

echo ""
echo "================================"
echo "✅ 数据库连接测试完成！"
echo ""
echo "💡 提示:"
echo "- 如果连接失败，请检查 PostgreSQL 服务状态"
echo "- 确保数据库 webbleen_api 已创建"
echo "- 使用 './start.sh -d' 启动 Docker 服务进行测试"
