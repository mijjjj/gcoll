---
name: gcoll-front-module
description: gcoll Web 控制台前端模块开发流程。用于修改 `frontend/web` 中的 Vue 3 页面、TypeScript API 模块、Pinia 状态、Vue Router 路由、Naive UI 表单表格抽屉、插件自定义配置页面、i18n 文本、空状态、错误状态或权限状态。
---

# gcoll Front Module

使用本技能前先阅读 `frontend/web/AGENTS.md`、`docs/current/04-前端设计规范.md` 和 `references/design-guidelines.md`。

## 工作流

1. 搜索同域页面、组件、store、API 模块和样式模式。
2. API 类型与方法对齐后端 `api/<domain>/v1` 和实际 HTTP 路径。
3. 页面使用 `<script setup lang="ts">`、Composition API、Naive UI 和现有布局。
4. 用户可见文本进入现有 i18n 结构；中英文 key 结构保持一致。
5. 完成后检查浅色、黑夜、空状态、错误状态、只读权限状态和 `1280x800`。
6. 从仓库根目录运行 `pnpm --dir frontend/web build`。

## 强制规则

- 不硬编码用户可见文本，除非当前项目对应模块尚未引入 i18n 且本次任务明确不触及国际化结构。
- 不在页面 catch 中重复弹出已由拦截器统一处理的 API 错误。
- 成功反馈由业务代码显式展示。
- 表格、分栏、抽屉和表单优先，避免装饰性卡片墙。
- 图标按钮使用 lucide 或已有图标库，并补充 tooltip 或可访问文本。
- 响应数据使用可选链和默认值防御 `null`、`undefined`。
- 列表 key 使用稳定业务 ID。
- 服务端和桌面端共享“简洁边缘桌面工作台”设计语言；桌面端更强调本机状态、连接测试、点位表和调试日志。
- 插件自定义页面必须保留宿主标题、保存状态、权限提示、审计提示和主题同步。
- 状态颜色不得作为唯一信息来源，必须同时提供文字或标签。

## 按需读取

- 创建管理列表或 CRUD 页面时读取 `references/module-page-pattern.md`。
- 修改菜单、路由标题、权限状态或 i18n 时读取 `references/permissions-i18n.md`。
- 设计新页面、调整布局密度、桌面端页面或复杂状态时读取 `references/design-guidelines.md`。

## 验证

```powershell
pnpm --dir frontend/web build
```
