# Railway 部署指南

## 部署步骤

### 1. 准备数据库

Railway 提供了 PostgreSQL 服务，你需要：

1. 在 Railway 控制台创建一个新的 PostgreSQL 服务
2. 获取数据库连接信息（会自动生成 `DATABASE_URL` 环境变量）

### 2. 部署应用

#### 方法一：通过 Railway CLI

```bash
# 安装 Railway CLI
npm install -g @railway/cli

# 登录 Railway
railway login

# 初始化项目
railway init

# 连接数据库服务
railway add postgresql

# 设置环境变量
railway variables set RUN_MODE=release

# 部署
railway up
```

#### 方法二：通过 GitHub 集成

1. 将代码推送到 GitHub 仓库
2. 在 Railway 控制台选择 "Deploy from GitHub repo"
3. 选择你的仓库
4. 添加 PostgreSQL 服务
5. 设置环境变量

### 3. 环境变量配置

在 Railway 控制台中设置以下环境变量：

```
DATABASE_URL=postgres://username:password@host:port/database?sslmode=require
RUN_MODE=release
```

**注意**：Railway 会自动提供 `DATABASE_URL`，你只需要设置 `RUN_MODE`。

### 4. 域名配置

Railway 会自动为你的服务分配一个域名，格式如：
`https://your-service-name-production.up.railway.app`

你也可以在 Railway 控制台中配置自定义域名。

## 监控和日志

- **日志查看**：在 Railway 控制台的 "Deployments" 页面查看实时日志
- **监控**：Railway 提供基本的监控指标
- **重启**：可以在控制台手动重启服务

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 `DATABASE_URL` 环境变量是否正确
   - 确保 PostgreSQL 服务正在运行

2. **端口问题**
   - Railway 会自动设置 `PORT` 环境变量
   - 确保应用监听正确的端口

3. **构建失败**
   - 检查 Dockerfile 语法
   - 查看构建日志中的错误信息

### 调试命令

```bash
# 查看服务状态
railway status

# 查看日志
railway logs

# 连接到服务
railway shell
```

## 成本优化

- Railway 提供免费额度，适合小项目
- 监控资源使用情况，避免超出免费限制
- 考虑使用 Railway 的付费计划以获得更多资源

## 安全建议

1. **环境变量**：不要在代码中硬编码敏感信息
2. **数据库**：使用强密码和 SSL 连接
3. **HTTPS**：Railway 自动提供 HTTPS 支持
4. **访问控制**：合理设置服务访问权限

## 更新部署

```bash
# 推送代码到 GitHub
git push origin main

# Railway 会自动检测更改并重新部署
# 或者手动触发部署
railway up
```

## 备份

- Railway 提供数据库自动备份
- 定期导出重要数据
- 考虑设置数据备份策略
