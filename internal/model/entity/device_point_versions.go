// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// DevicePointVersions is the golang structure for table device_point_versions.
type DevicePointVersions struct {
	Id           string `json:"id"           orm:"id"            ` //
	PointId      string `json:"pointId"      orm:"point_id"      ` //
	DeviceId     string `json:"deviceId"     orm:"device_id"     ` //
	PluginId     string `json:"pluginId"     orm:"plugin_id"     ` //
	Version      int    `json:"version"      orm:"version"       ` //
	SnapshotJson string `json:"snapshotJson" orm:"snapshot_json" ` //
	ChangeNote   string `json:"changeNote"   orm:"change_note"   ` //
	CreatedAt    string `json:"createdAt"    orm:"created_at"    ` //
}
