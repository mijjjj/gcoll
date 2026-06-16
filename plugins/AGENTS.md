# plugins Agent 指南

## 模块定位

`plugins/` 保存 gcoll 插件 SDK、内置插件和示例插件：

- `plugins/sdk-go/`：Go 插件 SDK。
- `plugins/builtin/`：内置插件。
- `plugins/examples/`：示例插件。

修改插件协议、SDK、宿主交互、内置插件或示例时使用 `.codex/skills/gcoll-plugin-workflow`。

## 核心规则

- 插件默认独立进程运行。
- 插件通信只使用 gRPC。
- 插件配置由宿主保存，插件目录不保存最终运行配置。
- 插件不得直接访问宿主数据库。
- 南向插件只提交采集记录，不直接转发目标系统。
- 北向插件只接收过滤后的记录，不主动采集设备。
- 敏感值必须通过宿主密钥存储引用。

## 协议和 SDK

- 修改 `api/proto/` 后必须同步 SDK、示例插件和 `docs/current/03-插件与点位协议.md`。
- SDK 面向第三方后必须显式版本化。
- 示例插件只演示推荐模式，不引入绕过宿主配置、权限或密钥存储的捷径。
- 内置 Modbus TCP 插件使用 `plugins/builtin/modbus_tcp/`，点位扩展字段放入通用点位表 `metadata`，不得绕过宿主点位 API。
- Modbus TCP 读取点位必须按区域和地址排序并批量读取，批量上限不得超过协议限制；网络高延迟或错误时必须降低批量长度。
- Modbus TCP 调试日志只能保存有限窗口和报文摘要，不得保存敏感配置明文或长期采集明细。

## 验证

- 修改 Go SDK 或 Go 示例后运行对应目录的 `go test ./...`。
- 修改协议后运行代码生成命令并说明生成范围。
