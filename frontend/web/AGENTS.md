# frontend/web Agent 指南

## 模块定位

`frontend/web` 是 gcoll Web 控制台，基于 Vue 3、TypeScript、Vite、Pinia、Vue Router、Naive UI、ECharts/uPlot 预留和 lucide 图标。

开发页面、路由、状态管理、API 调用或交互时参考 `.codex/skills/gcoll-front-module`，并先阅读 `docs/current/04-前端设计规范.md`。

## 核心规则

- 服务端控制台和桌面端共享设计语言，桌面端更强调本机状态和调试。
- 表格、分栏、抽屉、表单和密集信息布局优先，避免卡片墙和营销式英雄区。
- 插件自定义页面外层必须保留宿主工具栏、保存状态、权限提示和审计提示。
- 用户可见文本应进入现有 i18n 结构；新增 key 时保持中英文结构一致。
- API 模块与后端 `api/<domain>/v1` 契约保持一致，响应数据用可选链和默认值防御空值。
- 图标按钮使用 lucide 或已有图标库，并提供 tooltip 或可访问文本。
- 页面实现后检查浅色模式、黑夜模式、空状态、错误状态、只读权限状态和 `1280x800` 视口。

## 代码组织

- 页面优先使用 `<script setup lang="ts">` 和 Composition API。
- 页面组件使用 PascalCase；API、store、composable 文件使用 camelCase。
- 导入顺序：Vue/Router/Pinia/i18n，Naive UI，API/类型，组件，composables，utils。
- 响应式代码按组合式 API、状态、计算属性、方法、生命周期分段组织。
- 列表渲染使用稳定业务 ID 作为 key，不使用数组下标。

## 常用命令

在仓库根目录执行：

```powershell
pnpm --dir frontend/web build
pnpm --dir frontend/web dev
pnpm --dir frontend/web preview
```
