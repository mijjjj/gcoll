DROP INDEX IF EXISTS idx_runtime_events_device_time;
DROP INDEX IF EXISTS idx_runtime_events_time;
DROP TABLE IF EXISTS runtime_events;

DROP INDEX IF EXISTS idx_collection_tasks_plugin;
DROP INDEX IF EXISTS idx_collection_tasks_device;
DROP TABLE IF EXISTS collection_tasks;

DROP TABLE IF EXISTS device_point_versions;

DROP INDEX IF EXISTS idx_device_points_plugin;
DROP INDEX IF EXISTS idx_device_points_device;
DROP TABLE IF EXISTS device_points;

DROP TABLE IF EXISTS plugin_device_config_versions;

DROP INDEX IF EXISTS idx_plugin_device_configs_active;
DROP TABLE IF EXISTS plugin_device_configs;

DROP INDEX IF EXISTS idx_devices_plugin;
DROP INDEX IF EXISTS idx_devices_group;
DROP INDEX IF EXISTS idx_devices_code;
DROP TABLE IF EXISTS devices;

DROP TABLE IF EXISTS device_groups;

DROP INDEX IF EXISTS idx_plugin_config_schemas_plugin;
DROP TABLE IF EXISTS plugin_config_schemas;

DROP INDEX IF EXISTS idx_plugin_versions_plugin_version;
DROP TABLE IF EXISTS plugin_versions;

DROP TABLE IF EXISTS plugins;
