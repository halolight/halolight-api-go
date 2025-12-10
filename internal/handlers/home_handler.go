package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/halolight/halolight-api-go/pkg/config"
)

// HomeHandler handles home page and health check
type HomeHandler struct {
	cfg config.Config
}

// NewHomeHandler creates a new home handler
func NewHomeHandler(cfg config.Config) *HomeHandler {
	return &HomeHandler{cfg: cfg}
}

// Home returns a beautiful HTML homepage
func (h *HomeHandler) Home(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, h.getHomePage())
}

// Health returns health check response
func (h *HomeHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":        "HaloLight API",
		"version":     "1.0.0",
		"environment": h.cfg.AppEnv,
		"status":      "ok",
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"database":    "Use /health/database to check database connection",
	})
}

func (h *HomeHandler) getHomePage() string {
	return `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="description" content="HaloLight API - 基于 Gin + GORM 的高性能 Go 后端服务">
  <meta name="keywords" content="Go, Gin, API, GORM, PostgreSQL, JWT, RBAC">
  <title>HaloLight API | Go Backend Service</title>
  <link rel="icon" href="data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>🔵</text></svg>">
  <script src="https://cdn.tailwindcss.com"></script>
  <script>
    tailwind.config = {
      theme: {
        extend: {
          colors: {
            primary: '#00ADD8',
            secondary: '#0891b2',
            accent: '#6366f1',
          }
        }
      }
    }
  </script>
  <style>
    :root {
      --primary: #00ADD8;
      --primary-dark: #0891b2;
      --secondary: #0891b2;
      --accent: #6366f1;
      --bg-dark: #0f172a;
      --gradient: linear-gradient(135deg, var(--primary) 0%, var(--secondary) 50%, var(--accent) 100%);
    }
    .bg-gradient-animated::before {
      content: '';
      position: absolute;
      top: -50%;
      left: -50%;
      width: 200%;
      height: 200%;
      background: radial-gradient(circle at 30% 30%, rgba(0, 173, 216, 0.15) 0%, transparent 50%),
                  radial-gradient(circle at 70% 70%, rgba(8, 145, 178, 0.1) 0%, transparent 50%),
                  radial-gradient(circle at 50% 50%, rgba(99, 102, 241, 0.05) 0%, transparent 50%);
      animation: rotate 30s linear infinite;
    }
    @keyframes rotate {
      from { transform: rotate(0deg); }
      to { transform: rotate(360deg); }
    }
    .text-gradient {
      background: var(--gradient);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
    .bg-gradient-brand { background: var(--gradient); }
    .btn-gradient {
      background: var(--gradient);
      box-shadow: 0 4px 14px rgba(0, 173, 216, 0.4);
    }
    .btn-gradient:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(0, 173, 216, 0.5);
    }
    .card-hover:hover {
      border-color: var(--primary);
      transform: translateY(-4px);
      box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
    }
    .module-hover:hover {
      border-color: var(--primary);
      background: rgba(0, 173, 216, 0.1);
    }
  </style>
</head>
<body class="bg-slate-900 text-slate-50 min-h-screen overflow-x-hidden font-sans">
  <div class="fixed inset-0 bg-slate-900 -z-10 bg-gradient-animated"></div>

  <!-- Navigation -->
  <nav class="fixed top-0 left-0 right-0 z-50 py-4 backdrop-blur-xl bg-slate-900/80 border-b border-slate-700/50">
    <div class="max-w-7xl mx-auto px-6 flex justify-between items-center">
      <div class="text-2xl font-bold text-gradient">🔵 HaloLight API</div>
      <div class="hidden md:flex items-center gap-6">
        <a href="#features" class="text-slate-400 hover:text-white text-sm font-medium transition-colors">Features</a>
        <a href="#modules" class="text-slate-400 hover:text-white text-sm font-medium transition-colors">Modules</a>
        <a href="/docs" class="text-slate-400 hover:text-white text-sm font-medium transition-colors">API Docs</a>
        <a href="https://github.com/halolight/halolight-api-go" target="_blank" class="text-slate-400 hover:text-white text-sm font-medium transition-colors">GitHub</a>
        <span class="px-3 py-1 text-xs font-semibold rounded-full bg-cyan-500/20 text-cyan-400 border border-cyan-500/30">v1.0.0</span>
      </div>
    </div>
  </nav>

  <!-- Hero Section -->
  <section class="min-h-screen flex items-center pt-20">
    <div class="max-w-7xl mx-auto px-6">
      <div class="max-w-3xl">
        <div class="inline-flex items-center gap-2 px-4 py-2 bg-slate-800/80 border border-slate-700/50 rounded-full text-sm text-slate-400 mb-6">
          <span class="text-cyan-500">🔵</span> High-Performance Go Backend Service
        </div>
        <h1 class="text-4xl md:text-6xl font-extrabold leading-tight mb-6">
          高性能 Go API<br>
          <span class="text-gradient">Gin + GORM 驱动</span>
        </h1>
        <p class="text-xl text-slate-400 leading-relaxed mb-8">
          基于 Gin 1.10 + GORM 2 的高性能后端服务，提供完整的 JWT 双令牌认证、RBAC 权限管理，
          90+ RESTful API 端点开箱即用，支持高并发场景。
        </p>
        <div class="flex flex-col sm:flex-row gap-4 mb-12">
          <a href="/docs" class="btn-gradient inline-flex items-center justify-center gap-2 px-7 py-4 text-white font-semibold rounded-xl transition-all">
            📖 查看 API 文档
          </a>
          <a href="https://halolight.docs.h7ml.cn/guide/api-go" class="inline-flex items-center justify-center gap-2 px-7 py-4 bg-slate-800/80 text-white font-semibold rounded-xl border border-slate-700/50 hover:border-primary hover:bg-slate-800 transition-all" target="_blank">
            📚 在线使用指南
          </a>
          <a href="/health" class="inline-flex items-center justify-center gap-2 px-7 py-4 bg-slate-800/80 text-white font-semibold rounded-xl border border-slate-700/50 hover:border-primary hover:bg-slate-800 transition-all">
            💚 健康检查
          </a>
        </div>
        <!-- Tech Stack -->
        <div class="flex flex-wrap gap-3 pt-8 border-t border-slate-700/50">
          <div class="flex items-center gap-2 px-4 py-2 bg-slate-800/80 border border-slate-700/50 rounded-lg text-sm text-slate-400">
            <span>🔵</span> Go 1.22
          </div>
          <div class="flex items-center gap-2 px-4 py-2 bg-slate-800/80 border border-slate-700/50 rounded-lg text-sm text-slate-400">
            <span>⚡</span> Gin 1.10
          </div>
          <div class="flex items-center gap-2 px-4 py-2 bg-slate-800/80 border border-slate-700/50 rounded-lg text-sm text-slate-400">
            <span>🗄️</span> GORM 2
          </div>
          <div class="flex items-center gap-2 px-4 py-2 bg-slate-800/80 border border-slate-700/50 rounded-lg text-sm text-slate-400">
            <span>🐘</span> PostgreSQL 16
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Stats Section -->
  <section class="py-16">
    <div class="max-w-7xl mx-auto px-6">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div class="text-center p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl">
          <div class="text-5xl font-extrabold text-gradient mb-2">12</div>
          <div class="text-slate-400">业务模块</div>
        </div>
        <div class="text-center p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl">
          <div class="text-5xl font-extrabold text-gradient mb-2">90+</div>
          <div class="text-slate-400">API 端点</div>
        </div>
        <div class="text-center p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl">
          <div class="text-5xl font-extrabold text-gradient mb-2">&lt;10ms</div>
          <div class="text-slate-400">响应时间</div>
        </div>
        <div class="text-center p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl">
          <div class="text-5xl font-extrabold text-gradient mb-2">MIT</div>
          <div class="text-slate-400">开源协议</div>
        </div>
      </div>
    </div>
  </section>

  <!-- Features Section -->
  <section id="features" class="py-24">
    <div class="max-w-7xl mx-auto px-6">
      <div class="text-center mb-16">
        <h2 class="text-4xl font-bold mb-4">核心特性</h2>
        <p class="text-slate-400 text-lg max-w-2xl mx-auto">企业级架构设计，高性能开箱即用</p>
      </div>
      <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div class="p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl transition-all card-hover">
          <div class="w-12 h-12 flex items-center justify-center bg-gradient-brand rounded-xl text-2xl mb-5">🔐</div>
          <h3 class="text-xl font-semibold mb-3">JWT 双令牌认证</h3>
          <p class="text-slate-400 leading-relaxed">AccessToken + RefreshToken 机制，支持自动刷新，安全可靠的身份验证方案。</p>
        </div>
        <div class="p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl transition-all card-hover">
          <div class="w-12 h-12 flex items-center justify-center bg-gradient-brand rounded-xl text-2xl mb-5">🛡️</div>
          <h3 class="text-xl font-semibold mb-3">RBAC 权限控制</h3>
          <p class="text-slate-400 leading-relaxed">基于角色的访问控制，支持通配符权限（users:*, *），灵活的权限管理。</p>
        </div>
        <div class="p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl transition-all card-hover">
          <div class="w-12 h-12 flex items-center justify-center bg-gradient-brand rounded-xl text-2xl mb-5">⚡</div>
          <h3 class="text-xl font-semibold mb-3">高性能并发</h3>
          <p class="text-slate-400 leading-relaxed">基于 Goroutine 的高并发处理，轻松应对大流量场景，响应时间 &lt;10ms。</p>
        </div>
        <div class="p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl transition-all card-hover">
          <div class="w-12 h-12 flex items-center justify-center bg-gradient-brand rounded-xl text-2xl mb-5">🗄️</div>
          <h3 class="text-xl font-semibold mb-3">GORM ORM</h3>
          <p class="text-slate-400 leading-relaxed">强大的 ORM 库，自动迁移、关联查询、软删除，数据库操作简单高效。</p>
        </div>
        <div class="p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl transition-all card-hover">
          <div class="w-12 h-12 flex items-center justify-center bg-gradient-brand rounded-xl text-2xl mb-5">✅</div>
          <h3 class="text-xl font-semibold mb-3">Gin Binding 验证</h3>
          <p class="text-slate-400 leading-relaxed">使用 Gin binding 自动验证请求数据，确保数据完整性和安全性。</p>
        </div>
        <div class="p-8 bg-slate-800/50 border border-slate-700/50 rounded-2xl transition-all card-hover">
          <div class="w-12 h-12 flex items-center justify-center bg-gradient-brand rounded-xl text-2xl mb-5">🐳</div>
          <h3 class="text-xl font-semibold mb-3">Docker 部署</h3>
          <p class="text-slate-400 leading-relaxed">多阶段构建优化，Docker Compose 一键部署，内置健康检查机制。</p>
        </div>
      </div>
    </div>
  </section>

  <!-- Modules Section -->
  <section id="modules" class="py-24">
    <div class="max-w-7xl mx-auto px-6">
      <div class="text-center mb-16">
        <h2 class="text-4xl font-bold mb-4">API 模块</h2>
        <p class="text-slate-400 text-lg max-w-2xl mx-auto">12 个核心业务模块，覆盖常见企业应用场景</p>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
        <a href="/docs#/Auth" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">🔑</div>
          <div><h4 class="font-semibold text-white">Auth</h4><span class="text-sm text-slate-400">7 个端点</span></div>
        </a>
        <a href="/docs#/Users" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">👥</div>
          <div><h4 class="font-semibold text-white">Users</h4><span class="text-sm text-slate-400">7 个端点</span></div>
        </a>
        <a href="/docs#/Roles" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">🎭</div>
          <div><h4 class="font-semibold text-white">Roles</h4><span class="text-sm text-slate-400">6 个端点</span></div>
        </a>
        <a href="/docs#/Permissions" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">🔒</div>
          <div><h4 class="font-semibold text-white">Permissions</h4><span class="text-sm text-slate-400">4 个端点</span></div>
        </a>
        <a href="/docs#/Teams" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">👨‍👩‍👧‍👦</div>
          <div><h4 class="font-semibold text-white">Teams</h4><span class="text-sm text-slate-400">7 个端点</span></div>
        </a>
        <a href="/docs#/Documents" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">📄</div>
          <div><h4 class="font-semibold text-white">Documents</h4><span class="text-sm text-slate-400">11 个端点</span></div>
        </a>
        <a href="/docs#/Files" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">📁</div>
          <div><h4 class="font-semibold text-white">Files</h4><span class="text-sm text-slate-400">14 个端点</span></div>
        </a>
        <a href="/docs#/Folders" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">📂</div>
          <div><h4 class="font-semibold text-white">Folders</h4><span class="text-sm text-slate-400">5 个端点</span></div>
        </a>
        <a href="/docs#/Calendar" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">📅</div>
          <div><h4 class="font-semibold text-white">Calendar</h4><span class="text-sm text-slate-400">9 个端点</span></div>
        </a>
        <a href="/docs#/Notifications" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">🔔</div>
          <div><h4 class="font-semibold text-white">Notifications</h4><span class="text-sm text-slate-400">5 个端点</span></div>
        </a>
        <a href="/docs#/Messages" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">💬</div>
          <div><h4 class="font-semibold text-white">Messages</h4><span class="text-sm text-slate-400">5 个端点</span></div>
        </a>
        <a href="/docs#/Dashboard" class="p-5 bg-slate-800/50 border border-slate-700/50 rounded-xl flex items-center gap-4 transition-all module-hover no-underline">
          <div class="text-2xl">📊</div>
          <div><h4 class="font-semibold text-white">Dashboard</h4><span class="text-sm text-slate-400">9 个端点</span></div>
        </a>
      </div>
    </div>
  </section>

  <!-- CTA Section -->
  <section class="py-24">
    <div class="max-w-7xl mx-auto px-6">
      <div class="relative p-16 bg-gradient-brand rounded-3xl overflow-hidden">
        <div class="relative text-center">
          <h2 class="text-4xl font-bold mb-4">开始使用 HaloLight API</h2>
          <p class="text-lg opacity-90 mb-8">查看完整文档，快速集成到你的项目中</p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <a href="/docs" class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-white text-cyan-600 font-semibold rounded-xl hover:shadow-xl transition-all">
              📖 Swagger 文档
            </a>
            <a href="https://halolight.docs.h7ml.cn/guide/api-go" class="inline-flex items-center justify-center gap-2 px-8 py-4 bg-white/20 text-white font-semibold rounded-xl border border-white/40 hover:bg-white/30 transition-all" target="_blank">
              📚 完整使用指南
            </a>
          </div>
        </div>
      </div>
    </div>
  </section>

  <!-- Footer -->
  <footer class="py-12 border-t border-slate-700/50">
    <div class="max-w-7xl mx-auto px-6 text-center">
      <div class="flex flex-wrap justify-center gap-8 mb-6">
        <a href="/docs" class="text-slate-400 hover:text-white text-sm transition-colors">API 文档</a>
        <a href="https://halolight.docs.h7ml.cn/guide/api-go" target="_blank" class="text-slate-400 hover:text-white text-sm transition-colors">在线使用指南</a>
        <a href="https://github.com/halolight/halolight-api-go" target="_blank" class="text-slate-400 hover:text-white text-sm transition-colors">GitHub</a>
        <a href="https://github.com/halolight/halolight-api-go/issues" target="_blank" class="text-slate-400 hover:text-white text-sm transition-colors">问题反馈</a>
      </div>
      <p class="text-slate-400 text-sm">
        Built with ❤️ by <a href="https://github.com/h7ml" target="_blank" class="text-cyan-400 hover:underline">h7ml</a> |
        Powered by Gin 1.10 & GORM 2
      </p>
      <p class="text-slate-500 text-sm mt-2">
        Version 1.0.0 | Environment: ` + h.cfg.AppEnv + `
      </p>
    </div>
  </footer>
</body>
</html>`
}
