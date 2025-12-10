# CLAUDE.md

This file provides guidance to Claude Code when working with this repository.

## 项目概述

HaloLight Go 高性能后端 API 服务，基于 Gin + GORM 构建，提供完整的 RESTful API，包含 12 个业务模块。

## 技术栈

- **框架**: Gin 1.10
- **ORM**: GORM 2
- **数据库**: PostgreSQL 16
- **认证**: JWT 双令牌 (golang-jwt/jwt/v5)
- **验证**: go-playground/validator
- **ID 生成**: ULID
- **运行时**: Go 1.22+

## 常用命令

```bash
make dev          # Air 热重载开发
make build        # 编译
make run          # 运行
make test         # 测试
make test-coverage # 测试覆盖率
make lint         # golangci-lint 检查
make swagger      # 生成 Swagger
make docker-up    # Docker Compose 启动
make docker-down  # Docker Compose 停止
```

## 项目结构

```
├── cmd/server/           # 入口
│   └── main.go
├── internal/
│   ├── handlers/         # HTTP 处理器 (12 个模块)
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── role_handler.go
│   │   ├── permission_handler.go
│   │   ├── team_handler.go
│   │   ├── document_handler.go
│   │   ├── file_handler.go
│   │   ├── folder_handler.go
│   │   ├── calendar_handler.go
│   │   ├── notification_handler.go
│   │   ├── message_handler.go
│   │   └── dashboard_handler.go
│   ├── models/           # 数据模型 (17+ 模型)
│   ├── services/         # 业务逻辑层
│   ├── middleware/       # 中间件
│   │   ├── auth.go       # JWT 认证
│   │   └── cors.go       # CORS
│   ├── repository/       # 数据访问层
│   └── routes/           # 路由配置
│       └── router.go
└── pkg/
    ├── config/           # 配置管理
    ├── database/         # 数据库连接
    └── utils/            # 工具函数
        ├── jwt.go
        └── hash.go
```

## API 模块

| 模块 | 路由前缀 | 端点数 | 说明 |
|------|----------|--------|------|
| Auth | `/api/auth` | 7 | 登录、注册、刷新令牌、登出、获取当前用户、忘记/重置密码 |
| Users | `/api/users` | 7 | CRUD、分页搜索、状态更新、批量删除 |
| Roles | `/api/roles` | 6 | CRUD、权限分配 |
| Permissions | `/api/permissions` | 4 | CRUD |
| Teams | `/api/teams` | 7 | CRUD、成员管理 |
| Documents | `/api/documents` | 11 | CRUD、分享、标签、移动 |
| Files | `/api/files` | 14 | 上传、下载、存储信息、移动、复制、收藏 |
| Folders | `/api/folders` | 5 | 树形结构管理 |
| Calendar | `/api/calendar/events` | 9 | 事件、参会人管理 |
| Notifications | `/api/notifications` | 5 | 通知管理 |
| Messages | `/api/messages` | 5 | 会话、消息 |
| Dashboard | `/api/dashboard` | 9 | 统计数据 |

## 认证机制

### JWT 双令牌

```go
// 登录返回
{
    "accessToken": "eyJ...",   // JWT_EXPIRE_MINUTES 分钟有效
    "refreshToken": "eyJ...",  // 30天有效
    "user": { ... }
}

// 刷新令牌
POST /api/auth/refresh
{ "refreshToken": "eyJ..." }
```

### 认证中间件

```go
// 使用 AuthMiddleware
users := api.Group("/users")
users.Use(middleware.AuthMiddleware(cfg))
{
    users.GET("", userHandler.List)
}

// 从 context 获取用户 ID
userID := c.GetString("userID")
```

## 统一响应格式

```go
// 成功响应
gin.H{
    "success": true,
    "data":    data,
    "message": "操作成功",
}

// 错误响应
gin.H{
    "success": false,
    "message": "错误信息",
}

// 分页响应
gin.H{
    "success": true,
    "data":    items,
    "meta": gin.H{
        "total":      total,
        "page":       page,
        "limit":      limit,
        "totalPages": totalPages,
    },
}
```

## 环境变量

```bash
# 应用配置
APP_ENV=development
APP_PORT=8000

# JWT 配置
JWT_SECRET=your-super-secret-key
JWT_EXPIRE_MINUTES=60

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=halolight
DB_SSLMODE=disable
```

## 数据库模型

主要模型包括：
- **User** - 用户（含角色、权限关联）
- **Role/Permission/RolePermission** - RBAC 权限系统
- **RefreshToken** - 刷新令牌存储
- **Team/TeamMember** - 团队管理
- **Document/DocumentShare/Tag** - 文档系统
- **File/Folder** - 文件系统
- **CalendarEvent/EventAttendee/EventReminder** - 日历系统
- **Conversation/ConversationParticipant/Message** - 消息系统
- **Notification** - 通知系统
- **ActivityLog** - 活动日志

## 开发规范

1. **分层架构**: Handler → Service → Repository → Database
2. **依赖注入**: 通过构造函数注入依赖
3. **错误处理**: 使用 gin.H 返回统一格式
4. **验证**: 使用 Gin binding tag 验证请求
5. **ID 生成**: 使用 ULID 生成唯一 ID
