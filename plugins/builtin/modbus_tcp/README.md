# Modbus TCP 采集插件

本插件是 gcoll 的首个本地南向插件，用于通过 Modbus TCP 读取和写入工业设备点位。

## 支持能力

- 读取线圈、离散输入、保持寄存器和输入寄存器。
- 写入线圈和保持寄存器。
- 从宿主读取设备配置和通用点位表。
- 按区域和地址对读取点位重新排序，合并连续点位为批量读取请求。
- 根据近期读取延迟和错误动态调整线圈、寄存器批量读取长度。
- 支持 `change` 和 `all` 两种上报模式。
- 支持调试模式日志缓冲，由宿主采集并在控制台展示。

## 点位 metadata 约定

通用点位表字段保持平台统一，Modbus 扩展信息进入 `metadata`：

| 字段 | 说明 |
| --- | --- |
| `area` | `coil`、`discrete_input`、`holding_register`、`input_register` |
| `mode` | `read` 或 `write` |
| `quantity` | 点位占用线圈或寄存器数量 |
| `valueType` | 插件侧解码使用的值类型，例如 `bool`、`uint16`、`float32` |
| `byteOrder` | `big` 或 `little` |
| `wordOrder` | `big` 或 `little` |
| `scale` | 数值缩放系数 |
| `offset` | 数值偏移 |

插件目录不保存最终运行配置，设备配置和点位表必须由宿主持久化和下发。
