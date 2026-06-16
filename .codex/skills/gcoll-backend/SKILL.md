---
name: gcoll-backend
description: gcoll GoFrame 后端开发流程。用于修改 `api/`、`internal/`、`manifest/` 中的 HTTP API、Controller、Service、运行时装配、中间件、配置加载、错误处理、枚举、数据库访问或模块 CRUD；涉及表结构时同时使用 `gcoll-db-migrate-dao`。
---

# gcoll Backend

使用本技能前先阅读 `internal/AGENTS.md` 和相关 `docs/current/` 规范。

## 工作流

1. 先搜索现有 API、controller、service、middleware、测试和文档模式。
2. 判断变更是否影响 HTTP API、插件协议、SDK、事件格式或数据库迁移。
3. API 结构放在 `api/<domain>/v1/`，Controller 放在 `internal/controller/<domain>/`。
4. 业务编排放在 `internal/service/<domain>/`，MVP 阶段不新增 `logic/` 兼容层。
5. 运行时装配、路由和中间件挂载放在 `internal/cmd` 或既有启动模块。
6. 修改后从仓库根目录运行 `go test ./...`。

## 强制规则

- 结构体、接口、函数、字段注释使用中文。
- 管理 API 默认通过中间件返回 `code`、`message`、`data`，控制器不重复封装。
- 控制面不得进入高频采集路径；数据面不得依赖 HTTP 请求生命周期。
- 高频采集路径不得逐条同步写库。
- 缺失必需配置、设备资源、插件资源、密钥或外部端点时直接返回错误，不隐式兜底。
- 敏感配置不得进入普通 JSON 或日志明文。
- GoFrame 生成代码只在引入对应生成体系后使用，不手工维护 `internal/dao`、`internal/model/do`、`internal/model/entity`。

## 按需读取

- 创建或调整后端模块时读取 `references/module-patterns.md`。
- 实现静态关联查询时读取 `references/relation-patterns.md`。

## 验证

```powershell
go test ./...
```

如果运行了 `gf gen dao`、`gf gen ctrl`、`gf gen service` 或 `gf gen enums`，在最终说明中列出生成命令和影响范围。
