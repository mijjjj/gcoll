// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// RuntimeEvents is the golang structure for table runtime_events.
type RuntimeEvents struct {
	Id        string `json:"id"        orm:"id"         ` //
	Time      string `json:"time"      orm:"time"       ` //
	Level     string `json:"level"     orm:"level"      ` //
	Source    string `json:"source"    orm:"source"     ` //
	PluginId  string `json:"pluginId"  orm:"plugin_id"  ` //
	DeviceId  string `json:"deviceId"  orm:"device_id"  ` //
	TaskId    string `json:"taskId"    orm:"task_id"    ` //
	Message   string `json:"message"   orm:"message"    ` //
	TraceId   string `json:"traceId"   orm:"trace_id"   ` //
	CreatedAt string `json:"createdAt" orm:"created_at" ` //
}
