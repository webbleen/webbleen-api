# Webbleen-API

这是一个为 Webbleen 博客网站专门设计的精简统计 API 服务。

## 功能特性

- **访问统计**: 记录和统计网站访问量
- **页面统计**: 跟踪页面浏览量和热门页面
- **用户行为分析**: 分析用户设备、浏览器、地理位置等
- **内容统计**: 统计博客文章、标签、分类数量
- **趋势分析**: 提供访问趋势和日统计数据

## API 接口

### 记录访问
```
POST /api/visit
```
参数：
- `page`: 页面路径
- `session_id`: 会话ID
- `country`: 国家
- `city`: 城市
- `device`: 设备类型
- `browser`: 浏览器
- `os`: 操作系统

### 获取访问统计
```
GET /api/stats/visits
```
返回今日访问量、累计访问量、页面浏览量等

### 获取内容统计
```
GET /api/stats/content
```
返回文章数量、标签数量、分类数量等

### 获取热门页面
```
GET /api/stats/pages?limit=10
```
返回访问量最高的页面列表

### 获取访问趋势
```
GET /api/stats/trend?days=30
```
返回指定天数的访问趋势数据

### 获取用户行为分析
```
GET /api/stats/behavior
```
返回设备、浏览器、操作系统、地理位置等统计

### 获取日统计
```
GET /api/stats/daily?limit=30
```
返回每日统计数据

### 更新内容统计
```
POST /api/stats/content
```
参数：
- `articles`: 文章数量
- `tags`: 标签数量
- `categories`: 分类数量

## 配置

### 环境变量（必需）

设置以下环境变量：

```bash
# PostgreSQL 数据库连接
export DATABASE_URL="postgres://username:password@localhost:5432/webbleen_api?sslmode=disable"
export RUN_MODE="debug"
```

**注意**：`DATABASE_URL` 环境变量是必需的，不再支持配置文件方式。

## 运行

1. 确保 PostgreSQL 数据库运行
2. 创建数据库 `webbleen_api`
3. 设置环境变量并运行服务：

```bash
# 设置环境变量
export DATABASE_URL="postgres://username:password@localhost:5432/webbleen_api?sslmode=disable"

# 运行服务
go run main.go
```

或构建后运行：

```bash
go build -o webbleen-api main.go
./webbleen-api
```

## 数据库表结构

- `visit_record`: 访问记录表
- `page_view`: 页面访问统计表
- `daily_stats`: 日统计表
- `content_stats`: 内容统计表

## 与 Hugo 博客集成

在 Hugo 博客中添加统计脚本：

```javascript
// 记录页面访问
fetch('/api/visit', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
        page: window.location.pathname,
        session_id: getSessionId(),
        // 其他参数...
    })
});
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t webbleen-api .

# 运行容器
docker run -d \
  --name webbleen-api \
  -p 8000:8000 \
  -e DATABASE_URL="postgres://username:password@host.docker.internal:5432/webbleen_api?sslmode=disable" \
  webbleen-api
```

### Docker Compose 部署

创建 `docker-compose.yml`：

```yaml
version: '3.8'
services:
  webbleen-api:
    build: .
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=root:root@tcp(mysql:3306)/webbleen_api?charset=utf8&parseTime=True&loc=Local
      - TABLE_PREFIX=stats_
      - RUN_MODE=release
    depends_on:
      - mysql
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    environment:
      - MYSQL_ROOT_PASSWORD=root
      - MYSQL_DATABASE=webbleen_api
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

运行：

```bash
docker-compose up -d
```

## Railway 部署

### 快速部署

1. **准备数据库**：在 Railway 控制台创建 PostgreSQL 服务
2. **部署应用**：连接 GitHub 仓库或使用 Railway CLI
3. **设置环境变量**：
   ```
   RUN_MODE=release
   ```
   （`DATABASE_URL` 会自动提供）

### 详细步骤

参考 [RAILWAY_DEPLOYMENT.md](./RAILWAY_DEPLOYMENT.md) 获取完整的部署指南。

### Railway CLI 部署

```bash
# 安装 Railway CLI
npm install -g @railway/cli

# 登录并部署
railway login
railway init
railway add postgresql
railway up
```