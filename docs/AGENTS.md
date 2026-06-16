# docs Agent 指南

## 适用范围

`docs/` 保存 gcoll 的当前设计、历史 ADR、AI 开发入口和设计资料。开始修改文档前先阅读 `docs/README.md`。

## 文档权威性

- 当前有效设计只写入 `docs/current/`。
- `docs/adr/` 只记录已确认且难以逆转的历史决策背景，不直接覆盖当前规范。
- `docs/design/` 保存界面设计资料和视觉参考；前端实现以 `docs/current/04-前端设计规范.md` 为准。
- `docs/ai/AI开发入口.md` 是后续 AI 开发任务的入口说明。

## 写作规则

- 文档正文、说明文本、代码注释示例必须使用中文。
- 命令、路径、配置键、协议字段、API 路径、JSON/YAML 字段名保持英文或约定格式。
- 新增设计必须先更新 `docs/current/` 中对应规范。
- 不为同一主题新增多个平级草稿文档；需要保留取舍背景时新增 ADR。
- 涉及 HTTP API、插件协议、SDK、事件格式或数据库迁移的破坏性变更，必须显式说明版本策略或兼容选择。

## 验证

- 修改 Mermaid 图后检查语法和节点可读性。
- 修改前端设计规范后，确认 `frontend/web/AGENTS.md` 与相关 skill 没有冲突。
- 修改插件协议后，同步检查 `plugins/AGENTS.md` 与 `.codex/skills/gcoll-plugin-workflow`。
