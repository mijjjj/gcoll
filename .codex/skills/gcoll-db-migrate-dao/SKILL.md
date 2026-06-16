---
name: gcoll-db-migrate-dao
description: gcoll 数据库迁移与 GoFrame DAO 生成流程。用于新增或调整 `manifest/migrations` 迁移、数据库表字段、SQLite/PostgreSQL 存储结构、`internal/dao`、`internal/model/do`、`internal/model/entity` 或任何需要执行 `gf gen dao` 的后端变更。
---

# gcoll 数据库迁移与 DAO 流程

当表结构或字段变更时按以下顺序执行：

1. 先阅读 `docs/current/02-系统架构.md` 和 `docs/current/05-开发任务规范.md`。
2. 在 `manifest/migrations/` 新增成对迁移 SQL：`*.up.sql` 与 `*.down.sql`。
3. 迁移必须兼顾桌面 SQLite 和服务器 PostgreSQL 的当前策略；如果无法兼容，先在当前规范中说明数据库边界。
4. 使用 `manifest/config/config.yaml` 中的数据库连接在仓库根目录执行迁移。
5. 继续在仓库根目录执行 `gf gen dao`，生成 `internal/dao`、`internal/model/do`、`internal/model/entity`。
6. 后续代码优先使用生成的 `do` 和 `entity`。
7. 若生成结构无法满足业务场景，再在 `internal/model` 增加业务结构体。

## 强制要求

- 数据库迁移必须显式版本化。
- 采集明细不落库。
- 敏感配置不进入普通 JSON 字段明文。
- 设备插件配置、点位表、插件配置结构需要历史追溯时必须支持版本化。
- 不手工修改生成的 DAO、DO、Entity 文件。
- 运行迁移和 `gf gen dao` 前确认当前工作目录是仓库根目录。

## 验证

```powershell
gf gen dao
go test ./...
```

如果本地缺少迁移工具、数据库或 `gf` 命令，最终说明中明确未执行的验证和原因。
