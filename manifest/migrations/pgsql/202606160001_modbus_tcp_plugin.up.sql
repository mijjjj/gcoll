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
  debug_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
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

-- Modbus 扩展点位表保存协议解释字段，通用点位表仍保存平台通用字段。
CREATE TABLE IF NOT EXISTS modbus_tcp_point_profiles (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL,
  point_id TEXT NOT NULL,
  plugin_id TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  area TEXT NOT NULL,
  address INTEGER NOT NULL,
  quantity INTEGER NOT NULL DEFAULT 1,
  mode TEXT NOT NULL DEFAULT 'read',
  value_type TEXT NOT NULL,
  byte_order TEXT NOT NULL DEFAULT 'big',
  word_order TEXT NOT NULL DEFAULT 'big',
  scale DOUBLE PRECISION NOT NULL DEFAULT 1,
  offset DOUBLE PRECISION NOT NULL DEFAULT 0,
  report_mode TEXT NOT NULL DEFAULT 'change',
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CHECK (plugin_id = 'com.gcoll.modbus-tcp'),
  CHECK (area IN ('coil', 'discrete_input', 'holding_register', 'input_register')),
  CHECK (address BETWEEN 0 AND 65535),
  CHECK (quantity BETWEEN 1 AND 2000),
  CHECK (mode IN ('read', 'write')),
  CHECK (byte_order IN ('big', 'little')),
  CHECK (word_order IN ('big', 'little')),
  CHECK (report_mode IN ('change', 'all')),
  CHECK (NOT (mode = 'write' AND area IN ('discrete_input', 'input_register')))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_modbus_tcp_point_profiles_point
  ON modbus_tcp_point_profiles(point_id, plugin_id)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_modbus_tcp_point_profiles_read_order
  ON modbus_tcp_point_profiles(device_id, area, address)
  WHERE deleted_at IS NULL AND enabled = TRUE;

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
  fields_json JSONB NOT NULL DEFAULT '{}'::JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR'))
);

CREATE INDEX IF NOT EXISTS idx_modbus_tcp_debug_logs_device_time
  ON modbus_tcp_debug_logs(device_id, created_at DESC);

COMMENT ON TABLE modbus_tcp_device_profiles IS 'Modbus TCP 设备级协议配置，保存宿主下发给插件的连接和采集参数。';
COMMENT ON COLUMN modbus_tcp_device_profiles.id IS '配置记录唯一标识。';
COMMENT ON COLUMN modbus_tcp_device_profiles.device_id IS '所属设备 ID。';
COMMENT ON COLUMN modbus_tcp_device_profiles.plugin_id IS '插件 ID，固定为 com.gcoll.modbus-tcp。';
COMMENT ON COLUMN modbus_tcp_device_profiles.version IS '当前设备协议配置版本号。';
COMMENT ON COLUMN modbus_tcp_device_profiles.host IS 'Modbus TCP 设备主机地址。';
COMMENT ON COLUMN modbus_tcp_device_profiles.port IS 'Modbus TCP 端口。';
COMMENT ON COLUMN modbus_tcp_device_profiles.unit_id IS 'Modbus 单元 ID。';
COMMENT ON COLUMN modbus_tcp_device_profiles.timeout_ms IS '单次请求超时时间，单位毫秒。';
COMMENT ON COLUMN modbus_tcp_device_profiles.poll_interval_ms IS '采集轮询间隔，单位毫秒。';
COMMENT ON COLUMN modbus_tcp_device_profiles.report_mode IS '上报模式，支持 change 或 all。';
COMMENT ON COLUMN modbus_tcp_device_profiles.max_coil_batch IS '线圈和离散输入单次批量读取上限。';
COMMENT ON COLUMN modbus_tcp_device_profiles.max_register_batch IS '寄存器单次批量读取上限。';
COMMENT ON COLUMN modbus_tcp_device_profiles.low_latency_ms IS '自适应批量策略的低延迟阈值。';
COMMENT ON COLUMN modbus_tcp_device_profiles.high_latency_ms IS '自适应批量策略的高延迟阈值。';
COMMENT ON COLUMN modbus_tcp_device_profiles.debug_enabled IS '是否启用调试模式。';
COMMENT ON COLUMN modbus_tcp_device_profiles.enabled IS '该设备协议配置是否启用。';
COMMENT ON COLUMN modbus_tcp_device_profiles.created_at IS '记录创建时间。';
COMMENT ON COLUMN modbus_tcp_device_profiles.updated_at IS '记录更新时间。';
COMMENT ON COLUMN modbus_tcp_device_profiles.deleted_at IS '软删除时间。';

COMMENT ON TABLE modbus_tcp_point_profiles IS 'Modbus TCP 点位协议扩展配置，保存通用点位表无法表达的区域、地址和解码参数。';
COMMENT ON COLUMN modbus_tcp_point_profiles.id IS '点位扩展配置唯一标识。';
COMMENT ON COLUMN modbus_tcp_point_profiles.device_id IS '所属设备 ID。';
COMMENT ON COLUMN modbus_tcp_point_profiles.point_id IS '关联的通用点位 ID。';
COMMENT ON COLUMN modbus_tcp_point_profiles.plugin_id IS '插件 ID，固定为 com.gcoll.modbus-tcp。';
COMMENT ON COLUMN modbus_tcp_point_profiles.version IS '当前点位扩展配置版本号。';
COMMENT ON COLUMN modbus_tcp_point_profiles.area IS 'Modbus 数据区。';
COMMENT ON COLUMN modbus_tcp_point_profiles.address IS '从 0 开始的 Modbus 地址。';
COMMENT ON COLUMN modbus_tcp_point_profiles.quantity IS '点位占用的线圈或寄存器数量。';
COMMENT ON COLUMN modbus_tcp_point_profiles.mode IS '点位读写模式，支持 read 或 write。';
COMMENT ON COLUMN modbus_tcp_point_profiles.value_type IS '插件侧解码使用的值类型。';
COMMENT ON COLUMN modbus_tcp_point_profiles.byte_order IS '字节序。';
COMMENT ON COLUMN modbus_tcp_point_profiles.word_order IS '字序。';
COMMENT ON COLUMN modbus_tcp_point_profiles.scale IS '数值缩放系数。';
COMMENT ON COLUMN modbus_tcp_point_profiles.offset IS '数值偏移量。';
COMMENT ON COLUMN modbus_tcp_point_profiles.report_mode IS '点位上报模式。';
COMMENT ON COLUMN modbus_tcp_point_profiles.enabled IS '点位扩展配置是否启用。';
COMMENT ON COLUMN modbus_tcp_point_profiles.created_at IS '记录创建时间。';
COMMENT ON COLUMN modbus_tcp_point_profiles.updated_at IS '记录更新时间。';
COMMENT ON COLUMN modbus_tcp_point_profiles.deleted_at IS '软删除时间。';

COMMENT ON TABLE modbus_tcp_debug_logs IS 'Modbus TCP 有限窗口调试日志，不保存长期采集明细。';
COMMENT ON COLUMN modbus_tcp_debug_logs.id IS '调试日志唯一标识。';
COMMENT ON COLUMN modbus_tcp_debug_logs.device_id IS '关联设备 ID。';
COMMENT ON COLUMN modbus_tcp_debug_logs.task_id IS '关联采集任务 ID。';
COMMENT ON COLUMN modbus_tcp_debug_logs.point_id IS '关联点位 ID。';
COMMENT ON COLUMN modbus_tcp_debug_logs.trace_id IS '调用链追踪 ID。';
COMMENT ON COLUMN modbus_tcp_debug_logs.level IS '日志级别。';
COMMENT ON COLUMN modbus_tcp_debug_logs.message IS '日志消息。';
COMMENT ON COLUMN modbus_tcp_debug_logs.area IS '本次请求涉及的 Modbus 数据区。';
COMMENT ON COLUMN modbus_tcp_debug_logs.address IS '本次请求起始地址。';
COMMENT ON COLUMN modbus_tcp_debug_logs.latency_ms IS '本次请求耗时，单位毫秒。';
COMMENT ON COLUMN modbus_tcp_debug_logs.raw_hex IS '原始响应摘要十六进制文本。';
COMMENT ON COLUMN modbus_tcp_debug_logs.fields_json IS '调试扩展字段 JSON。';
COMMENT ON COLUMN modbus_tcp_debug_logs.created_at IS '日志创建时间。';
