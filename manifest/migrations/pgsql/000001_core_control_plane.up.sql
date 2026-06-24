-- 控制面核心表：插件清单快照、设备、设备插件配置、通用点位、采集任务和运行事件。
CREATE TABLE IF NOT EXISTS schema_migrations (
  version TEXT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE schema_migrations IS '记录已经成功应用的数据库迁移版本。';
COMMENT ON COLUMN schema_migrations.version IS '迁移版本号，通常来自迁移文件名。';
COMMENT ON COLUMN schema_migrations.applied_at IS '迁移应用完成时间。';

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
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CHECK (type IN ('system', 'southbound', 'northbound')),
  CHECK (runtime IN ('process')),
  CHECK (protocol IN ('grpc')),
  CHECK (status IN ('installed', 'enabled', 'disabled', 'running', 'stopped', 'failed'))
);

COMMENT ON TABLE plugins IS '插件清单快照表，保存导入或安装插件的控制面元数据；运行时事实来源仍是宿主内存注册表。';
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

CREATE TABLE IF NOT EXISTS plugin_versions (
  id TEXT PRIMARY KEY,
  plugin_id TEXT NOT NULL,
  version TEXT NOT NULL,
  package_path TEXT NOT NULL DEFAULT '',
  manifest_json JSONB NOT NULL,
  permissions_json JSONB NOT NULL DEFAULT '[]'::JSONB,
  capabilities_json JSONB NOT NULL DEFAULT '[]'::JSONB,
  checksum TEXT NOT NULL DEFAULT '',
  active BOOLEAN NOT NULL DEFAULT FALSE,
  installed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE plugin_versions IS '插件版本表，保存每个插件版本的清单、权限和能力声明。';
COMMENT ON COLUMN plugin_versions.id IS '插件版本唯一标识。';
COMMENT ON COLUMN plugin_versions.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN plugin_versions.version IS '插件版本号。';
COMMENT ON COLUMN plugin_versions.package_path IS '插件包或插件目录路径。';
COMMENT ON COLUMN plugin_versions.manifest_json IS '插件清单 JSON。';
COMMENT ON COLUMN plugin_versions.permissions_json IS '插件权限声明 JSON。';
COMMENT ON COLUMN plugin_versions.capabilities_json IS '插件能力声明 JSON。';
COMMENT ON COLUMN plugin_versions.checksum IS '插件包校验和。';
COMMENT ON COLUMN plugin_versions.active IS '该版本是否为活跃版本。';
COMMENT ON COLUMN plugin_versions.installed_at IS '版本安装时间。';
COMMENT ON COLUMN plugin_versions.created_at IS '记录创建时间。';
COMMENT ON COLUMN plugin_versions.updated_at IS '记录更新时间。';

CREATE UNIQUE INDEX IF NOT EXISTS idx_plugin_versions_plugin_version
  ON plugin_versions(plugin_id, version);

CREATE TABLE IF NOT EXISTS plugin_config_schemas (
  id TEXT PRIMARY KEY,
  plugin_id TEXT NOT NULL,
  plugin_version_id TEXT NOT NULL,
  schema_version INTEGER NOT NULL DEFAULT 1,
  schema_json JSONB NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE plugin_config_schemas IS '插件配置结构表，保存插件版本对应的配置表单和服务端校验结构。';
COMMENT ON COLUMN plugin_config_schemas.id IS '配置结构唯一标识。';
COMMENT ON COLUMN plugin_config_schemas.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN plugin_config_schemas.plugin_version_id IS '所属插件版本 ID。';
COMMENT ON COLUMN plugin_config_schemas.schema_version IS '配置结构版本号。';
COMMENT ON COLUMN plugin_config_schemas.schema_json IS '配置结构 JSON。';
COMMENT ON COLUMN plugin_config_schemas.created_at IS '记录创建时间。';

CREATE INDEX IF NOT EXISTS idx_plugin_config_schemas_plugin
  ON plugin_config_schemas(plugin_id, plugin_version_id);

CREATE TABLE IF NOT EXISTS device_groups (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE device_groups IS '设备分组表，用于控制台设备列表归类。';
COMMENT ON COLUMN device_groups.id IS '设备分组唯一标识。';
COMMENT ON COLUMN device_groups.name IS '设备分组名称。';
COMMENT ON COLUMN device_groups.sort_order IS '设备分组排序值。';
COMMENT ON COLUMN device_groups.created_at IS '记录创建时间。';
COMMENT ON COLUMN device_groups.updated_at IS '记录更新时间。';

CREATE TABLE IF NOT EXISTS devices (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  code TEXT NOT NULL,
  group_id TEXT NOT NULL REFERENCES device_groups(id),
  plugin_id TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'offline',
  enabled BOOLEAN NOT NULL DEFAULT FALSE,
  report_mode TEXT NOT NULL DEFAULT 'change',
  last_seen_at TIMESTAMPTZ,
  description TEXT NOT NULL DEFAULT '',
  metadata_json JSONB NOT NULL DEFAULT '{}'::JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CHECK (status IN ('online', 'offline', 'disabled', 'error')),
  CHECK (report_mode IN ('change', 'all'))
);

COMMENT ON TABLE devices IS '设备档案表，保存设备基础信息和绑定的南向插件 ID。';
COMMENT ON COLUMN devices.id IS '设备唯一标识。';
COMMENT ON COLUMN devices.name IS '设备名称。';
COMMENT ON COLUMN devices.code IS '设备编码。';
COMMENT ON COLUMN devices.group_id IS '所属设备分组 ID。';
COMMENT ON COLUMN devices.plugin_id IS '设备绑定的南向插件 ID，由插件宿主注册表解释。';
COMMENT ON COLUMN devices.status IS '设备状态。';
COMMENT ON COLUMN devices.enabled IS '设备是否启用。';
COMMENT ON COLUMN devices.report_mode IS '设备上报模式。';
COMMENT ON COLUMN devices.last_seen_at IS '设备最近在线或采集时间。';
COMMENT ON COLUMN devices.description IS '设备说明。';
COMMENT ON COLUMN devices.metadata_json IS '设备扩展元数据 JSON。';
COMMENT ON COLUMN devices.created_at IS '记录创建时间。';
COMMENT ON COLUMN devices.updated_at IS '记录更新时间。';

CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_code
  ON devices(code);

CREATE INDEX IF NOT EXISTS idx_devices_group
  ON devices(group_id);

CREATE INDEX IF NOT EXISTS idx_devices_plugin
  ON devices(plugin_id);

CREATE TABLE IF NOT EXISTS plugin_device_configs (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL REFERENCES devices(id),
  plugin_id TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  config_json JSONB NOT NULL,
  report_mode TEXT NOT NULL DEFAULT 'change',
  enabled BOOLEAN NOT NULL DEFAULT FALSE,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CHECK (report_mode IN ('change', 'all'))
);

COMMENT ON TABLE plugin_device_configs IS '设备插件配置表，保存设备维度最终运行配置。';
COMMENT ON COLUMN plugin_device_configs.id IS '设备插件配置唯一标识。';
COMMENT ON COLUMN plugin_device_configs.device_id IS '所属设备 ID。';
COMMENT ON COLUMN plugin_device_configs.plugin_id IS '所属插件 ID，由插件宿主注册表解释。';
COMMENT ON COLUMN plugin_device_configs.version IS '当前配置版本号。';
COMMENT ON COLUMN plugin_device_configs.config_json IS '设备插件配置 JSON，敏感值只保存密钥引用。';
COMMENT ON COLUMN plugin_device_configs.report_mode IS '配置生效的上报模式。';
COMMENT ON COLUMN plugin_device_configs.enabled IS '配置是否启用。';
COMMENT ON COLUMN plugin_device_configs.active IS '是否为当前活跃配置。';
COMMENT ON COLUMN plugin_device_configs.created_at IS '记录创建时间。';
COMMENT ON COLUMN plugin_device_configs.updated_at IS '记录更新时间。';

CREATE UNIQUE INDEX IF NOT EXISTS idx_plugin_device_configs_active
  ON plugin_device_configs(device_id, plugin_id)
  WHERE active = TRUE;

CREATE TABLE IF NOT EXISTS plugin_device_config_versions (
  id TEXT PRIMARY KEY,
  config_id TEXT NOT NULL REFERENCES plugin_device_configs(id),
  device_id TEXT NOT NULL REFERENCES devices(id),
  plugin_id TEXT NOT NULL,
  version INTEGER NOT NULL,
  config_json JSONB NOT NULL,
  change_note TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE plugin_device_config_versions IS '设备插件配置版本表，保存设备配置历史快照。';
COMMENT ON COLUMN plugin_device_config_versions.id IS '配置版本唯一标识。';
COMMENT ON COLUMN plugin_device_config_versions.config_id IS '所属设备插件配置 ID。';
COMMENT ON COLUMN plugin_device_config_versions.device_id IS '所属设备 ID。';
COMMENT ON COLUMN plugin_device_config_versions.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN plugin_device_config_versions.version IS '配置版本号。';
COMMENT ON COLUMN plugin_device_config_versions.config_json IS '配置版本快照 JSON。';
COMMENT ON COLUMN plugin_device_config_versions.change_note IS '版本变更说明。';
COMMENT ON COLUMN plugin_device_config_versions.created_at IS '版本创建时间。';

CREATE TABLE IF NOT EXISTS device_points (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL REFERENCES devices(id),
  plugin_id TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  address TEXT NOT NULL,
  value_type TEXT NOT NULL,
  unit TEXT NOT NULL DEFAULT '',
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  tags_json JSONB NOT NULL DEFAULT '{}'::JSONB,
  metadata_json JSONB NOT NULL DEFAULT '{}'::JSONB,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CHECK (value_type IN ('bool', 'int', 'float', 'string', 'bytes', 'datetime', 'json'))
);

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

CREATE INDEX IF NOT EXISTS idx_device_points_device
  ON device_points(device_id);

CREATE INDEX IF NOT EXISTS idx_device_points_plugin
  ON device_points(plugin_id);

CREATE TABLE IF NOT EXISTS device_point_versions (
  id TEXT PRIMARY KEY,
  point_id TEXT NOT NULL REFERENCES device_points(id),
  device_id TEXT NOT NULL REFERENCES devices(id),
  plugin_id TEXT NOT NULL,
  version INTEGER NOT NULL,
  snapshot_json JSONB NOT NULL,
  change_note TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE device_point_versions IS '点位版本表，保存通用点位表历史快照。';
COMMENT ON COLUMN device_point_versions.id IS '点位版本唯一标识。';
COMMENT ON COLUMN device_point_versions.point_id IS '所属点位 ID。';
COMMENT ON COLUMN device_point_versions.device_id IS '所属设备 ID。';
COMMENT ON COLUMN device_point_versions.plugin_id IS '所属插件 ID。';
COMMENT ON COLUMN device_point_versions.version IS '点位版本号。';
COMMENT ON COLUMN device_point_versions.snapshot_json IS '点位版本快照 JSON。';
COMMENT ON COLUMN device_point_versions.change_note IS '版本变更说明。';
COMMENT ON COLUMN device_point_versions.created_at IS '版本创建时间。';

CREATE TABLE IF NOT EXISTS collection_tasks (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  device_id TEXT NOT NULL REFERENCES devices(id),
  south_plugin_id TEXT NOT NULL,
  report_mode TEXT NOT NULL DEFAULT 'change',
  status TEXT NOT NULL DEFAULT 'stopped',
  rule_hit_rate TEXT NOT NULL DEFAULT '0%',
  rate TEXT NOT NULL DEFAULT '0 条/秒',
  last_collected_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW(),
  CHECK (report_mode IN ('change', 'all')),
  CHECK (status IN ('running', 'stopped', 'disabled', 'error'))
);

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

CREATE INDEX IF NOT EXISTS idx_collection_tasks_device
  ON collection_tasks(device_id);

CREATE INDEX IF NOT EXISTS idx_collection_tasks_plugin
  ON collection_tasks(south_plugin_id);

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
  created_at TIMESTAMPTZ DEFAULT NOW(),
  CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR'))
);

COMMENT ON TABLE runtime_events IS '运行事件表，保存低频运行事件、调试摘要和控制面日志。';
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

CREATE INDEX IF NOT EXISTS idx_runtime_events_time
  ON runtime_events(time DESC);

CREATE INDEX IF NOT EXISTS idx_runtime_events_device_time
  ON runtime_events(device_id, time DESC);
