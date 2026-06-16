-- Modbus TCP 设备配置由宿主持久化，插件只通过协议读取运行配置。
CREATE TABLE IF NOT EXISTS modbus_tcp_device_profiles (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL,
  plugin_id TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  host TEXT NOT NULL,
  port INTEGER NOT NULL DEFAULT 502,
  unit_id INTEGER NOT NULL,
  timeout_ms INTEGER NOT NULL DEFAULT 2000,
  poll_interval_ms INTEGER NOT NULL DEFAULT 1000,
  report_mode TEXT NOT NULL DEFAULT 'change',
  max_coil_batch INTEGER NOT NULL DEFAULT 512,
  max_register_batch INTEGER NOT NULL DEFAULT 64,
  low_latency_ms INTEGER NOT NULL DEFAULT 80,
  high_latency_ms INTEGER NOT NULL DEFAULT 1000,
  debug_enabled INTEGER NOT NULL DEFAULT 0,
  enabled INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT,
  CHECK (plugin_id = 'com.gcoll.modbus-tcp'),
  CHECK (port BETWEEN 1 AND 65535),
  CHECK (unit_id BETWEEN 0 AND 247),
  CHECK (timeout_ms >= 100),
  CHECK (poll_interval_ms >= 100),
  CHECK (report_mode IN ('change', 'all')),
  CHECK (max_coil_batch BETWEEN 1 AND 2000),
  CHECK (max_register_batch BETWEEN 1 AND 125)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_modbus_tcp_device_profiles_device
  ON modbus_tcp_device_profiles(device_id, plugin_id)
  WHERE deleted_at IS NULL;

-- 调试日志只保存有限窗口，采集明细不落库。
CREATE TABLE IF NOT EXISTS modbus_tcp_debug_logs (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL,
  task_id TEXT,
  point_id TEXT,
  trace_id TEXT,
  level TEXT NOT NULL,
  message TEXT NOT NULL,
  area TEXT,
  address INTEGER,
  latency_ms INTEGER,
  raw_hex TEXT,
  fields_json TEXT NOT NULL DEFAULT '{}',
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR'))
);

CREATE INDEX IF NOT EXISTS idx_modbus_tcp_debug_logs_device_time
  ON modbus_tcp_debug_logs(device_id, created_at DESC);

-- 表注释：modbus_tcp_device_profiles：Modbus TCP 设备级协议配置，保存宿主下发给插件的连接和采集参数。
-- 字段注释：modbus_tcp_device_profiles.id：配置记录唯一标识。
-- 字段注释：modbus_tcp_device_profiles.device_id：所属设备 ID。
-- 字段注释：modbus_tcp_device_profiles.plugin_id：插件 ID，固定为 com.gcoll.modbus-tcp。
-- 字段注释：modbus_tcp_device_profiles.version：当前设备协议配置版本号。
-- 字段注释：modbus_tcp_device_profiles.host：Modbus TCP 设备主机地址。
-- 字段注释：modbus_tcp_device_profiles.port：Modbus TCP 端口。
-- 字段注释：modbus_tcp_device_profiles.unit_id：Modbus 单元 ID。
-- 字段注释：modbus_tcp_device_profiles.timeout_ms：单次请求超时时间，单位毫秒。
-- 字段注释：modbus_tcp_device_profiles.poll_interval_ms：采集轮询间隔，单位毫秒。
-- 字段注释：modbus_tcp_device_profiles.report_mode：上报模式，支持 change 或 all。
-- 字段注释：modbus_tcp_device_profiles.max_coil_batch：线圈和离散输入单次批量读取上限。
-- 字段注释：modbus_tcp_device_profiles.max_register_batch：寄存器单次批量读取上限。
-- 字段注释：modbus_tcp_device_profiles.low_latency_ms：自适应批量策略的低延迟阈值。
-- 字段注释：modbus_tcp_device_profiles.high_latency_ms：自适应批量策略的高延迟阈值。
-- 字段注释：modbus_tcp_device_profiles.debug_enabled：是否启用调试模式。
-- 字段注释：modbus_tcp_device_profiles.enabled：该设备协议配置是否启用。
-- 字段注释：modbus_tcp_device_profiles.created_at：记录创建时间。
-- 字段注释：modbus_tcp_device_profiles.updated_at：记录更新时间。
-- 字段注释：modbus_tcp_device_profiles.deleted_at：软删除时间。
-- 表注释：modbus_tcp_debug_logs：Modbus TCP 有限窗口调试日志，不保存长期采集明细。
-- 字段注释：modbus_tcp_debug_logs.id：调试日志唯一标识。
-- 字段注释：modbus_tcp_debug_logs.device_id：关联设备 ID。
-- 字段注释：modbus_tcp_debug_logs.task_id：关联采集任务 ID。
-- 字段注释：modbus_tcp_debug_logs.point_id：关联点位 ID。
-- 字段注释：modbus_tcp_debug_logs.trace_id：调用链追踪 ID。
-- 字段注释：modbus_tcp_debug_logs.level：日志级别。
-- 字段注释：modbus_tcp_debug_logs.message：日志消息。
-- 字段注释：modbus_tcp_debug_logs.area：本次请求涉及的 Modbus 数据区。
-- 字段注释：modbus_tcp_debug_logs.address：本次请求起始地址。
-- 字段注释：modbus_tcp_debug_logs.latency_ms：本次请求耗时，单位毫秒。
-- 字段注释：modbus_tcp_debug_logs.raw_hex：原始响应摘要十六进制文本。
-- 字段注释：modbus_tcp_debug_logs.fields_json：调试扩展字段 JSON。
-- 字段注释：modbus_tcp_debug_logs.created_at：日志创建时间。
