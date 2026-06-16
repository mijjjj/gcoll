// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// PluginDeviceConfigVersions is the golang structure for table plugin_device_config_versions.
type PluginDeviceConfigVersions struct {
	Id         string `json:"id"         orm:"id"          ` //
	ConfigId   string `json:"configId"   orm:"config_id"   ` //
	DeviceId   string `json:"deviceId"   orm:"device_id"   ` //
	PluginId   string `json:"pluginId"   orm:"plugin_id"   ` //
	Version    int    `json:"version"    orm:"version"     ` //
	ConfigJson string `json:"configJson" orm:"config_json" ` //
	ChangeNote string `json:"changeNote" orm:"change_note" ` //
	CreatedAt  string `json:"createdAt"  orm:"created_at"  ` //
}
