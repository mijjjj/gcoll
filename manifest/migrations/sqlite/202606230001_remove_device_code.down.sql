ALTER TABLE devices ADD COLUMN code TEXT;

UPDATE devices
SET code = id
WHERE code IS NULL OR code = '';

CREATE UNIQUE INDEX IF NOT EXISTS idx_devices_code
  ON devices(code);
