# gcoll Agent 指南

## 项目定位

gcoll 是一套面向设备数据采集、点位管理、插件编排和北向转发的通用采集平台。项目同时覆盖服务器端、Wails 桌面端、Vue 控制台、插件协议、插件 SDK 和内置插件示例。

## 模块入口

- `api/`：HTTP/OpenAPI 与 gRPC/Protocol Buffers 契约。修改公开契约时必须同步 `docs/current/03-插件与点位协议.md` 或相关当前规范。
- `internal/`：GoFrame 后端核心运行时。修改 Go 代码时使用 `goframe-v2` 和 `.codex/skills/gcoll-backend`。
- `frontend/web/`：Vue 3 + TypeScript + Vite + Naive UI Web 控制台。开发页面时使用 `.codex/skills/gcoll-front-module`。
- `frontend/desktop/`：Wails 桌面端前端入口预留，遵循 `docs/current/04-前端设计规范.md`。
- `plugins/`：插件 SDK、内置插件和示例。修改插件协议、宿主交互或示例时使用 `.codex/skills/gcoll-plugin-workflow`。
- `manifest/`：GoFrame 配置、国际化和后续数据库迁移入口。
- `docs/`：当前有效设计与历史 ADR。遵循 `docs/AGENTS.md`。

## 语言要求

本项目生成的任何文档和代码注释都必须使用中文。

适用范围包括但不限于：

- 产品需求文档。
- 架构设计文档。
- 插件协议文档。
- API 说明文档。
- 部署文档。
- 测试文档。
- README、变更日志、ADR。
- 代码中的行注释、块注释、接口注释、结构体字段注释、函数注释。

例外情况：

- 编程语言关键字、标准库名称、第三方库名称、命令、文件路径、协议字段、配置键、API 路径、错误码、日志字段等应保持原始英文或约定格式。
- 面向外部生态的机器可读字段，例如 JSON、YAML、OpenAPI、Protocol Buffers、数据库字段名，可以使用英文命名，但说明文本必须使用中文。
- 引用第三方官方名称或原文术语时，可以保留英文，并在首次出现时补充中文解释。

## 工程原则

- 优先采用当前主流、清晰、可维护的实现方式。
- MVP 阶段内部契约以最新设计为准，不为尚不存在的旧版本增加兼容分支。
- 一旦插件协议、HTTP API、SDK、事件格式对第三方开发者公开，就必须进行显式版本管理。
- 涉及外部公开契约的破坏性变更，实施前必须明确是否需要兼容策略。

## AI 开发入口

- 开始任何需求、设计或实现前，先阅读 `docs/README.md`。
- 当前有效设计以 `docs/current/` 为准；`docs/adr/` 只记录历史决策背景。
- 需要让 AI 按项目方式持续开发时，优先使用 `.codex/skills/gcoll-development`。
- 本地验证方式遵循 `.codex/skills/gcoll-development`：默认只运行后端测试和前端构建，不启动前后端服务。
- 如果旧文档与 `docs/current/` 冲突，以 `docs/current/` 为准。

## 跨模块硬规则

- 控制面不得进入高频采集数据路径；数据面不得依赖 HTTP 请求生命周期。
- 采集明细不落库，最新点位值保存在运行时缓存。
- 插件不得直接访问宿主数据库；插件配置由宿主保存，敏感值通过宿主密钥存储引用。
- 南向插件只提交采集记录，不直接转发目标系统；北向插件只接收过滤后的记录，不主动采集设备。
- 设备配置、点位表、插件配置结构必须分离，并在需要历史追溯时显式版本化。
- 修改 HTTP API、插件协议、SDK、事件格式或数据库迁移前，先判断是否属于公开契约破坏性变更。
- 生成或修改代码后运行对应范围验证：后端优先 `go test ./...`，前端优先 `pnpm --dir frontend/web build`。
- 不要为项目任何目录或文件新增测试文件；禁止生成 `*_test.go`、前端测试规格文件或其他测试文件。
- 具体南向协议的连接参数、点位 metadata、读取计划、写入限制和调试摘要必须留在插件实现内；主程序不得引入协议专用表、协议专用 service 分支或协议专用前端页面。
- 数据库迁移若 SQLite 与 PostgreSQL 写法不同，必须分别维护在 `manifest/migrations/sqlite/` 与 `manifest/migrations/pgsql/`，同一版本文件名保持一致。
- 数据库表和字段必须添加中文注释；PostgreSQL 使用 `COMMENT ON TABLE/COLUMN`，SQLite 在迁移 SQL 中使用 `-- 表注释：` 和 `-- 字段注释：` 注释块记录。
- GoFrame 数据访问必须以 `internal/dao` 作为模型入口，以 `internal/model/do` 作为写入与更新结构，以 `internal/model/entity` 作为读取结果结构。
- 写入、更新和条件拼装优先使用 GoFrame DAO/Model 能力与 `do` 结构，避免手写 `map`，以便自动忽略空值字段。
- 事务必须通过 GoFrame 的 `Transaction` 包裹；事务回调内统一通过回调传入的 `ctx` 透传事务，不直接调用 `tx.Insert`、`tx.Update`、`tx.Delete`、`tx.Exec`，而是继续使用 `dao.*.Ctx(ctx)` 或 `db.Ctx(ctx)` 执行操作。
- 执行 `Delete`、`Update` 或其他依赖业务键的数据库操作前，必须先在业务层校验关键条件不为空；不要把“条件为空”的判断留给 ORM 或数据库层兜底。
- 标准库与第三方库存在等效能力时，优先使用 GoFrame 已提供并适合当前场景的工具库、组件和辅助方法。

## 项目技能

- `.codex/skills/gcoll-development`：进行 gcoll 需求、设计、实现和交付时的总入口。
- `.codex/skills/gcoll-backend`：修改 GoFrame 后端 API、Controller、Service、运行时装配和中间件时使用。
- `.codex/skills/gcoll-front-module`：开发 Web 控制台页面、路由、状态、API 调用和 Naive UI 交互时使用。
- `.codex/skills/gcoll-plugin-workflow`：修改插件协议、插件宿主、SDK、内置插件或示例插件时使用。
- `.codex/skills/gcoll-db-migrate-dao`：引入或调整数据库表、字段、迁移和 GoFrame DAO 生成时使用。
- `.codex/skills/goframe-enum-refresh`：新增或调整 GoFrame 枚举常量和 `v:"enums"` 校验时使用。
- `.codex/skills/no-default-fallback`：修改运行时配置、设备配置、插件配置、密钥、外部端点或默认/兜底逻辑时使用。
