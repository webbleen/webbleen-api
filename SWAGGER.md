# Swagger API 文档

## 📖 概述

Webbleen 博客 API 服务提供了完整的 Swagger 文档，方便开发者查看和测试 API 接口。

## 🔗 访问地址

- **Swagger UI**: `https://api.webbleen.com/swagger/index.html`
- **JSON 格式**: `https://api.webbleen.com/swagger/doc.json`
- **YAML 格式**: `https://api.webbleen.com/swagger/doc.yaml`

## 📋 API 接口列表

### 统计相关接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/stats/visit` | POST | 记录访问 |
| `/stats/visits` | GET | 获取访问统计概览 |
| `/stats/content` | GET | 获取内容统计 |
| `/stats/content` | POST | 更新内容统计 |
| `/stats/pages` | GET | 获取热门页面 |
| `/stats/trend` | GET | 获取访问趋势 |
| `/stats/behavior` | GET | 获取用户行为分析 |
| `/stats/daily` | GET | 获取日统计 |

## 🛠️ 本地开发

### 生成 Swagger 文档

```bash
# 安装 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
swag init

# 运行服务
go run main.go
```

### 访问本地 Swagger UI

```
http://localhost:8000/swagger/index.html
```

## 📝 文档更新

当修改 API 接口或添加新的接口时，需要重新生成 Swagger 文档：

```bash
swag init
```

## 🔧 配置说明

Swagger 配置在 `main.go` 中定义：

```go
// @title Webbleen 博客 API 服务
// @version 1.0
// @description Webbleen 博客 API 服务，提供统计、内容管理等功能
// @host api.webbleen.com
// @BasePath /
// @schemes https
```

## 📊 数据模型

### VisitRecord - 访问记录
- `ip`: 访问者 IP 地址
- `user_agent`: 用户代理
- `page`: 访问页面
- `session_id`: 会话 ID
- `device`: 设备类型
- `browser`: 浏览器
- `os`: 操作系统
- `country`: 国家
- `city`: 城市

### ContentStats - 内容统计
- `total_articles`: 文章总数
- `total_tags`: 标签总数
- `total_categories`: 分类总数
- `last_update`: 最后更新时间

## 🚀 部署说明

Swagger 文档会自动包含在部署的 API 服务中，无需额外配置。

## 📞 支持

如有问题，请联系：
- 邮箱: contact@webbleen.com
- 网站: https://webbleen.com
