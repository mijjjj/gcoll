# docs Agent 指南

## 适用范围

`docs/` 保存 gcoll 的当前设计和历史 ADR。开始修改文档前先阅读 `docs/README.md`。

## 文档权威性

- 当前有效设计只写入 `docs/current/`。
- `docs/adr/` 只记录已确认且难以逆转的历史决策背景，不直接覆盖当前规范。
- AI 执行流程维护在 `.codex/skills/gcoll-development/SKILL.md`。
- 前端视觉执行细节维护在 `.codex/skills/gcoll-front-module/references/design-guidelines.md`。

## 写作规则

- 文档正文、说明文本、代码注释示例必须使用中文。
- 命令、路径、配置键、协议字段、API 路径、JSON/YAML 字段名保持英文或约定格式。
- 新增设计必须先更新 `docs/current/` 中对应规范。
- 不为同一主题新增多个平级草稿文档；需要保留取舍背景时新增 ADR。
- 涉及 HTTP API、插件协议、SDK、事件格式或数据库迁移的破坏性变更，必须显式说明版本策略或兼容选择。
- 不再新增分散的 AI 入口或界面设计草稿目录；对应内容应整理到项目 Agent、skill 或 `docs/current/`。

## 验证

- 修改 Mermaid 图后检查语法和节点可读性。
- 修改前端设计规范后，确认 `frontend/web/AGENTS.md` 与相关 skill 没有冲突。
- 修改插件协议后，同步检查 `plugins/AGENTS.md` 与 `.codex/skills/gcoll-plugin-workflow`。
