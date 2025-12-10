# HaloLight API Go

[![License](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8.svg?logo=go)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.10-008ECF.svg)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-1.25-00ADD8.svg)](https://gorm.io/)

HaloLight 后台管理系统的 **Go 高性能后端 API 服务**，基于 Gin + GORM 构建的 RESTful API。

- API 文档: <https://api-go.halolight.h7ml.cn>
- GitHub: <https://github.com/halolight/halolight-api-go>

## 特性

- **Gin Web Framework** - 高性能 HTTP 框架
- **GORM 2** - 强大的 ORM 库
- **JWT 双令牌认证** - Access Token + Refresh Token
- **RBAC 权限系统** - 基于角色的访问控制
- **PostgreSQL 16** - 生产级数据库
- **12 个业务模块** - 完整的后台管理 API
- **90+ RESTful 端点** - 覆盖常见业务场景
- **Docker 支持** - 容器化部署
- **完善的错误处理** - 统一的错误响应格式
- **输入验证** - 基于 Gin binding 的自动验证

## 技术栈

- **Go** 1.21+
- **Gin** 1.10 - Web 框架
- **GORM** 1.25 - ORM
- **PostgreSQL** 15+ - 数据库
- **JWT** (golang-jwt/jwt/v5) - 认证
- **bcrypt** - 密码哈希

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- PostgreSQL 15+
- Docker (可选)

### 安装

```bash
# 克隆仓库
git clone https://github.com/halolight/halolight-api-go.git
cd halolight-api-go

# 安装依赖
go mod tidy

# 复制环境变量文件
cp .env.example .env

# 编辑 .env 文件，配置数据库连接等信息
# vim .env

# 启动开发服务器
make dev

# 或直接运行
go run ./cmd/server
```

### 使用 Docker

```bash
# 启动所有服务（API + PostgreSQL）
make docker-up

# 查看日志
make docker-logs

# 停止服务
make docker-down
```

## 可用命令

```bash
make dev            # 启动开发服务器
make build          # 编译二进制文件
make run            # 编译并运行
make test           # 运行测试
make test-coverage  # 运行测试并生成覆盖率报告
make lint           # 运行代码检查
make tidy           # 整理 Go 模块依赖
make docker-build   # 构建 Docker 镜像
make docker-up      # 启动 Docker 容器
make docker-down    # 停止 Docker 容器
make docker-logs    # 查看 Docker 日志
make clean          # 清理构建产物
make help           # 显示帮助信息
```

## 项目结构

```
halolight-api-go/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── handlers/                # HTTP 处理器
│   │   ├── auth_handler.go      # 认证处理器
│   │   └── user_handler.go      # 用户处理器
│   ├── models/                  # 数据模型
│   │   └── user.go              # 用户模型
│   ├── services/                # 业务逻辑层
│   │   ├── auth_service.go      # 认证服务
│   │   └── user_service.go      # 用户服务
│   ├── middleware/              # 中间件
│   │   ├── auth.go              # JWT 认证中间件
│   │   └── cors.go              # CORS 中间件
│   ├── repository/              # 数据访问层
│   │   └── user_repository.go   # 用户数据仓库
│   └── routes/                  # 路由定义
│       └── router.go            # 路由配置
├── pkg/
│   ├── config/                  # 配置管理
│   │   └── config.go
│   ├── database/                # 数据库连接
│   │   └── database.go
│   └── utils/                   # 工具函数
│       ├── jwt.go               # JWT 工具
│       └── hash.go              # 密码哈希工具
├── .env.example                 # 环境变量示例
├── .gitignore
├── Dockerfile                   # Docker 配置
├── docker-compose.yml           # Docker Compose 配置
├── Makefile                     # Make 命令
├── go.mod                       # Go 模块定义
└── README.md
```

## API 端点

### 认证 (Public)

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/auth/register` | 用户注册 |
| POST | `/api/auth/login` | 用户登录 |
| POST | `/api/auth/refresh` | 刷新令牌 |
| POST | `/api/auth/forgot-password` | 忘记密码 |
| POST | `/api/auth/reset-password` | 重置密码 |

### 认证 (Protected)

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/auth/me` | 获取当前用户 |
| POST | `/api/auth/logout` | 登出 |

### 用户管理 (Protected)

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/users` | 用户列表（分页） |
| GET | `/api/users/:id` | 用户详情 |
| POST | `/api/users` | 创建用户 |
| PATCH | `/api/users/:id` | 更新用户 |
| PATCH | `/api/users/:id/status` | 更新状态 |
| POST | `/api/users/batch-delete` | 批量删除 |
| DELETE | `/api/users/:id` | 删除用户 |

### 其他模块 (Protected)

- **Roles** (`/api/roles`) - 角色 CRUD + 权限分配
- **Permissions** (`/api/permissions`) - 权限 CRUD
- **Teams** (`/api/teams`) - 团队 CRUD + 成员管理
- **Documents** (`/api/documents`) - 文档 CRUD + 分享/标签
- **Files** (`/api/files`) - 文件上传/下载/管理
- **Folders** (`/api/folders`) - 文件夹树形结构
- **Calendar** (`/api/calendar/events`) - 日历事件管理
- **Notifications** (`/api/notifications`) - 通知管理
- **Messages** (`/api/messages`) - 消息会话
- **Dashboard** (`/api/dashboard`) - 统计数据

### 健康检查

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/health` | 健康检查 |

## API 示例

### 注册

```bash
curl -X POST http://localhost:8000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "password123"
  }'
```

### 登录

```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### 获取用户列表（需要认证）

```bash
curl -X GET http://localhost:8000/api/users?page=1&page_size=20 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 环境变量

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `APP_ENV` | 应用环境 | `development` |
| `APP_PORT` | 服务端口 | `8000` |
| `JWT_SECRET` | JWT 密钥 | `change-me-in-production` |
| `JWT_EXPIRE_MINUTES` | JWT 过期时间（分钟） | `60` |
| `DB_HOST` | 数据库主机 | `localhost` |
| `DB_PORT` | 数据库端口 | `5432` |
| `DB_USER` | 数据库用户 | `postgres` |
| `DB_PASSWORD` | 数据库密码 | `postgres` |
| `DB_NAME` | 数据库名称 | `halolight` |
| `DB_SSLMODE` | SSL 模式 | `disable` |

## 架构设计

### 分层架构

```
┌─────────────────────────────────────┐
│         HTTP Handlers               │  ← 请求入口，参数验证
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│          Services                   │  ← 业务逻辑，事务管理
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│        Repository                   │  ← 数据访问，CRUD 操作
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│         Database (GORM)             │  ← PostgreSQL
└─────────────────────────────────────┘
```

### 认证流程

1. 用户登录 → 验证凭据
2. 生成 JWT token（包含 user_id）
3. 客户端在后续请求中携带 token
4. AuthMiddleware 验证 token
5. 从 token 提取 user_id 并注入到 context
6. 业务逻辑可通过 context 获取当前用户

## 开发指南

### 添加新的 API 端点

1. **定义模型** (`internal/models/`)
2. **创建 Repository** (`internal/repository/`)
3. **实现 Service** (`internal/services/`)
4. **添加 Handler** (`internal/handlers/`)
5. **注册路由** (`internal/routes/router.go`)

### 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 查看覆盖率报告（会在浏览器打开）
open coverage.html
```

### 代码规范

```bash
# 格式化代码
go fmt ./...

# 运行 linter（需要安装 golangci-lint）
make lint
```

## 部署

### Docker 部署

```bash
# 构建镜像
make docker-build

# 使用 docker-compose 启动
make docker-up
```

### 二进制部署

```bash
# 编译生产版本
make build

# 运行二进制
./bin/halolight-api
```

### 环境变量配置

生产环境建议通过系统环境变量或配置管理工具（如 Kubernetes ConfigMap/Secret）注入配置，而不是使用 `.env` 文件。

## 安全最佳实践

- ✅ 密码使用 bcrypt 哈希（cost=10）
- ✅ JWT 签名验证（HS256）
- ✅ CORS 中间件配置
- ✅ 输入验证（Gin binding）
- ✅ SQL 注入防护（GORM 参数化查询）
- ✅ 错误信息不泄露敏感数据
- ⚠️ 生产环境务必修改 `JWT_SECRET`
- ⚠️ 生产环境建议使用 HTTPS
- ⚠️ 考虑添加速率限制中间件

## 性能

- **并发处理**: Gin 基于 goroutine，支持高并发
- **数据库连接池**: GORM 自动管理连接池
- **内存占用**: ~20-30MB（空载）
- **响应时间**: < 10ms (P99, 无数据库查询)

## 故障排查

### 常见问题

**问题**: 无法连接数据库

```bash
# 检查 PostgreSQL 是否运行
docker-compose ps

# 查看数据库日志
docker-compose logs db
```

**问题**: JWT token 无效

- 确认 `JWT_SECRET` 配置正确
- 检查 token 是否过期
- 验证 Authorization header 格式: `Bearer <token>`

**问题**: 依赖下载失败

```bash
# 设置 Go 代理（中国用户）
export GOPROXY=https://goproxy.cn,direct

# 重新下载
go mod tidy
```

## 相关链接

- [HaloLight 文档](https://halolight.docs.h7ml.cn)
- [Gin 文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Go JWT 文档](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

[MIT](LICENSE)

---

Made with ❤️ by HaloLight Team
