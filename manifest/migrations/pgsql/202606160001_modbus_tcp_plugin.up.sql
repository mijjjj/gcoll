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
