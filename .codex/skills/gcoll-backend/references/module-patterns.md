# 后端模块模式

## 生成和实现顺序

1. 数据库表结构变更先使用 `gcoll-db-migrate-dao`。
2. API 请求和响应结构放在 `api/<domain>/v1/`。
3. Controller 放在 `internal/controller/<domain>/`，只处理 HTTP 入口适配。
4. Service 放在 `internal/service/<domain>/`，承载领域规则和模块编排。
5. 路由和中间件挂载放在 `internal/cmd` 的既有 Server 组装代码。
6. 枚举变更使用 `goframe-enum-refresh`。

## API 层

- 使用 `g.Meta` 标注 `path`、`method`、`summary`、`tags`。
- 路径参数使用 GoFrame 约定标签。
- 校验标签保持清晰，错误提示使用中文。
- 请求和响应结构与数据库、前端和文档同步。

## Controller 层

- Controller 只做参数转换、上下文透传和 service 调用。
- 不在 Controller 中写复杂业务规则、数据库访问或数据面循环。
- 保持返回结构与 API 响应定义一致。

## Service 层

- 优先复用现有 service 和模块边界。
- 控制面 service 处理配置、设备、点位、插件管理、任务和审计。
- 数据面 service 处理插件宿主、采集入口、队列、点位缓存、规则过滤和北向转发。
- 跨模块调用需要保持依赖方向清晰，避免控制面和数据面互相缠绕。

## 错误和日志

- 错误消息面向用户时使用中文。
- 日志字段使用英文 key，避免记录敏感值。
- 底层错误应保留上下文，上层业务错误应明确可操作原因。

## 路由

- 管理 API 挂在 `/api/v1` 下。
- 优先使用 GoFrame 标准路由和 `group.Bind(...)`。
- 中间件统一处理响应、认证、跨域和错误映射。
