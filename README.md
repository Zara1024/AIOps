# 智维云枢 (CloudOps Hub)

面向中大型企业的一站式 **AI 驱动智能运维管理平台**，集成资产管理、Kubernetes 编排、监控告警、自动化任务、AI 智能分析等核心运维能力。

## 技术栈

### 后端
- **语言**: Go 1.24+
- **框架**: Gin + GORM + gRPC
- **数据库**: PostgreSQL 16 + Redis 7
- **消息队列**: NATS JetStream

### 前端
- **框架**: Vue 3.5+ (TypeScript)
- **构建**: Vite 6
- **UI**: Element Plus
- **状态管理**: Pinia

### AI
- **LLM 适配**: OpenAI / Claude / 通义千问 / DeepSeek / Ollama
- **核心引擎**: Agentic Loop + Tool Calling
- **知识库**: pgvector RAG

## 快速开发

```bash
# 后端
cd cloudops-server
go run cmd/server/main.go

# 前端
cd cloudops-web
npm install
npm run dev
```

## 默认账号

- 用户名: `admin`
- 密码: `Admin@2026`

## 项目结构

```
AIOps/
├── cloudops-server/     # Go 后端
│   ├── cmd/server/      # 启动入口
│   ├── internal/        # 内部模块
│   │   ├── auth/        # 认证模块
│   │   ├── system/      # 系统管理
│   │   └── router/      # 路由注册
│   └── pkg/             # 公共工具包
├── cloudops-web/        # Vue3 前端
└── docker-compose.dev.yml
```

## 开源协议

MIT License
