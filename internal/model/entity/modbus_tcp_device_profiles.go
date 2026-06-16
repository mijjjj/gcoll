// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// ModbusTcpDeviceProfiles is the golang structure for table modbus_tcp_device_profiles.
type ModbusTcpDeviceProfiles struct {
	Id               string `json:"id"               orm:"id"                 ` //
	DeviceId         string `json:"deviceId"         orm:"device_id"          ` //
	PluginId         string `json:"pluginId"         orm:"plugin_id"          ` //
	Version          int    `json:"version"          orm:"version"            ` //
	Host             string `json:"host"             orm:"host"               ` //
	Port             int    `json:"port"             orm:"port"               ` //
	UnitId           int    `json:"unitId"           orm:"unit_id"            ` //
	TimeoutMs        int    `json:"timeoutMs"        orm:"timeout_ms"         ` //
	PollIntervalMs   int    `json:"pollIntervalMs"   orm:"poll_interval_ms"   ` //
	ReportMode       string `json:"reportMode"       orm:"report_mode"        ` //
	MaxCoilBatch     int    `json:"maxCoilBatch"     orm:"max_coil_batch"     ` //
	MaxRegisterBatch int    `json:"maxRegisterBatch" orm:"max_register_batch" ` //
	LowLatencyMs     int    `json:"lowLatencyMs"     orm:"low_latency_ms"     ` //
	HighLatencyMs    int    `json:"highLatencyMs"    orm:"high_latency_ms"    ` //
	DebugEnabled     int    `json:"debugEnabled"     orm:"debug_enabled"      ` //
	Enabled          int    `json:"enabled"          orm:"enabled"            ` //
	CreatedAt        string `json:"createdAt"        orm:"created_at"         ` //
	UpdatedAt        string `json:"updatedAt"        orm:"updated_at"         ` //
	DeletedAt        string `json:"deletedAt"        orm:"deleted_at"         ` //
}
