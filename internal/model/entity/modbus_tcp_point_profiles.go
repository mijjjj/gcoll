// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// ModbusTcpPointProfiles is the golang structure for table modbus_tcp_point_profiles.
type ModbusTcpPointProfiles struct {
	Id         string  `json:"id"         orm:"id"          ` //
	DeviceId   string  `json:"deviceId"   orm:"device_id"   ` //
	PointId    string  `json:"pointId"    orm:"point_id"    ` //
	PluginId   string  `json:"pluginId"   orm:"plugin_id"   ` //
	Version    int     `json:"version"    orm:"version"     ` //
	Area       string  `json:"area"       orm:"area"        ` //
	Address    int     `json:"address"    orm:"address"     ` //
	Quantity   int     `json:"quantity"   orm:"quantity"    ` //
	Mode       string  `json:"mode"       orm:"mode"        ` //
	ValueType  string  `json:"valueType"  orm:"value_type"  ` //
	ByteOrder  string  `json:"byteOrder"  orm:"byte_order"  ` //
	WordOrder  string  `json:"wordOrder"  orm:"word_order"  ` //
	Scale      float32 `json:"scale"      orm:"scale"       ` //
	Offset     float32 `json:"offset"     orm:"offset"      ` //
	ReportMode string  `json:"reportMode" orm:"report_mode" ` //
	Enabled    int     `json:"enabled"    orm:"enabled"     ` //
	CreatedAt  string  `json:"createdAt"  orm:"created_at"  ` //
	UpdatedAt  string  `json:"updatedAt"  orm:"updated_at"  ` //
	DeletedAt  string  `json:"deletedAt"  orm:"deleted_at"  ` //
}
