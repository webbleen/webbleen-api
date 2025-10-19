-- webbleen-api 数据库初始化脚本
-- 用于 PostgreSQL 数据库初始化

-- 创建数据库（如果不存在）
-- 注意：这通常需要超级用户权限
-- CREATE DATABASE webbleen_api;

-- 连接到 webbleen_api 数据库
\c webbleen_api;

-- 创建用户（如果不存在）
-- 注意：这通常需要超级用户权限
-- CREATE USER postgres WITH PASSWORD 'password';

-- 授予权限
-- GRANT ALL PRIVILEGES ON DATABASE webbleen_api TO postgres;

-- 创建扩展（如果需要）
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建表（这些表将由 GORM 自动创建，这里仅作为参考）

-- 访问统计表
CREATE TABLE IF NOT EXISTS visit_records (
    id SERIAL PRIMARY KEY,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    referer TEXT,
    page_url TEXT NOT NULL,
    page_title TEXT,
    session_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 工具表
CREATE TABLE IF NOT EXISTS tools (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    icon VARCHAR(255),
    url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 工具使用统计表
CREATE TABLE IF NOT EXISTS tool_usages (
    id SERIAL PRIMARY KEY,
    tool_id INTEGER REFERENCES tools(id),
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- API 密钥表
CREATE TABLE IF NOT EXISTS api_keys (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 聊天会话表
CREATE TABLE IF NOT EXISTS chat_sessions (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 聊天消息表
CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'assistant')),
    content TEXT NOT NULL,
    type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES chat_sessions(session_id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_visit_records_created_at ON visit_records(created_at);
CREATE INDEX IF NOT EXISTS idx_visit_records_ip_address ON visit_records(ip_address);
CREATE INDEX IF NOT EXISTS idx_tool_usages_tool_id ON tool_usages(tool_id);
CREATE INDEX IF NOT EXISTS idx_tool_usages_created_at ON tool_usages(created_at);
CREATE INDEX IF NOT EXISTS idx_chat_messages_session_id ON chat_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(created_at);

-- 插入示例数据（可选）
INSERT INTO tools (name, description, category, icon, url) VALUES
('AI 聊天助手', '智能聊天助手，支持代码审查和分析', 'AI', '🤖', '/chat'),
('代码审查', 'AI 驱动的代码质量检查工具', 'AI', '🔍', '/api/ai/chat'),
('架构分析', '深度代码架构分析工具', 'AI', '🏗️', '/api/ai/chat')
ON CONFLICT DO NOTHING;

-- 显示表信息
\dt
