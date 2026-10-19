# Fibera

基于 Go 与 Fiber v3 的后端服务骨架，提供配置加载、日志、数据库自动同步、JWT 鉴权及用户模块 CRUD 示例。

## 技术栈

- Web 框架：Fiber v3
- ORM：GORM（MySQL）
- 缓存：Redis（可选）
- 配置：Viper + godotenv，配置内嵌于二进制
- 日志：zap + lumberjack
- 鉴权：JWT

## 目录结构

```
.
├── cmd/main.go                    # 程序入口
├── configs/                       # 环境配置文件
├── internal/
│   ├── app/                       # 应用启动与装配
│   ├── biz/user/                  # 用户业务模块（handler、service、router）
│   ├── common/                    # 通用响应结构
│   ├── config/                    # 配置加载与定义
│   ├── infrastructure/            # 数据库、日志、Redis
│   ├── middleware/                # JWT、请求日志中间件
│   └── model/                     # 数据表模型
├── utils/security/                # 通用工具
└── embed.go                       # 内嵌配置文件
```

## 快速开始

环境要求：Go 1.26+、MySQL，Redis 按需启用。

1. 创建环境变量文件：

```bash
cp .env.example .env
```

2. 按需修改 [configs/config_dev.yaml](configs/config_dev.yaml) 中的数据库连接信息。

3. 启动服务：

```bash
go run ./cmd
```

默认监听 `7000` 端口。首次运行建议保持 `database.sync: true`，用于自动建表。

## 配置说明

程序通过 `.env` 中的 `APP_ENV` 选择配置，默认 `dev`，对应 `configs/config_{env}.yaml`。若二进制文件同级目录存在 `configs/config.yaml`，将优先加载该外部配置。

| 配置项 | 说明 |
| --- | --- |
| `app.port` | 服务监听端口 |
| `database.sync` | 是否在启动时同步表结构 |
| `database.max_idle_conns` / `max_open_conns` | 连接池参数 |
| `redis.enable` | 是否启用 Redis 缓存 |
| `log.level` / `format` / `output_method` | 日志级别、格式与输出方式 |
| `jwt.expiration_time` | Token 有效期 |

## 接口说明

除 `/login`、`/public` 外，所有接口均需携带请求头 `Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/users` | 查询用户列表 |
| GET | `/api/users/:id` | 查询用户详情 |
| POST | `/api/users` | 创建用户 |
| PUT | `/api/users/:id` | 更新用户 |
| DELETE | `/api/users/:id` | 删除用户 |

统一响应结构：

```json
{ "code": 200, "message": "success", "data": {} }
```

失败时 `data` 为 `null`，`code` 与 HTTP 状态码一致：

```json
{ "code": 404, "message": "用户不存在", "data": null }
```

## 扩展约定

### 新增数据表

在 `internal/model` 下创建模型文件，实现 `TableName()` 并在 `init` 中注册：

```go
func init() {
	db.Register(new(SysUser))
}
```

启动时同步逻辑会遍历全部已注册模型，新增表无需修改同步代码。

### 新增业务模块

在 `internal/biz` 下创建模块目录，提供统一的注册入口：

```go
func RegisterRoutes(app *fiber.App) {
	// 构造 service 与 handler，并注册路由
}
```

随后在 [internal/app/bootstrap.go](internal/app/bootstrap.go) 的 `registerRoutes` 中调用即可。
