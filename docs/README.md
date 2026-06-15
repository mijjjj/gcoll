# gcoll 文档入口

本目录是后续 AI 和开发者进入项目的唯一文档入口。当前有效设计集中在 `docs/current/`，历史决策保留在 `docs/adr/`。

## 阅读顺序

1. `docs/current/01-产品范围.md`：确认产品定位、MVP 边界和里程碑。
2. `docs/current/02-系统架构.md`：确认运行时、模块边界、数据流和部署策略。
3. `docs/current/03-插件与点位协议.md`：确认插件模型、配置模型、点位表和协议边界。
4. `docs/current/04-前端设计规范.md`：确认 Vue 前端和 Wails 桌面端的界面规范。
5. `docs/current/05-开发任务规范.md`：确认后续 AI 开发时的任务拆分、验收和测试要求。
6. `docs/ai/AI开发入口.md`：给 AI 执行具体开发任务时使用。

## 文档权威性

- `docs/current/` 是当前实现依据。
- `docs/adr/` 是决策记录，用于解释为什么这么做，不直接覆盖当前规范。
- 新增设计必须先更新 `docs/current/` 中对应规范；只有出现不可逆或对外契约相关的取舍时，才新增 ADR。
- 不再维护分散的旧版产品、架构、设计草稿；避免同一规则在多处重复出现。

## 当前技术基线

- 后端：Go + GoFrame。
- 桌面端：Wails 3。
- 前端：Vue 3 + TypeScript + Vite + Naive UI。
- 状态管理：Pinia。
- 路由：Vue Router。
- 图表：ECharts 或 uPlot。
- 插件通信：gRPC + Protocol Buffers。
- 桌面默认数据库：SQLite。
- 服务器生产数据库：PostgreSQL。

## 最新优先原则

MVP 阶段内部契约按最新设计推进，不为尚未发布的旧内部设计添加兼容分支。插件协议、HTTP API、SDK、事件格式一旦面向第三方公开，就必须显式版本化；破坏性变更实施前必须确认是否需要兼容策略。
