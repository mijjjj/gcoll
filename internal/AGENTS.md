# internal Agent 指南

## 模块定位

`internal/` 是 gcoll 的 GoFrame 后端核心运行时，包含启动装配、控制器、中间件、领域服务、数据面队列、插件宿主和基础设施模块。

修改本目录 Go 代码时使用 `goframe-v2` 技能，并参考 `.codex/skills/gcoll-backend`。涉及数据库表或字段变更时同时使用 `.codex/skills/gcoll-db-migrate-dao`。涉及枚举校验时同时使用 `.codex/skills/goframe-enum-refresh`。

## 分层规则

- `internal/cmd` 只负责 Server 装配、配置加载、中间件挂载和路由注册。
- `internal/controller/<domain>` 只做 HTTP 入参、上下文和响应编排，不承载复杂业务规则。
- `internal/service/<domain>` 承载领域服务和模块编排；MVP 阶段不额外引入 `logic/` 兼容层。
- `internal/service/<domain>/` 下默认使用与模块同名的 `<domain>.go` 作为主实现文件；不要保留仅占位的 `doc.go`，也不要使用泛化的 `service.go` 作为模块主文件名。
- `internal/service/middleware` 统一处理响应包裹、认证、错误映射和跨域等横切逻辑。
- `internal/dao`、`internal/model/do`、`internal/model/entity` 引入数据库表后由 GoFrame CLI 生成，不手工维护生成代码。

## 模块边界

- 控制面模块：`auth`、`device`、`point`、`pluginconfig`、`pluginmgmt`、`marketplace`、`scheduler`。
- 数据面模块：`pluginhost`、`collector`、`runtimequeue`、`pointcache`、`pipeline`、`delivery`。
- 基础设施模块：`storage`、`secret`、`observability`、`eventbus`、`i18n`。
- 控制面不得进入高频采集路径；数据面不得依赖 HTTP 请求生命周期。

## 代码规则

- 先搜索现有 service、controller、middleware 和测试模式，再新增实现。
- 结构体、接口、函数、字段注释使用中文。
- 错误应保留上下文；日志字段、配置键、协议字段保持英文。
- 高频路径不得逐条同步写库，不无限制创建 goroutine。
- 缺失必需配置、密钥、插件或设备资源时直接返回明确错误，不隐式选择默认资源。

## 验证

在仓库根目录执行：

```powershell
go test ./...
```
