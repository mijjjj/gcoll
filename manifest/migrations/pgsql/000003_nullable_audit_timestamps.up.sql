ALTER TABLE plugins
  ALTER COLUMN updated_at DROP NOT NULL;

COMMENT ON COLUMN plugins.updated_at IS '插件更新时间，允许为空。';

ALTER TABLE plugin_versions
  ALTER COLUMN created_at DROP NOT NULL,
  ALTER COLUMN updated_at DROP NOT NULL;

COMMENT ON COLUMN plugin_versions.created_at IS '记录创建时间，允许为空。';
COMMENT ON COLUMN plugin_versions.updated_at IS '记录更新时间，允许为空。';

ALTER TABLE plugin_config_schemas
  ALTER COLUMN created_at DROP NOT NULL;

COMMENT ON COLUMN plugin_config_schemas.created_at IS '记录创建时间，允许为空。';

ALTER TABLE device_groups
  ALTER COLUMN created_at DROP NOT NULL,
  ALTER COLUMN updated_at DROP NOT NULL;

COMMENT ON COLUMN device_groups.created_at IS '记录创建时间，允许为空。';
COMMENT ON COLUMN device_groups.updated_at IS '记录更新时间，允许为空。';

ALTER TABLE devices
  ALTER COLUMN created_at DROP NOT NULL,
  ALTER COLUMN updated_at DROP NOT NULL;

COMMENT ON COLUMN devices.created_at IS '记录创建时间，允许为空。';
COMMENT ON COLUMN devices.updated_at IS '记录更新时间，允许为空。';

ALTER TABLE plugin_device_configs
  ALTER COLUMN created_at DROP NOT NULL,
  ALTER COLUMN updated_at DROP NOT NULL;

COMMENT ON COLUMN plugin_device_configs.created_at IS '记录创建时间，允许为空。';
COMMENT ON COLUMN plugin_device_configs.updated_at IS '记录更新时间，允许为空。';

ALTER TABLE plugin_device_config_versions
  ALTER COLUMN created_at DROP NOT NULL;

COMMENT ON COLUMN plugin_device_config_versions.created_at IS '版本创建时间，允许为空。';

ALTER TABLE device_points
  ALTER COLUMN created_at DROP NOT NULL,
  ALTER COLUMN updated_at DROP NOT NULL;

COMMENT ON COLUMN device_points.created_at IS '记录创建时间，允许为空。';
COMMENT ON COLUMN device_points.updated_at IS '记录更新时间，允许为空。';

ALTER TABLE device_point_versions
  ALTER COLUMN created_at DROP NOT NULL;

COMMENT ON COLUMN device_point_versions.created_at IS '版本创建时间，允许为空。';

ALTER TABLE collection_tasks
  ALTER COLUMN created_at DROP NOT NULL,
  ALTER COLUMN updated_at DROP NOT NULL;

COMMENT ON COLUMN collection_tasks.created_at IS '记录创建时间，允许为空。';
COMMENT ON COLUMN collection_tasks.updated_at IS '记录更新时间，允许为空。';

ALTER TABLE runtime_events
  ALTER COLUMN created_at DROP NOT NULL;

COMMENT ON COLUMN runtime_events.created_at IS '记录创建时间，允许为空。';
