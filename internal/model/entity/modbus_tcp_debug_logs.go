// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// ModbusTcpDebugLogs is the golang structure for table modbus_tcp_debug_logs.
type ModbusTcpDebugLogs struct {
	Id         string `json:"id"         orm:"id"          ` //
	DeviceId   string `json:"deviceId"   orm:"device_id"   ` //
	TaskId     string `json:"taskId"     orm:"task_id"     ` //
	PointId    string `json:"pointId"    orm:"point_id"    ` //
	TraceId    string `json:"traceId"    orm:"trace_id"    ` //
	Level      string `json:"level"      orm:"level"       ` //
	Message    string `json:"message"    orm:"message"     ` //
	Area       string `json:"area"       orm:"area"        ` //
	Address    int    `json:"address"    orm:"address"     ` //
	LatencyMs  int    `json:"latencyMs"  orm:"latency_ms"  ` //
	RawHex     string `json:"rawHex"     orm:"raw_hex"     ` //
	FieldsJson string `json:"fieldsJson" orm:"fields_json" ` //
	CreatedAt  string `json:"createdAt"  orm:"created_at"  ` //
}
