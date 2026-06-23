DROP INDEX IF EXISTS idx_devices_code;

ALTER TABLE devices DROP COLUMN code;
