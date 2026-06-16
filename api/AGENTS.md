# api Agent 指南

## 适用范围

`api/` 保存 gcoll 的公开契约：

- `api/<domain>/v1/`：GoFrame HTTP/OpenAPI 请求与响应结构。
- `api/proto/`：插件 gRPC 与 Protocol Buffers 协议。
- `api/openapi/`：导出的 OpenAPI 资料或生成物。

## 契约规则

- HTTP API、插件协议、SDK 和事件格式一旦面向第三方公开，必须显式版本化。
- MVP 内部契约按最新设计推进；公开契约破坏性变更实施前必须明确兼容策略。
- 请求和响应结构体注释使用中文；JSON 字段名、协议字段名保持英文或约定格式。
- HTTP 管理 API 默认使用统一响应包裹：`code`、`message`、`data`。
- 插件协议与管理 API 分离演进，不混用 HTTP 生命周期处理高频采集数据。

## GoFrame API 约定

- API 定义放在 `api/<domain>/v1/`。
- 优先使用 `g.Meta` 描述路径、方法、摘要和标签。
- 路径分组遵循 `docs/current/02-系统架构.md` 中的 API 策略。
- API 结构变更必须同步 controller、service、前端调用和文档。

## 插件协议约定

- 南向插件只提交采集记录。
- 北向插件只接收过滤后的记录。
- 插件不得直接访问宿主数据库。
- 敏感配置只通过宿主密钥存储引用传递。

## 验证

- 修改 HTTP API 后运行 `go test ./...`。
- 修改 proto 后同步更新生成代码、插件 SDK、示例插件和 `docs/current/03-插件与点位协议.md`。
