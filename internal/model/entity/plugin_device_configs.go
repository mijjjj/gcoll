// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PluginDeviceConfigs is the golang structure for table plugin_device_configs.
type PluginDeviceConfigs struct {
	Id         string `json:"id"         orm:"id"          ` //
	DeviceId   string `json:"deviceId"   orm:"device_id"   ` //
	PluginId   string `json:"pluginId"   orm:"plugin_id"   ` //
	Version    int    `json:"version"    orm:"version"     ` //
	ConfigJson string `json:"configJson" orm:"config_json" ` //
	ReportMode string `json:"reportMode" orm:"report_mode" ` //
	Enabled    int    `json:"enabled"    orm:"enabled"     ` //
	Active     int    `json:"active"     orm:"active"      ` //
	CreatedAt  string `json:"createdAt"  orm:"created_at"  ` //
	UpdatedAt  string `json:"updatedAt"  orm:"updated_at"  ` //
	DeletedAt  string `json:"deletedAt"  orm:"deleted_at"  ` //
}
