UPDATE plugins
SET updated_at = NOW()
WHERE updated_at IS NULL;

ALTER TABLE plugins
  ALTER COLUMN updated_at SET NOT NULL;

COMMENT ON COLUMN plugins.updated_at IS '插件更新时间。';

UPDATE plugin_versions
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

ALTER TABLE plugin_versions
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET NOT NULL;

COMMENT ON COLUMN plugin_versions.created_at IS '记录创建时间。';
COMMENT ON COLUMN plugin_versions.updated_at IS '记录更新时间。';

UPDATE plugin_config_schemas
SET created_at = NOW()
WHERE created_at IS NULL;

ALTER TABLE plugin_config_schemas
  ALTER COLUMN created_at SET NOT NULL;

COMMENT ON COLUMN plugin_config_schemas.created_at IS '记录创建时间。';

UPDATE device_groups
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

ALTER TABLE device_groups
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET NOT NULL;

COMMENT ON COLUMN device_groups.created_at IS '记录创建时间。';
COMMENT ON COLUMN device_groups.updated_at IS '记录更新时间。';

UPDATE devices
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

ALTER TABLE devices
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET NOT NULL;

COMMENT ON COLUMN devices.created_at IS '记录创建时间。';
COMMENT ON COLUMN devices.updated_at IS '记录更新时间。';

UPDATE plugin_device_configs
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

ALTER TABLE plugin_device_configs
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET NOT NULL;

COMMENT ON COLUMN plugin_device_configs.created_at IS '记录创建时间。';
COMMENT ON COLUMN plugin_device_configs.updated_at IS '记录更新时间。';

UPDATE plugin_device_config_versions
SET created_at = NOW()
WHERE created_at IS NULL;

ALTER TABLE plugin_device_config_versions
  ALTER COLUMN created_at SET NOT NULL;

COMMENT ON COLUMN plugin_device_config_versions.created_at IS '版本创建时间。';

UPDATE device_points
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

ALTER TABLE device_points
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET NOT NULL;

COMMENT ON COLUMN device_points.created_at IS '记录创建时间。';
COMMENT ON COLUMN device_points.updated_at IS '记录更新时间。';

UPDATE device_point_versions
SET created_at = NOW()
WHERE created_at IS NULL;

ALTER TABLE device_point_versions
  ALTER COLUMN created_at SET NOT NULL;

COMMENT ON COLUMN device_point_versions.created_at IS '版本创建时间。';

UPDATE collection_tasks
SET created_at = COALESCE(created_at, NOW()),
    updated_at = COALESCE(updated_at, NOW())
WHERE created_at IS NULL OR updated_at IS NULL;

ALTER TABLE collection_tasks
  ALTER COLUMN created_at SET NOT NULL,
  ALTER COLUMN updated_at SET NOT NULL;

COMMENT ON COLUMN collection_tasks.created_at IS '记录创建时间。';
COMMENT ON COLUMN collection_tasks.updated_at IS '记录更新时间。';

UPDATE runtime_events
SET created_at = NOW()
WHERE created_at IS NULL;

ALTER TABLE runtime_events
  ALTER COLUMN created_at SET NOT NULL;

COMMENT ON COLUMN runtime_events.created_at IS '记录创建时间。';
