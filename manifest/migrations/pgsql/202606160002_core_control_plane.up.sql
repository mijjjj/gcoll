-- 控制面核心表：插件、设备、设备配置和通用点位表。
CREATE TABLE IF NOT EXISTS schema_migrations (
  version TEXT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS plugins (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  runtime TEXT NOT NULL,
  protocol TEXT NOT NULL,
  status TEXT NOT NULL,
  active_version TEXT NOT NULL,
  source TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  installed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CHECK (type IN ('system', 'southbound', 'northbound')),
  CHECK (runtime IN ('process')),
  CHECK (protocol IN ('grpc')),
  CHECK (status IN ('installed', 'enabled', 'disabled', 'running', 'stopped', 'failed'))
);

CREATE TABLE IF NOT EXISTS plugin_versions (
  id TEXT PRIMARY KEY,
  plugin_id TEXT NOT NULL REFERENCES plugins(id),
  version TEXT NOT NULL,
  package_path TEXT NOT NULL DEFAULT '',
  manifest_json JSONB NOT NULL,
  permissions_json JSONB NOT NULL DEFAULT '[]'::JSONB,
  capabilities_json JSONB NOT NULL DEFAULT '[]'::JSONB,
  checksum TEXT NOT NULL DEFAULT '',
  active BOOLEAN NOT NULL DEFAULT FALSE,
  installed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_plugin_versions_plugin_version
  ON plugin_versions(plugin_id, version)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS plugin_config_schemas (
  id TEXT PRIMARY KEY,
  plugin_id TEXT NOT NULL REFERENCES plugins(id),
  plugin_version_id TEXT NOT NULL REFERENCES plugin_versions(id),
  schema_version INTEGER NOT NULL DEFAULT 1,
  schema_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_plugin_config_schemas_plugin
  ON plugin_config_schemas(plugin_id, plugin_version_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS device_groups (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS devices (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  code TEXT NOT NULL,
  group_id TEXT NOT NULL REFERENCES device_groups(id),
  plugin_id TEXT NOT NULL REFERENCES plugins(id),
  status TEXT NOT NULL DEFAULT 'offline',
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  report_mode TEXT NOT NULL DEFAULT 'change',
  last_seen_at TIMESTAMPTZ,
  description TEXT NOT NULL DEFAULT '',
  metadata_json JSONB NOT NULL DEFAULT '{}'::JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CHECK (status IN ('online', 'offline', 'disabled', 'error')),
  CHECK (report_mode IN ('change', 'all'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_code
  ON devices(code)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_devices_group
  ON devices(group_id)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_devices_plugin
  ON devices(plugin_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS plugin_device_configs (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL REFERENCES devices(id),
  plugin_id TEXT NOT NULL REFERENCES plugins(id),
  version INTEGER NOT NULL DEFAULT 1,
  config_json JSONB NOT NULL,
  report_mode TEXT NOT NULL DEFAULT 'change',
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CHECK (report_mode IN ('change', 'all'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_plugin_device_configs_active
  ON plugin_device_configs(device_id, plugin_id, active)
  WHERE deleted_at IS NULL AND active = TRUE;

CREATE TABLE IF NOT EXISTS plugin_device_config_versions (
  id TEXT PRIMARY KEY,
  config_id TEXT NOT NULL REFERENCES plugin_device_configs(id),
  device_id TEXT NOT NULL REFERENCES devices(id),
  plugin_id TEXT NOT NULL REFERENCES plugins(id),
  version INTEGER NOT NULL,
  config_json JSONB NOT NULL,
  change_note TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS device_points (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL REFERENCES devices(id),
  plugin_id TEXT NOT NULL REFERENCES plugins(id),
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  address TEXT NOT NULL,
  value_type TEXT NOT NULL,
  unit TEXT NOT NULL DEFAULT '',
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  tags_json JSONB NOT NULL DEFAULT '{}'::JSONB,
  metadata_json JSONB NOT NULL DEFAULT '{}'::JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CHECK (value_type IN ('bool', 'int', 'float', 'string', 'bytes', 'datetime', 'json'))
);

CREATE INDEX IF NOT EXISTS idx_device_points_device
  ON device_points(device_id)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_device_points_plugin
  ON device_points(plugin_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS device_point_versions (
  id TEXT PRIMARY KEY,
  point_id TEXT NOT NULL REFERENCES device_points(id),
  device_id TEXT NOT NULL REFERENCES devices(id),
  plugin_id TEXT NOT NULL REFERENCES plugins(id),
  version INTEGER NOT NULL,
  snapshot_json JSONB NOT NULL,
  change_note TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS collection_tasks (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  device_id TEXT NOT NULL REFERENCES devices(id),
  south_plugin_id TEXT NOT NULL REFERENCES plugins(id),
  report_mode TEXT NOT NULL DEFAULT 'change',
  status TEXT NOT NULL DEFAULT 'stopped',
  rule_hit_rate TEXT NOT NULL DEFAULT '0%',
  rate TEXT NOT NULL DEFAULT '0 条/秒',
  last_collected_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CHECK (report_mode IN ('change', 'all')),
  CHECK (status IN ('running', 'stopped', 'disabled', 'error'))
);

CREATE INDEX IF NOT EXISTS idx_collection_tasks_device
  ON collection_tasks(device_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS runtime_events (
  id TEXT PRIMARY KEY,
  time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  level TEXT NOT NULL,
  source TEXT NOT NULL,
  plugin_id TEXT,
  device_id TEXT,
  task_id TEXT,
  message TEXT NOT NULL,
  trace_id TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR'))
);

CREATE INDEX IF NOT EXISTS idx_runtime_events_time
  ON runtime_events(time DESC);

INSERT INTO plugins (id, name, type, runtime, protocol, status, active_version, source, description, enabled, installed_at, updated_at)
VALUES ('com.gcoll.modbus-tcp', 'Modbus TCP 采集', 'southbound', 'process', 'grpc', 'running', '0.1.0', 'builtin', '内置 Modbus TCP 南向采集插件。', TRUE, '2026-06-16 10:15:00+08', '2026-06-16 10:15:00+08')
ON CONFLICT(id) DO UPDATE SET
  name = excluded.name,
  type = excluded.type,
  runtime = excluded.runtime,
  protocol = excluded.protocol,
  status = excluded.status,
  active_version = excluded.active_version,
  source = excluded.source,
  description = excluded.description,
  enabled = excluded.enabled,
  updated_at = excluded.updated_at;

INSERT INTO plugin_versions (id, plugin_id, version, package_path, manifest_json, permissions_json, capabilities_json, active, installed_at)
VALUES (
  'pv-com-gcoll-modbus-tcp-0-1-0',
  'com.gcoll.modbus-tcp',
  '0.1.0',
  'plugins/builtin/modbus_tcp',
  '{"schemaVersion":1,"id":"com.gcoll.modbus-tcp","name":"Modbus TCP 采集","type":"southbound","version":"0.1.0","runtime":"process","protocol":"grpc"}'::JSONB,
  '["network.outbound","config.read","runtime.events"]'::JSONB,
  '["southbound.collector.modbus-tcp"]'::JSONB,
  TRUE,
  '2026-06-16 10:15:00+08'
)
ON CONFLICT(plugin_id, version) WHERE deleted_at IS NULL DO UPDATE SET
  manifest_json = excluded.manifest_json,
  permissions_json = excluded.permissions_json,
  capabilities_json = excluded.capabilities_json,
  active = excluded.active,
  updated_at = NOW();

INSERT INTO plugin_config_schemas (id, plugin_id, plugin_version_id, schema_version, schema_json)
VALUES (
  'pcs-com-gcoll-modbus-tcp-0-1-0',
  'com.gcoll.modbus-tcp',
  'pv-com-gcoll-modbus-tcp-0-1-0',
  1,
  '{"type":"object","required":["host","port","unitId"],"properties":{"host":{"type":"string","title":"设备地址"},"port":{"type":"number","title":"端口"},"unitId":{"type":"number","title":"单元 ID"},"timeoutMs":{"type":"number","title":"超时时间"},"pollIntervalMs":{"type":"number","title":"轮询间隔"},"reportMode":{"type":"string","enum":["change","all"],"title":"上报模式"}}}'::JSONB
)
ON CONFLICT(id) DO UPDATE SET schema_json = excluded.schema_json;

INSERT INTO device_groups (id, name, sort_order)
VALUES
  ('edge', '边缘现场', 10),
  ('test', '测试分组', 20)
ON CONFLICT(id) DO UPDATE SET
  name = excluded.name,
  sort_order = excluded.sort_order,
  updated_at = NOW();

INSERT INTO devices (id, name, code, group_id, plugin_id, status, enabled, report_mode, last_seen_at, description)
VALUES
  ('dev-edge-gw-a01', '边缘网关 A01', 'DEV-EDGE-A01', 'edge', 'com.gcoll.modbus-tcp', 'online', TRUE, 'change', '2026-06-16 10:30:18+08', '用于验证 Modbus TCP 采集、过滤、转发闭环的本地网关。'),
  ('dev-sim-line-b02', '模拟产线 B02', 'DEV-SIM-B02', 'edge', 'com.gcoll.modbus-tcp', 'offline', FALSE, 'change', NULL, '保留给现场调试的模拟设备。')
ON CONFLICT(id) DO UPDATE SET
  name = excluded.name,
  code = excluded.code,
  group_id = excluded.group_id,
  plugin_id = excluded.plugin_id,
  status = excluded.status,
  enabled = excluded.enabled,
  report_mode = excluded.report_mode,
  last_seen_at = excluded.last_seen_at,
  description = excluded.description,
  updated_at = NOW();

INSERT INTO plugin_device_configs (id, device_id, plugin_id, version, config_json, report_mode, enabled, active)
VALUES
  ('pdc-dev-edge-gw-a01-modbus', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 1, '{"host":"192.168.10.25","port":502,"unitId":1,"timeoutMs":2000,"pollIntervalMs":1000,"reportMode":"change","debugEnabled":true,"maxCoilBatch":512,"maxRegisterBatch":64,"lowLatencyMs":80,"highLatencyMs":1000}'::JSONB, 'change', TRUE, TRUE),
  ('pdc-dev-sim-line-b02-modbus', 'dev-sim-line-b02', 'com.gcoll.modbus-tcp', 1, '{"host":"192.168.10.88","port":502,"unitId":2,"timeoutMs":3000,"pollIntervalMs":1500,"reportMode":"change","debugEnabled":false,"maxCoilBatch":256,"maxRegisterBatch":32,"lowLatencyMs":100,"highLatencyMs":1200}'::JSONB, 'change', FALSE, TRUE)
ON CONFLICT(id) DO UPDATE SET
  config_json = excluded.config_json,
  report_mode = excluded.report_mode,
  enabled = excluded.enabled,
  active = excluded.active,
  updated_at = NOW();

INSERT INTO plugin_device_config_versions (id, config_id, device_id, plugin_id, version, config_json, change_note)
VALUES
  ('pdcv-dev-edge-gw-a01-modbus-1', 'pdc-dev-edge-gw-a01-modbus', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 1, '{"host":"192.168.10.25","port":502,"unitId":1,"timeoutMs":2000,"pollIntervalMs":1000,"reportMode":"change","debugEnabled":true,"maxCoilBatch":512,"maxRegisterBatch":64,"lowLatencyMs":80,"highLatencyMs":1000}'::JSONB, '初始化内置示例配置'),
  ('pdcv-dev-sim-line-b02-modbus-1', 'pdc-dev-sim-line-b02-modbus', 'dev-sim-line-b02', 'com.gcoll.modbus-tcp', 1, '{"host":"192.168.10.88","port":502,"unitId":2,"timeoutMs":3000,"pollIntervalMs":1500,"reportMode":"change","debugEnabled":false,"maxCoilBatch":256,"maxRegisterBatch":32,"lowLatencyMs":100,"highLatencyMs":1200}'::JSONB, '初始化内置示例配置')
ON CONFLICT(id) DO NOTHING;

INSERT INTO device_points (id, device_id, plugin_id, name, description, address, value_type, unit, enabled, tags_json, metadata_json)
VALUES
  ('pt-temperature', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 'TEMP_01', '环境温度', 'holding_register:40001', 'float', '℃', TRUE, '{"area":"A","kind":"environment"}'::JSONB, '{"area":"holding_register","mode":"read","quantity":2,"valueType":"float32","byteOrder":"big","wordOrder":"big","scale":1,"offset":0}'::JSONB),
  ('pt-pressure', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 'PRESS_01', '管线压力', 'holding_register:40003', 'float', 'kPa', TRUE, '{"area":"A","kind":"process"}'::JSONB, '{"area":"holding_register","mode":"read","quantity":2,"valueType":"float32","byteOrder":"big","wordOrder":"big","scale":1,"offset":0}'::JSONB),
  ('pt-motor-state', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 'MOTOR_RUN', '电机运行状态', 'coil:00001', 'bool', '', TRUE, '{"area":"A","kind":"status"}'::JSONB, '{"area":"coil","mode":"read","quantity":1,"valueType":"bool","byteOrder":"big","wordOrder":"big","scale":1,"offset":0}'::JSONB),
  ('pt-energy', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 'ENERGY_TOTAL', '累计能耗', 'input_register:30001', 'float', 'kWh', TRUE, '{"area":"A","kind":"meter"}'::JSONB, '{"area":"input_register","mode":"read","quantity":2,"valueType":"float32","byteOrder":"big","wordOrder":"little","scale":1,"offset":0}'::JSONB),
  ('pt-speed-set', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 'SPEED_SET', '速度设定值', 'holding_register:40110', 'int', 'rpm', TRUE, '{"area":"A","kind":"write"}'::JSONB, '{"area":"holding_register","mode":"write","quantity":1,"valueType":"uint16","byteOrder":"big","wordOrder":"big","scale":1,"offset":0}'::JSONB),
  ('pt-emergency', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 'EMERGENCY_STOP', '急停输入状态', 'discrete_input:10001', 'bool', '', TRUE, '{"area":"A","kind":"safety"}'::JSONB, '{"area":"discrete_input","mode":"read","quantity":1,"valueType":"bool","byteOrder":"big","wordOrder":"big","scale":1,"offset":0}'::JSONB)
ON CONFLICT(id) DO UPDATE SET
  name = excluded.name,
  description = excluded.description,
  address = excluded.address,
  value_type = excluded.value_type,
  unit = excluded.unit,
  enabled = excluded.enabled,
  tags_json = excluded.tags_json,
  metadata_json = excluded.metadata_json,
  updated_at = NOW();

INSERT INTO device_point_versions (id, point_id, device_id, plugin_id, version, snapshot_json, change_note)
SELECT 'dpv-' || id || '-1', id, device_id, plugin_id, 1,
       jsonb_build_object('name', name, 'address', address, 'valueType', value_type),
       '初始化通用点位表'
FROM device_points
WHERE device_id = 'dev-edge-gw-a01'
ON CONFLICT(id) DO NOTHING;

INSERT INTO collection_tasks (id, name, device_id, south_plugin_id, report_mode, status, rule_hit_rate, rate, last_collected_at)
VALUES ('task-modbus-a01', '样例 Modbus TCP 采集链路', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 'change', 'running', '72%', '128 条/秒', '2026-06-16 10:30:18+08')
ON CONFLICT(id) DO UPDATE SET
  name = excluded.name,
  device_id = excluded.device_id,
  south_plugin_id = excluded.south_plugin_id,
  report_mode = excluded.report_mode,
  status = excluded.status,
  rule_hit_rate = excluded.rule_hit_rate,
  rate = excluded.rate,
  last_collected_at = excluded.last_collected_at,
  updated_at = NOW();

INSERT INTO runtime_events (id, time, level, source, plugin_id, device_id, task_id, message, trace_id)
VALUES
  ('evt-001', '2026-06-16 10:30:18+08', 'INFO', 'collector', 'com.gcoll.modbus-tcp', 'dev-edge-gw-a01', 'task-modbus-a01', '已接收 128 条采集记录并写入内存缓冲。', 'trace-demo-001'),
  ('evt-002', '2026-06-16 10:30:18+08', 'INFO', 'pipeline', NULL, 'dev-edge-gw-a01', 'task-modbus-a01', '规则过滤命中 92 条记录，准备交给北向转发。', 'trace-demo-001')
ON CONFLICT(id) DO UPDATE SET
  time = excluded.time,
  level = excluded.level,
  source = excluded.source,
  plugin_id = excluded.plugin_id,
  device_id = excluded.device_id,
  task_id = excluded.task_id,
  message = excluded.message,
  trace_id = excluded.trace_id;

INSERT INTO modbus_tcp_device_profiles (
  id, device_id, plugin_id, version, host, port, unit_id, timeout_ms, poll_interval_ms,
  report_mode, max_coil_batch, max_register_batch, low_latency_ms, high_latency_ms,
  debug_enabled, enabled
)
VALUES
  ('mdp-dev-edge-gw-a01', 'dev-edge-gw-a01', 'com.gcoll.modbus-tcp', 1, '192.168.10.25', 502, 1, 2000, 1000, 'change', 512, 64, 80, 1000, TRUE, TRUE),
  ('mdp-dev-sim-line-b02', 'dev-sim-line-b02', 'com.gcoll.modbus-tcp', 1, '192.168.10.88', 502, 2, 3000, 1500, 'change', 256, 32, 100, 1200, FALSE, FALSE)
ON CONFLICT(device_id, plugin_id) WHERE deleted_at IS NULL DO UPDATE SET
  host = excluded.host,
  port = excluded.port,
  unit_id = excluded.unit_id,
  timeout_ms = excluded.timeout_ms,
  poll_interval_ms = excluded.poll_interval_ms,
  report_mode = excluded.report_mode,
  max_coil_batch = excluded.max_coil_batch,
  max_register_batch = excluded.max_register_batch,
  low_latency_ms = excluded.low_latency_ms,
  high_latency_ms = excluded.high_latency_ms,
  debug_enabled = excluded.debug_enabled,
  enabled = excluded.enabled,
  updated_at = NOW();

INSERT INTO modbus_tcp_debug_logs (id, device_id, task_id, point_id, trace_id, level, message, area, address, latency_ms, raw_hex, fields_json, created_at)
VALUES
  ('modbus-log-001', 'dev-edge-gw-a01', 'task-modbus-a01', NULL, 'trace-demo-001', 'DEBUG', '批量读取成功', 'holding_register', 0, 22, '41CC000042CA999A', '{}'::JSONB, '2026-06-16 10:30:18.101+08'),
  ('modbus-log-002', 'dev-edge-gw-a01', 'task-modbus-a01', NULL, 'trace-demo-001', 'DEBUG', '线圈读取成功', 'coil', 0, 16, '01', '{}'::JSONB, '2026-06-16 10:30:18.118+08'),
  ('modbus-log-003', 'dev-edge-gw-a01', 'task-modbus-a01', NULL, 'trace-demo-001', 'INFO', '自适应读取上限保持稳定', 'holding_register', 0, 0, '', '{}'::JSONB, '2026-06-16 10:30:19.002+08')
ON CONFLICT(id) DO UPDATE SET
  message = excluded.message,
  area = excluded.area,
  address = excluded.address,
  latency_ms = excluded.latency_ms,
  raw_hex = excluded.raw_hex,
  created_at = excluded.created_at;

COMMENT ON TABLE schema_migrations IS '记录已经成功应用的数据库迁移版本。';
COMMENT ON COLUMN schema_migrations.version IS '迁移版本号，通常来自迁移文件名。';
COMMENT ON COLUMN schema_migrations.applied_at IS '迁移应用完成时间。';

COMMENT ON TABLE plugins IS '插件主表，保存插件当前活跃版本和运行状态。';
COMMENT ON COLUMN plugins.id IS '插件唯一标识。';
COMMENT ON COLUMN plugins.name IS '插件显示名称。';
COMMENT ON COLUMN plugins.type IS '插件类型，支持 system、southbound、northbound。';
COMMENT ON COLUMN plugins.runtime IS '插件运行时类型。';
COMMENT ON COLUMN plugins.protocol IS '插件通信协议。';
COMMENT ON COLUMN plugins.status IS '插件当前管理状态。';
COMMENT ON COLUMN plugins.active_version IS '当前活跃插件版本。';
COMMENT ON COLUMN plugins.source IS '插件来源，例如 builtin 或 imported。';
COMMENT ON COLUMN plugins.description IS '插件说明。';
COMMENT ON COLUMN plugins.enabled IS '插件是否启用。';
COMMENT ON COLUMN plugins.installed_at IS '插件安装时间。';
COMMENT ON COLUMN plugins.updated_at IS '插件更新时间。';
COMMENT ON COLUMN plugins.deleted_at IS '软删除时间。';

COMMENT ON TABLE plugin_versions IS '插件版本表，保存每个插件版本的清单、权限和能力声明。';
COMMENT ON COLUMN plugin_versions.id IS '插件版本唯一标识。';
COMMENT ON COLUMN plugin_versions.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN plugin_versions.version IS '插件版本号。';
COMMENT ON COLUMN plugin_versions.package_path IS '插件包或内置插件目录路径。';
COMMENT ON COLUMN plugin_versions.manifest_json IS '插件清单 JSON。';
COMMENT ON COLUMN plugin_versions.permissions_json IS '插件权限声明 JSON。';
COMMENT ON COLUMN plugin_versions.capabilities_json IS '插件能力声明 JSON。';
COMMENT ON COLUMN plugin_versions.checksum IS '插件包校验和。';
COMMENT ON COLUMN plugin_versions.active IS '该版本是否为活跃版本。';
COMMENT ON COLUMN plugin_versions.installed_at IS '版本安装时间。';
COMMENT ON COLUMN plugin_versions.created_at IS '记录创建时间。';
COMMENT ON COLUMN plugin_versions.updated_at IS '记录更新时间。';
COMMENT ON COLUMN plugin_versions.deleted_at IS '软删除时间。';

COMMENT ON TABLE plugin_config_schemas IS '插件配置结构表，保存插件版本对应的配置表单和服务端校验结构。';
COMMENT ON COLUMN plugin_config_schemas.id IS '配置结构唯一标识。';
COMMENT ON COLUMN plugin_config_schemas.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN plugin_config_schemas.plugin_version_id IS '所属插件版本 ID。';
COMMENT ON COLUMN plugin_config_schemas.schema_version IS '配置结构版本号。';
COMMENT ON COLUMN plugin_config_schemas.schema_json IS '配置结构 JSON。';
COMMENT ON COLUMN plugin_config_schemas.created_at IS '记录创建时间。';
COMMENT ON COLUMN plugin_config_schemas.deleted_at IS '软删除时间。';

COMMENT ON TABLE device_groups IS '设备分组表，用于控制台设备列表归类。';
COMMENT ON COLUMN device_groups.id IS '设备分组唯一标识。';
COMMENT ON COLUMN device_groups.name IS '设备分组名称。';
COMMENT ON COLUMN device_groups.sort_order IS '设备分组排序值。';
COMMENT ON COLUMN device_groups.created_at IS '记录创建时间。';
COMMENT ON COLUMN device_groups.updated_at IS '记录更新时间。';
COMMENT ON COLUMN device_groups.deleted_at IS '软删除时间。';

COMMENT ON TABLE devices IS '设备档案表，保存设备基础信息和绑定的南向插件。';
COMMENT ON COLUMN devices.id IS '设备唯一标识。';
COMMENT ON COLUMN devices.name IS '设备名称。';
COMMENT ON COLUMN devices.code IS '设备编码。';
COMMENT ON COLUMN devices.group_id IS '所属设备分组 ID。';
COMMENT ON COLUMN devices.plugin_id IS '设备绑定的南向插件 ID。';
COMMENT ON COLUMN devices.status IS '设备状态。';
COMMENT ON COLUMN devices.enabled IS '设备是否启用。';
COMMENT ON COLUMN devices.report_mode IS '设备上报模式。';
COMMENT ON COLUMN devices.last_seen_at IS '设备最近在线或采集时间。';
COMMENT ON COLUMN devices.description IS '设备说明。';
COMMENT ON COLUMN devices.metadata_json IS '设备扩展元数据 JSON。';
COMMENT ON COLUMN devices.created_at IS '记录创建时间。';
COMMENT ON COLUMN devices.updated_at IS '记录更新时间。';
COMMENT ON COLUMN devices.deleted_at IS '软删除时间。';

COMMENT ON TABLE plugin_device_configs IS '设备插件配置表，保存设备维度最终运行配置。';
COMMENT ON COLUMN plugin_device_configs.id IS '设备插件配置唯一标识。';
COMMENT ON COLUMN plugin_device_configs.device_id IS '所属设备 ID。';
COMMENT ON COLUMN plugin_device_configs.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN plugin_device_configs.version IS '当前配置版本号。';
COMMENT ON COLUMN plugin_device_configs.config_json IS '设备插件配置 JSON。';
COMMENT ON COLUMN plugin_device_configs.report_mode IS '配置生效的上报模式。';
COMMENT ON COLUMN plugin_device_configs.enabled IS '配置是否启用。';
COMMENT ON COLUMN plugin_device_configs.active IS '是否为当前活跃配置。';
COMMENT ON COLUMN plugin_device_configs.created_at IS '记录创建时间。';
COMMENT ON COLUMN plugin_device_configs.updated_at IS '记录更新时间。';
COMMENT ON COLUMN plugin_device_configs.deleted_at IS '软删除时间。';

COMMENT ON TABLE plugin_device_config_versions IS '设备插件配置版本表，保存设备配置历史快照。';
COMMENT ON COLUMN plugin_device_config_versions.id IS '配置版本唯一标识。';
COMMENT ON COLUMN plugin_device_config_versions.config_id IS '所属设备插件配置 ID。';
COMMENT ON COLUMN plugin_device_config_versions.device_id IS '所属设备 ID。';
COMMENT ON COLUMN plugin_device_config_versions.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN plugin_device_config_versions.version IS '配置版本号。';
COMMENT ON COLUMN plugin_device_config_versions.config_json IS '配置版本快照 JSON。';
COMMENT ON COLUMN plugin_device_config_versions.change_note IS '版本变更说明。';
COMMENT ON COLUMN plugin_device_config_versions.created_at IS '版本创建时间。';

COMMENT ON TABLE device_points IS '通用点位表，保存平台统一点位字段。';
COMMENT ON COLUMN device_points.id IS '点位唯一标识。';
COMMENT ON COLUMN device_points.device_id IS '所属设备 ID。';
COMMENT ON COLUMN device_points.plugin_id IS '解释该点位的插件 ID。';
COMMENT ON COLUMN device_points.name IS '点位名称。';
COMMENT ON COLUMN device_points.description IS '点位说明。';
COMMENT ON COLUMN device_points.address IS '点位地址，由插件解释。';
COMMENT ON COLUMN device_points.value_type IS '平台通用值类型。';
COMMENT ON COLUMN device_points.unit IS '点位单位。';
COMMENT ON COLUMN device_points.enabled IS '点位是否启用。';
COMMENT ON COLUMN device_points.tags_json IS '点位标签 JSON。';
COMMENT ON COLUMN device_points.metadata_json IS '插件扩展元数据 JSON，不保存敏感值。';
COMMENT ON COLUMN device_points.created_at IS '记录创建时间。';
COMMENT ON COLUMN device_points.updated_at IS '记录更新时间。';
COMMENT ON COLUMN device_points.deleted_at IS '软删除时间。';

COMMENT ON TABLE device_point_versions IS '点位版本表，保存通用点位表历史快照。';
COMMENT ON COLUMN device_point_versions.id IS '点位版本唯一标识。';
COMMENT ON COLUMN device_point_versions.point_id IS '所属点位 ID。';
COMMENT ON COLUMN device_point_versions.device_id IS '所属设备 ID。';
COMMENT ON COLUMN device_point_versions.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN device_point_versions.version IS '点位版本号。';
COMMENT ON COLUMN device_point_versions.snapshot_json IS '点位版本快照 JSON。';
COMMENT ON COLUMN device_point_versions.change_note IS '版本变更说明。';
COMMENT ON COLUMN device_point_versions.created_at IS '版本创建时间。';

COMMENT ON TABLE collection_tasks IS '采集任务表，保存设备采集任务的控制面状态。';
COMMENT ON COLUMN collection_tasks.id IS '采集任务唯一标识。';
COMMENT ON COLUMN collection_tasks.name IS '采集任务名称。';
COMMENT ON COLUMN collection_tasks.device_id IS '采集任务关联设备 ID。';
COMMENT ON COLUMN collection_tasks.south_plugin_id IS '采集任务使用的南向插件 ID。';
COMMENT ON COLUMN collection_tasks.report_mode IS '采集任务上报模式。';
COMMENT ON COLUMN collection_tasks.status IS '采集任务状态。';
COMMENT ON COLUMN collection_tasks.rule_hit_rate IS '规则命中率展示值。';
COMMENT ON COLUMN collection_tasks.rate IS '采集速率展示值。';
COMMENT ON COLUMN collection_tasks.last_collected_at IS '最近采集时间。';
COMMENT ON COLUMN collection_tasks.created_at IS '记录创建时间。';
COMMENT ON COLUMN collection_tasks.updated_at IS '记录更新时间。';
COMMENT ON COLUMN collection_tasks.deleted_at IS '软删除时间。';

COMMENT ON TABLE runtime_events IS '运行事件表，保存低频运行事件和控制面日志。';
COMMENT ON COLUMN runtime_events.id IS '运行事件唯一标识。';
COMMENT ON COLUMN runtime_events.time IS '事件发生时间。';
COMMENT ON COLUMN runtime_events.level IS '事件级别。';
COMMENT ON COLUMN runtime_events.source IS '事件来源模块。';
COMMENT ON COLUMN runtime_events.plugin_id IS '关联插件 ID。';
COMMENT ON COLUMN runtime_events.device_id IS '关联设备 ID。';
COMMENT ON COLUMN runtime_events.task_id IS '关联任务 ID。';
COMMENT ON COLUMN runtime_events.message IS '事件消息。';
COMMENT ON COLUMN runtime_events.trace_id IS '调用链追踪 ID。';
COMMENT ON COLUMN runtime_events.created_at IS '记录创建时间。';
