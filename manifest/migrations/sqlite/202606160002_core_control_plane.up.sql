-- 控制面核心表：插件清单快照、设备、设备插件配置、通用点位、采集任务和运行事件。
CREATE TABLE IF NOT EXISTS schema_migrations (
  version TEXT PRIMARY KEY,
  applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
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
  enabled INTEGER NOT NULL DEFAULT 1,
  installed_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT,
  CHECK (type IN ('system', 'southbound', 'northbound')),
  CHECK (runtime IN ('process')),
  CHECK (protocol IN ('grpc')),
  CHECK (status IN ('installed', 'enabled', 'disabled', 'running', 'stopped', 'failed')),
  CHECK (enabled IN (0, 1))
);

CREATE TABLE IF NOT EXISTS plugin_versions (
  id TEXT PRIMARY KEY,
  plugin_id TEXT NOT NULL,
  version TEXT NOT NULL,
  package_path TEXT NOT NULL DEFAULT '',
  manifest_json TEXT NOT NULL,
  permissions_json TEXT NOT NULL DEFAULT '[]',
  capabilities_json TEXT NOT NULL DEFAULT '[]',
  checksum TEXT NOT NULL DEFAULT '',
  active INTEGER NOT NULL DEFAULT 0,
  installed_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT,
  CHECK (active IN (0, 1))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_plugin_versions_plugin_version
  ON plugin_versions(plugin_id, version)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS plugin_config_schemas (
  id TEXT PRIMARY KEY,
  plugin_id TEXT NOT NULL,
  plugin_version_id TEXT NOT NULL,
  schema_version INTEGER NOT NULL DEFAULT 1,
  schema_json TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT
);

CREATE INDEX IF NOT EXISTS idx_plugin_config_schemas_plugin
  ON plugin_config_schemas(plugin_id, plugin_version_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS device_groups (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  sort_order INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT
);

CREATE TABLE IF NOT EXISTS devices (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  code TEXT NOT NULL,
  group_id TEXT NOT NULL,
  plugin_id TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'offline',
  enabled INTEGER NOT NULL DEFAULT 0,
  report_mode TEXT NOT NULL DEFAULT 'change',
  last_seen_at TEXT,
  description TEXT NOT NULL DEFAULT '',
  metadata_json TEXT NOT NULL DEFAULT '{}',
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT,
  CHECK (status IN ('online', 'offline', 'disabled', 'error')),
  CHECK (enabled IN (0, 1)),
  CHECK (report_mode IN ('change', 'all')),
  FOREIGN KEY (group_id) REFERENCES device_groups(id)
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
  device_id TEXT NOT NULL,
  plugin_id TEXT NOT NULL,
  version INTEGER NOT NULL DEFAULT 1,
  config_json TEXT NOT NULL,
  report_mode TEXT NOT NULL DEFAULT 'change',
  enabled INTEGER NOT NULL DEFAULT 0,
  active INTEGER NOT NULL DEFAULT 1,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT,
  CHECK (report_mode IN ('change', 'all')),
  CHECK (enabled IN (0, 1)),
  CHECK (active IN (0, 1)),
  FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_plugin_device_configs_active
  ON plugin_device_configs(device_id, plugin_id, active)
  WHERE deleted_at IS NULL AND active = 1;

CREATE TABLE IF NOT EXISTS plugin_device_config_versions (
  id TEXT PRIMARY KEY,
  config_id TEXT NOT NULL,
  device_id TEXT NOT NULL,
  plugin_id TEXT NOT NULL,
  version INTEGER NOT NULL,
  config_json TEXT NOT NULL,
  change_note TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (config_id) REFERENCES plugin_device_configs(id),
  FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE TABLE IF NOT EXISTS device_points (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL,
  plugin_id TEXT NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  address TEXT NOT NULL,
  value_type TEXT NOT NULL,
  unit TEXT NOT NULL DEFAULT '',
  enabled INTEGER NOT NULL DEFAULT 1,
  tags_json TEXT NOT NULL DEFAULT '{}',
  metadata_json TEXT NOT NULL DEFAULT '{}',
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT,
  CHECK (value_type IN ('bool', 'int', 'float', 'string', 'bytes', 'datetime', 'json')),
  CHECK (enabled IN (0, 1)),
  FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE INDEX IF NOT EXISTS idx_device_points_device
  ON device_points(device_id)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_device_points_plugin
  ON device_points(plugin_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS device_point_versions (
  id TEXT PRIMARY KEY,
  point_id TEXT NOT NULL,
  device_id TEXT NOT NULL,
  plugin_id TEXT NOT NULL,
  version INTEGER NOT NULL,
  snapshot_json TEXT NOT NULL,
  change_note TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (point_id) REFERENCES device_points(id),
  FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE TABLE IF NOT EXISTS collection_tasks (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  device_id TEXT NOT NULL,
  south_plugin_id TEXT NOT NULL,
  report_mode TEXT NOT NULL DEFAULT 'change',
  status TEXT NOT NULL DEFAULT 'stopped',
  rule_hit_rate TEXT NOT NULL DEFAULT '0%',
  rate TEXT NOT NULL DEFAULT '0 条/秒',
  last_collected_at TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TEXT,
  CHECK (report_mode IN ('change', 'all')),
  CHECK (status IN ('running', 'stopped', 'disabled', 'error')),
  FOREIGN KEY (device_id) REFERENCES devices(id)
);

CREATE INDEX IF NOT EXISTS idx_collection_tasks_device
  ON collection_tasks(device_id)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_collection_tasks_plugin
  ON collection_tasks(south_plugin_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS runtime_events (
  id TEXT PRIMARY KEY,
  time TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  level TEXT NOT NULL,
  source TEXT NOT NULL,
  plugin_id TEXT,
  device_id TEXT,
  task_id TEXT,
  message TEXT NOT NULL,
  trace_id TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR'))
);

CREATE INDEX IF NOT EXISTS idx_runtime_events_time
  ON runtime_events(time DESC);

CREATE INDEX IF NOT EXISTS idx_runtime_events_device_time
  ON runtime_events(device_id, time DESC);

-- 表注释：schema_migrations：记录已经成功应用的数据库迁移版本。
-- 字段注释：schema_migrations.version：迁移版本号，通常来自迁移文件名。
-- 字段注释：schema_migrations.applied_at：迁移应用完成时间。
-- 表注释：plugins：插件清单快照表，保存导入或安装插件的控制面元数据；运行时事实来源仍是宿主内存注册表。
-- 字段注释：plugins.id：插件唯一标识。
-- 字段注释：plugins.name：插件显示名称。
-- 字段注释：plugins.type：插件类型，支持 system、southbound、northbound。
-- 字段注释：plugins.runtime：插件运行时类型。
-- 字段注释：plugins.protocol：插件通信协议。
-- 字段注释：plugins.status：插件当前管理状态。
-- 字段注释：plugins.active_version：当前活跃插件版本。
-- 字段注释：plugins.source：插件来源，例如 builtin 或 imported。
-- 字段注释：plugins.description：插件说明。
-- 字段注释：plugins.enabled：插件是否启用。
-- 字段注释：plugins.installed_at：插件安装时间。
-- 字段注释：plugins.updated_at：插件更新时间。
-- 字段注释：plugins.deleted_at：软删除时间。
-- 表注释：plugin_versions：插件版本表，保存每个插件版本的清单、权限和能力声明。
-- 字段注释：plugin_versions.id：插件版本唯一标识。
-- 字段注释：plugin_versions.plugin_id：所属插件 ID。
-- 字段注释：plugin_versions.version：插件版本号。
-- 字段注释：plugin_versions.package_path：插件包或插件目录路径。
-- 字段注释：plugin_versions.manifest_json：插件清单 JSON。
-- 字段注释：plugin_versions.permissions_json：插件权限声明 JSON。
-- 字段注释：plugin_versions.capabilities_json：插件能力声明 JSON。
-- 字段注释：plugin_versions.checksum：插件包校验和。
-- 字段注释：plugin_versions.active：该版本是否为活跃版本。
-- 字段注释：plugin_versions.installed_at：版本安装时间。
-- 字段注释：plugin_versions.created_at：记录创建时间。
-- 字段注释：plugin_versions.updated_at：记录更新时间。
-- 字段注释：plugin_versions.deleted_at：软删除时间。
-- 表注释：plugin_config_schemas：插件配置结构表，保存插件版本对应的配置表单和服务端校验结构。
-- 字段注释：plugin_config_schemas.id：配置结构唯一标识。
-- 字段注释：plugin_config_schemas.plugin_id：所属插件 ID。
-- 字段注释：plugin_config_schemas.plugin_version_id：所属插件版本 ID。
-- 字段注释：plugin_config_schemas.schema_version：配置结构版本号。
-- 字段注释：plugin_config_schemas.schema_json：配置结构 JSON。
-- 字段注释：plugin_config_schemas.created_at：记录创建时间。
-- 字段注释：plugin_config_schemas.deleted_at：软删除时间。
-- 表注释：device_groups：设备分组表，用于控制台设备列表归类。
-- 字段注释：device_groups.id：设备分组唯一标识。
-- 字段注释：device_groups.name：设备分组名称。
-- 字段注释：device_groups.sort_order：设备分组排序值。
-- 字段注释：device_groups.created_at：记录创建时间。
-- 字段注释：device_groups.updated_at：记录更新时间。
-- 字段注释：device_groups.deleted_at：软删除时间。
-- 表注释：devices：设备档案表，保存设备基础信息和绑定的南向插件 ID。
-- 字段注释：devices.id：设备唯一标识。
-- 字段注释：devices.name：设备名称。
-- 字段注释：devices.code：设备编码。
-- 字段注释：devices.group_id：所属设备分组 ID。
-- 字段注释：devices.plugin_id：设备绑定的南向插件 ID，由插件宿主注册表解释。
-- 字段注释：devices.status：设备状态。
-- 字段注释：devices.enabled：设备是否启用。
-- 字段注释：devices.report_mode：设备上报模式。
-- 字段注释：devices.last_seen_at：设备最近在线或采集时间。
-- 字段注释：devices.description：设备说明。
-- 字段注释：devices.metadata_json：设备扩展元数据 JSON。
-- 字段注释：devices.created_at：记录创建时间。
-- 字段注释：devices.updated_at：记录更新时间。
-- 字段注释：devices.deleted_at：软删除时间。
-- 表注释：plugin_device_configs：设备插件配置表，保存设备维度最终运行配置。
-- 字段注释：plugin_device_configs.id：设备插件配置唯一标识。
-- 字段注释：plugin_device_configs.device_id：所属设备 ID。
-- 字段注释：plugin_device_configs.plugin_id：所属插件 ID，由插件宿主注册表解释。
-- 字段注释：plugin_device_configs.version：当前配置版本号。
-- 字段注释：plugin_device_configs.config_json：设备插件配置 JSON，敏感值只保存密钥引用。
-- 字段注释：plugin_device_configs.report_mode：配置生效的上报模式。
-- 字段注释：plugin_device_configs.enabled：配置是否启用。
-- 字段注释：plugin_device_configs.active：是否为当前活跃配置。
-- 字段注释：plugin_device_configs.created_at：记录创建时间。
-- 字段注释：plugin_device_configs.updated_at：记录更新时间。
-- 字段注释：plugin_device_configs.deleted_at：软删除时间。
-- 表注释：plugin_device_config_versions：设备插件配置版本表，保存设备配置历史快照。
-- 字段注释：plugin_device_config_versions.id：配置版本唯一标识。
-- 字段注释：plugin_device_config_versions.config_id：所属设备插件配置 ID。
-- 字段注释：plugin_device_config_versions.device_id：所属设备 ID。
-- 字段注释：plugin_device_config_versions.plugin_id：所属插件 ID。
-- 字段注释：plugin_device_config_versions.version：配置版本号。
-- 字段注释：plugin_device_config_versions.config_json：配置版本快照 JSON。
-- 字段注释：plugin_device_config_versions.change_note：版本变更说明。
-- 字段注释：plugin_device_config_versions.created_at：版本创建时间。
-- 表注释：device_points：通用点位表，保存平台统一点位字段。
-- 字段注释：device_points.id：点位唯一标识。
-- 字段注释：device_points.device_id：所属设备 ID。
-- 字段注释：device_points.plugin_id：解释该点位的插件 ID。
-- 字段注释：device_points.name：点位名称。
-- 字段注释：device_points.description：点位说明。
-- 字段注释：device_points.address：点位地址，由插件解释。
-- 字段注释：device_points.value_type：平台通用值类型。
-- 字段注释：device_points.unit：点位单位。
-- 字段注释：device_points.enabled：点位是否启用。
-- 字段注释：device_points.tags_json：点位标签 JSON。
-- 字段注释：device_points.metadata_json：插件扩展元数据 JSON，不保存敏感值。
-- 字段注释：device_points.created_at：记录创建时间。
-- 字段注释：device_points.updated_at：记录更新时间。
-- 字段注释：device_points.deleted_at：软删除时间。
-- 表注释：device_point_versions：点位版本表，保存通用点位表历史快照。
-- 字段注释：device_point_versions.id：点位版本唯一标识。
-- 字段注释：device_point_versions.point_id：所属点位 ID。
-- 字段注释：device_point_versions.device_id：所属设备 ID。
-- 字段注释：device_point_versions.plugin_id：所属插件 ID。
-- 字段注释：device_point_versions.version：点位版本号。
-- 字段注释：device_point_versions.snapshot_json：点位版本快照 JSON。
-- 字段注释：device_point_versions.change_note：版本变更说明。
-- 字段注释：device_point_versions.created_at：版本创建时间。
-- 表注释：collection_tasks：采集任务表，保存设备采集任务的控制面状态。
-- 字段注释：collection_tasks.id：采集任务唯一标识。
-- 字段注释：collection_tasks.name：采集任务名称。
-- 字段注释：collection_tasks.device_id：采集任务关联设备 ID。
-- 字段注释：collection_tasks.south_plugin_id：采集任务使用的南向插件 ID。
-- 字段注释：collection_tasks.report_mode：采集任务上报模式。
-- 字段注释：collection_tasks.status：采集任务状态。
-- 字段注释：collection_tasks.rule_hit_rate：规则命中率展示值。
-- 字段注释：collection_tasks.rate：采集速率展示值。
-- 字段注释：collection_tasks.last_collected_at：最近采集时间。
-- 字段注释：collection_tasks.created_at：记录创建时间。
-- 字段注释：collection_tasks.updated_at：记录更新时间。
-- 字段注释：collection_tasks.deleted_at：软删除时间。
-- 表注释：runtime_events：运行事件表，保存低频运行事件、调试摘要和控制面日志。
-- 字段注释：runtime_events.id：运行事件唯一标识。
-- 字段注释：runtime_events.time：事件发生时间。
-- 字段注释：runtime_events.level：事件级别。
-- 字段注释：runtime_events.source：事件来源模块。
-- 字段注释：runtime_events.plugin_id：关联插件 ID。
-- 字段注释：runtime_events.device_id：关联设备 ID。
-- 字段注释：runtime_events.task_id：关联任务 ID。
-- 字段注释：runtime_events.message：事件消息。
-- 字段注释：runtime_events.trace_id：调用链追踪 ID。
-- 字段注释：runtime_events.created_at：记录创建时间。
