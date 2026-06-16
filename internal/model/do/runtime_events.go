// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// RuntimeEvents is the golang structure of table runtime_events for DAO operations like Where/Data.
type RuntimeEvents struct {
	g.Meta    `orm:"table:runtime_events, do:true"`
	Id        any //
	Time      any //
	Level     any //
	Source    any //
	PluginId  any //
	DeviceId  any //
	TaskId    any //
	Message   any //
	TraceId   any //
	CreatedAt any //
}
