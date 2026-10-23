# Fibera

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white) ![Fiber](https://img.shields.io/badge/Fiber-v3-00ACD7?logo=gofiber&logoColor=white) ![GORM](https://img.shields.io/badge/GORM-MySQL-2D6E9E?logo=mysql&logoColor=white) ![Redis](https://img.shields.io/badge/Redis-optional-DC382D?logo=redis&logoColor=white) ![Prometheus](https://img.shields.io/badge/Prometheus-metrics-E6522C?logo=prometheus&logoColor=white) ![Swagger](https://img.shields.io/badge/Swagger-2.0-85EA2D?logo=swagger&logoColor=black) ![Scalar](https://img.shields.io/badge/Scalar-API%20Reference-1A1A1A) ![License](https://img.shields.io/badge/License-MIT-yellow)

> 基于 Go 与 Fiber v3 的后端服务脚手架：分层清晰、开箱可跑，适合作为 API 服务或后台管理系统的起步模板。

Fibera 采用分层架构，依赖关系由外向内单向收敛：`cmd` 作为进程入口，仅触发启动流程；`internal/app` 承担装配职责，串联基础设施初始化、中间件与路由注册；`internal/biz` 按领域划分模块，各模块内聚 `handler`、`service`、`dto` 与路由注册；`internal/infrastructure` 封装数据库、日志、缓存等外部依赖，业务层不直接依赖其实现。

项目内置配置加载与内嵌、结构化日志、GORM 自动同步表结构、JWT 鉴权、Prometheus 监控、统一响应结构，以及由 swag 注解生成、Scalar 渲染的交互式 API 文档，并附带可直接运行的登录认证与用户管理示例，便于快速起项目或作为团队工程规范参考。

## 特性

- **配置内嵌**：多环境 YAML 与 `.env` 通过 `go:embed` 打进二进制，单文件即可分发部署。
- **分层清晰**：`biz` 按领域拆分，`handler` / `service` / `dto` / `router` 四件套职责明确，业务与框架解耦。
- **可观测**：zap 结构化日志 + Prometheus 指标（QPS、时延、并发数），`/metrics` 开箱可用。
- **鉴权内置**：JWT 中间件带匿名白名单，可按需启用 Redis 会话缓存。
- **文档自动化**：swag 注解产出 Swagger 规范，Scalar 渲染交互式文档，改注释即改文档。
- **统一响应**：`{code, message, data}` 约定，业务侧只关心数据本身。
- **易扩展**：新增数据表无需改动同步逻辑，新增模块只需一次注册。

## 技术栈

| 分类 | 选型 | 说明 |
| --- | --- | --- |
| 语言 | Go 1.26+ | 泛型、`min`/`max` 等现代特性可用 |
| Web 框架 | Fiber v3 | 基于 fasthttp 的高性能框架 |
| ORM | GORM + MySQL | 启动时可选自动同步表结构 |
| 缓存 | Redis（可选） | go-redis v9，启用后 token 写入 Redis |
| 配置 | Viper + godotenv | 多环境 YAML 内嵌至二进制 |
| 日志 | zap + lumberjack | JSON / console 两种格式，按大小切割归档 |
| 鉴权 | golang-jwt v5 | 请求头 `Authorization: Bearer <token>` |
| 监控 | Prometheus client_golang | `/metrics` 暴露 QPS、时延、并发数 |
| 中间件 | recover、logger、metrics、jwt | 执行顺序见[请求链路](#请求链路) |
| 文档 | swag v2 + gofiber-scalar | Swagger 注解 + Scalar UI |

## 快速开始

| 依赖 | 版本 | 是否必需 |
| --- | --- | --- |
| Go | 1.26+ | 必需 |
| MySQL | 5.7+ / 8.0 | 必需 |
| Redis | 6.0+ | 可选，`redis.enable` 打开时必需 |

1. 创建环境变量文件（`APP_ENV` 决定加载哪份配置）：

   ```bash
   cp .env.example .env
   ```

2. 修改 [configs/config_dev.yaml](configs/config_dev.yaml) 中的数据库连接信息，首次运行保持 `database.sync: true` 以便自动建表。

3. 启动服务：

   ```bash
   go run ./cmd
   ```

4. 验证：

   ```bash
   curl -s http://127.0.0.1:7000/metrics | head
   ```

   浏览器打开 <http://127.0.0.1:7000/docs> 可查看交互式 API 文档。

## 目录结构

```
.
├── cmd/main.go                 # 入口，仅调用 app.Start()，附 swag 通用注解
├── configs/                    # 多环境配置（编译期内嵌）
├── docs/                       # swag 生成产物，勿手改
├── internal/
│   ├── apidoc/                 # Scalar 文档路由注册
│   ├── app/                    # 启动装配：基础设施、中间件、路由
│   ├── biz/                    # 业务模块，按领域拆分
│   │   ├── auth/               # 登录认证：handler / service / dto / router
│   │   └── user/               # 用户管理：handler / service / dto / router
│   ├── common/                 # 统一响应结构
│   ├── config/                 # 配置定义、加载与缺省值
│   ├── infrastructure/         # 外部依赖：db / log / redis
│   ├── middleware/             # JWT、请求日志、Prometheus 指标
│   └── model/                  # GORM 数据模型与自注册
├── utils/security/             # 密码哈希、加密、敏感信息脱敏
└── embed.go                    # 配置文件与 .env 的内嵌入口
```

业务模块内部约定：

```
internal/biz/<module>/
├── dto.go        # 请求/响应结构与模型转换
├── handler.go    # 参数解析、调用 service、输出统一响应（含 swag 注解）
├── service.go    # 业务逻辑，不感知 HTTP
└── router.go     # 路由注册入口 RegisterRoutes
```

## 请求链路

```
请求 → recover → logger → metrics → jwt → 路由 → handler → service → 基础设施
```

- `recover`：panic 兜底，转为 500 响应而非中断连接。
- `logger`：zap 结构化记录方法、路径、状态码、耗时与客户端 IP。
- `metrics`：采集 QPS、时延与并发数，标签使用路由模板，避免基数膨胀。
- `jwt`：校验 Bearer Token；`/api/auth/login`、`/metrics` 与文档路径免鉴权。

## 配置说明

配置加载优先级：

1. 二进制同级目录下的 `configs/config.yaml`（外部配置，便于容器挂载覆盖）。
2. 内置配置 `configs/config_{APP_ENV}.yaml`，`APP_ENV` 取自 `.env`，缺省 `dev`。

> `.env` 与 YAML 均为编译期内嵌，修改后需重新编译（`go run` 会自动重新编译）。

| 配置项 | 说明 |
| --- | --- |
| `app.port` | 服务监听端口 |
| `app.version` | 版本号，启动日志与 API 文档中展示 |
| `database.driver` | 数据库类型，当前固定 `mysql` |
| `database.source` | DSN 连接串 |
| `database.sync` | 启动时是否同步表结构 |
| `database.max_idle_conns` / `max_open_conns` | 连接池空闲/最大连接数 |
| `redis.enable` | 是否启用 Redis（token 缓存随之生效） |
| `redis.addr` / `password` / `db` | Redis 地址、密码与库号 |
| `log.level` | 日志级别：`debug`、`info`、`warn`、`error` |
| `log.format` | 日志格式：`text`、`json` |
| `log.output_method` | 输出方式：`console`、`file`、`default`（同时输出） |
| `log.output_dir` | 日志目录，落在进程工作目录的上一级 |
| `log.max_size` / `max_backups` / `max_age` | 日志切割大小(MB)、备份数、保留天数 |
| `jwt.issuer` | 签发者 |
| `jwt.signing_method` | 签名算法：`HS256`、`HS384`、`HS512` |
| `jwt.expiration_time` | Token 有效期，如 `30m`、`24h` |
| `docs.enable` | 是否启用 API 文档 |
| `docs.path` / `base_path` | 文档访问路径、反向代理前缀 |
| `docs.title` / `theme` | 文档标题与 Scalar 主题 |
| `docs.cache_age` / `fallback_cache_age` | 页面与内置脚本缓存时间（秒） |
| `docs.force_offline` | 是否强制加载内置脚本，内网离线部署保持 `true` |
| `docs.proxy_url` / `custom_style` | 跨域代理地址、自定义 CSS 变量 |

## API 文档

接口定义由 swag 注解生成，页面由 Scalar 渲染，默认访问 `http://127.0.0.1:7000/docs`，规范文件为 `/docs/doc.json`。文档路径默认匿名可访问。

修改 handler 注解后重新生成：

```bash
go install github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc5
swag init -g cmd/main.go -o docs
```

注解约定：

- 通用信息（标题、版本、鉴权定义）写在 [cmd/main.go](cmd/main.go)，标题与版本在运行时以 `docs.title`、`app.version` 覆盖。
- 接口注解写在 handler 方法上方，响应统一使用 `common.Response`，通过 `{data=...}` 覆盖数据字段，例如 `@Success 200 {object} common.Response{data=[]user.UserResponse}`。
- 需鉴权的接口追加 `@Security BearerAuth`。

## 接口说明

除 `/api/auth/login`、`/metrics` 与文档路径外，所有接口均需携带请求头 `Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/api/auth/login` | 登录认证，返回 token 与用户信息 |
| GET | `/api/users` | 查询用户列表 |
| GET | `/api/users/:id` | 查询用户详情 |
| POST | `/api/users` | 创建用户 |
| PUT | `/api/users/:id` | 更新用户 |
| DELETE | `/api/users/:id` | 删除用户 |

调用示例：

```bash
# 登录获取 token（用户名与用户 ID 二选一，同时提供时以用户 ID 为准）
TOKEN=$(curl -s -X POST http://127.0.0.1:7000/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"userName":"<用户名>","password":"<密码>"}' | jq -r '.data.token')

# 携带 token 访问受保护接口
curl -s http://127.0.0.1:7000/api/users -H "Authorization: Bearer $TOKEN"
```

统一响应结构：

```json
{ "code": 200, "message": "success", "data": {} }
```

失败时 `data` 为 `null`，`code` 与 HTTP 状态码一致：

```json
{ "code": 404, "message": "用户不存在", "data": null }
```

## 监控

`/metrics` 以 Prometheus 文本格式暴露以下指标：

| 指标 | 类型 | 标签 | 说明 |
| --- | --- | --- | --- |
| `http_requests_total` | Counter | `method`、`path`、`status` | 请求总数 |
| `http_request_duration_seconds` | Histogram | `method`、`path`、`status` | 请求耗时分布 |
| `http_requests_in_flight` | Gauge | `method` | 当前并发处理数 |

`path` 取路由模板（如 `/api/users/:id`）而非真实路径，避免标签基数爆炸；`/metrics` 自身不参与采集。

## 扩展约定

### 新增数据表

在 `internal/model` 下创建模型文件，实现 `TableName()` 并在 `init` 中注册：

```go
func init() {
	db.Register(new(SysUser))
}
```

启动同步时会遍历全部已注册模型，新增表无需修改同步代码。

### 新增业务模块

在 `internal/biz` 下创建模块目录（`dto.go`、`service.go`、`handler.go`、`router.go`），提供统一入口：

```go
func RegisterRoutes(app *fiber.App) {
	// 构造 service 与 handler，并注册路由
}
```

随后在 [internal/app/bootstrap.go](internal/app/bootstrap.go) 的 `registerRoutes` 中调用即可。

### 新增接口

1. 在 `handler` 方法上方补齐 swag 注解（`@Summary`、`@Param`、`@Success`、`@Router` 等）。
2. 需要鉴权时追加 `@Security BearerAuth`。
3. 执行 `swag init -g cmd/main.go -o docs` 更新文档。

## 参与贡献

1. Fork 并新建分支：`git checkout -b feat/your-feature`。
2. 提交前执行 `gofmt -l .`、`go vet ./...`、`go build ./...`，保持零告警。
3. 发起 Pull Request，说明变更动机与验证方式。

## 许可证

[MIT License](LICENSE) © 2026 Fibera
