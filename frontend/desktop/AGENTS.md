# frontend/desktop Agent 指南

## 模块定位

`frontend/desktop` 是 Wails 桌面端前端入口预留。桌面端复用 Web 控制台设计语言，但优先呈现本机运行状态、插件目录、SQLite 状态、采集调试和本地 APIKey 授权。

## 开发规则

- 先阅读 `docs/current/04-前端设计规范.md`。
- 与 `frontend/web` 共用可复用组件和视觉规范，避免另起一套风格。
- 桌面端页面必须清晰区分本机资源和服务器资源。
- 本地 HTTP API 默认可开启，但必须使用 APIKey 授权。
- 插件目录、配置文件路径、日志路径等本机路径展示时避免泄露敏感内容。

## 验证

- 前端构建规则与 `frontend/web/AGENTS.md` 保持一致。
- 涉及 Wails 桌面能力时同步验证桌面启动、窗口尺寸和系统主题适配。
