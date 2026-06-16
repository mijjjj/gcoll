---
name: gcoll-plugin-workflow
description: gcoll 插件开发与协议维护流程。用于修改插件协议、`api/proto`、插件宿主、gRPC 握手、南向采集、北向转发、插件配置结构、插件 SDK、内置插件、示例插件或插件相关文档。
---

# gcoll Plugin Workflow

使用本技能前先阅读 `plugins/AGENTS.md` 和 `docs/current/03-插件与点位协议.md`。

## 工作流

1. 判断变更影响范围：协议、宿主、SDK、内置插件、示例插件、前端配置页或文档。
2. 修改协议前确认是否属于公开契约破坏性变更。
3. 协议结构变更先更新 `docs/current/03-插件与点位协议.md`。
4. 修改 `api/proto/` 后同步生成代码、SDK、示例插件和测试。
5. 修改插件配置结构时同步宿主校验、设备配置、密钥引用和前端配置页。
6. 完成后运行相关 Go 测试或构建命令。

## 强制规则

- 插件默认独立进程运行，通信只使用 gRPC。
- 插件不得直接访问宿主数据库。
- 插件目录不保存最终运行配置。
- 南向插件只提交采集记录，不直接转发目标系统。
- 北向插件只接收过滤后的记录，不主动采集设备。
- 敏感值必须通过宿主密钥存储引用。
- 背压状态必须能传递给南向插件。
- 插件自定义页面不得绕过宿主主题、权限、保存状态和审计提示。
- 内置 Modbus TCP 插件使用 `plugins/builtin/modbus_tcp/`，插件 ID 固定为 `com.gcoll.modbus-tcp`。
- Modbus TCP 设备配置和点位表只能由宿主保存并下发，插件目录不得保存最终运行配置。
- Modbus TCP 读取点位必须按区域和地址排序并批量读取，批量上限不得超过协议限制；网络高延迟、超时或错误时必须降低批量长度。
- Modbus TCP 写入只允许 `coil` 和 `holding_register`，不得写入 `discrete_input` 或 `input_register`。
- Modbus TCP 调试日志只保存有限窗口和报文摘要，不得把采集明细长期落库。

## 验证

- Go 代码变更运行 `go test ./...` 或对应子目录测试。
- 前端配置页变更运行 `pnpm --dir frontend/web build`。
- 协议变更在最终说明中列出同步过的 SDK、示例和文档。
