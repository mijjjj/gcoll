# manifest Agent 指南

## 模块定位

`manifest/` 保存 gcoll 运行配置、国际化资源和后续数据库迁移入口。

## 配置规则

- 配置键保持英文，说明文本使用中文。
- 服务器默认配置位于 `manifest/config/config.yaml`。
- 桌面默认数据库使用 SQLite；服务器生产数据库使用 PostgreSQL。
- 敏感值不得写入普通配置文件，必须通过宿主密钥存储引用。
- 缺失必需配置时直接报错，不隐式选择默认业务资源。

## 国际化规则

- `manifest/i18n/` 保存后端或运行时国际化资源。
- 新增用户可见错误消息时同步中英文文件。
- 错误码、日志字段、配置键保持英文。

## 数据库规则

- 引入数据库迁移后，迁移文件应放在 `manifest/migrations/`。
- SQLite 与 PostgreSQL SQL 写法不一致时，迁移应拆分到 `manifest/migrations/sqlite/` 和 `manifest/migrations/pgsql/`，同一变更使用相同版本号文件名。
- 表结构变更使用 `.codex/skills/gcoll-db-migrate-dao`。
- 数据库迁移必须为每张表和每个字段写中文注释；PostgreSQL 使用 `COMMENT ON TABLE/COLUMN`，SQLite 使用 SQL 注释块记录表和字段说明。
- 采集明细不落库；数据库只保存配置、任务、插件、用户、审计和必要运行状态。
