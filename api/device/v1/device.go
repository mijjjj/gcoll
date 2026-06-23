package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

// DevicePluginConfigPageReq 描述设备插件配置页请求。
type DevicePluginConfigPageReq struct {
	g.Meta   `path:"/devices/{deviceId}/protocol-config" method:"get" tags:"Devices" summary:"获取设备协议配置页数据"`
	DeviceId string `json:"deviceId" in:"path"`
}

// DevicePluginConfigPageRes 描述设备插件配置页响应。
type DevicePluginConfigPageRes struct {
	Plugin           commonv1.PluginItem             `json:"plugin"`
	Device           commonv1.DeviceItem             `json:"device"`
	Config           map[string]any                  `json:"config"`
	ConfigSchema     map[string]any                  `json:"configSchema"`
	CustomConfigPage commonv1.PluginCustomConfigPage `json:"customConfigPage"`
	CustomPointPage  commonv1.PluginCustomPointPage  `json:"customPointPage"`
	Configured       bool                            `json:"configured"`
	Points           []commonv1.PointItem            `json:"points"`
	RecentEvents     []commonv1.RuntimeEvent         `json:"recentEvents"`
	Operations       []commonv1.PluginOperation      `json:"operations"`
	Warnings         []string                        `json:"warnings"`
}

// UpdateDevicePluginConfigReq 描述保存设备插件配置请求。
type UpdateDevicePluginConfigReq struct {
	g.Meta   `path:"/devices/{deviceId}/protocol-config" method:"put" tags:"Devices" summary:"保存设备插件配置"`
	DeviceId string         `json:"deviceId" in:"path"`
	Config   map[string]any `json:"config" v:"required#设备插件配置不能为空"`
}

// UpdateDevicePluginConfigRes 描述保存设备插件配置响应。
type UpdateDevicePluginConfigRes struct {
	Config map[string]any `json:"config"`
}

// TestDevicePluginConnectionReq 描述测试设备插件连接请求。
type TestDevicePluginConnectionReq struct {
	g.Meta   `path:"/devices/{deviceId}/protocol-config/test" method:"post" tags:"Devices" summary:"测试设备插件连接"`
	DeviceId string         `json:"deviceId" in:"path"`
	Config   map[string]any `json:"config"`
}

// TestDevicePluginConnectionRes 描述测试设备插件连接响应。
type TestDevicePluginConnectionRes struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	LatencyMs int    `json:"latencyMs"`
	TraceId   string `json:"traceId"`
}

// DevicesReq 描述设备列表请求。
type DevicesReq struct {
	g.Meta `path:"/devices" method:"get" tags:"Devices" summary:"获取设备列表"`
}

// DevicesRes 描述设备列表响应。
type DevicesRes struct {
	Groups []commonv1.DeviceGroup `json:"groups"`
	Items  []commonv1.DeviceItem  `json:"items"`
}

// CreateDeviceGroupReq 描述新增设备分组请求。
type CreateDeviceGroupReq struct {
	g.Meta `path:"/device-groups" method:"post" tags:"Devices" summary:"新增设备分组"`
	Id     string `json:"id"`
	Name   string `json:"name" v:"required#设备分组名称不能为空"`
}

// CreateDeviceGroupRes 描述新增设备分组响应。
type CreateDeviceGroupRes struct {
	Group commonv1.DeviceGroup `json:"group"`
}

// DeleteDeviceGroupReq 描述删除设备分组请求。
type DeleteDeviceGroupReq struct {
	g.Meta  `path:"/device-groups/{groupId}" method:"delete" tags:"Devices" summary:"删除设备分组"`
	GroupId string `json:"groupId" in:"path"`
}

// DeleteDeviceGroupRes 描述删除设备分组响应。
type DeleteDeviceGroupRes struct {
	GroupId string `json:"groupId"`
}

// CreateDeviceReq 描述新增设备请求。
type CreateDeviceReq struct {
	g.Meta      `path:"/devices" method:"post" tags:"Devices" summary:"新增设备"`
	Id          string         `json:"id"`
	Name        string         `json:"name" v:"required#设备名称不能为空"`
	GroupId     string         `json:"groupId" v:"required#设备分组不能为空"`
	PluginId    string         `json:"pluginId" v:"required#设备插件不能为空"`
	Enabled     bool           `json:"enabled"`
	ReportMode  string         `json:"reportMode" v:"required|in:change,all#上报模式不能为空|上报模式只能是 change 或 all"`
	Description string         `json:"description"`
	Config      map[string]any `json:"config"`
}

// CreateDeviceRes 描述新增设备响应。
type CreateDeviceRes struct {
	Device commonv1.DeviceItem `json:"device"`
}

// MoveDeviceToGroupReq 描述移动设备所属分组请求。
type MoveDeviceToGroupReq struct {
	g.Meta   `path:"/devices/{deviceId}/group" method:"patch" tags:"Devices" summary:"移动设备所属分组"`
	DeviceId string `json:"deviceId" in:"path"`
	GroupId  string `json:"groupId" v:"required#目标设备分组不能为空"`
}

// MoveDeviceToGroupRes 描述移动设备所属分组响应。
type MoveDeviceToGroupRes struct {
	Device commonv1.DeviceItem `json:"device"`
}

// DeleteDeviceReq 描述删除设备请求。
type DeleteDeviceReq struct {
	g.Meta   `path:"/devices/{deviceId}" method:"delete" tags:"Devices" summary:"删除设备"`
	DeviceId string `json:"deviceId" in:"path"`
}

// DeleteDeviceRes 描述删除设备响应。
type DeleteDeviceRes struct {
	DeviceId string `json:"deviceId"`
}

// DevicePointsReq 描述设备点位列表请求。
type DevicePointsReq struct {
	g.Meta   `path:"/devices/{deviceId}/points" method:"get" tags:"Devices" summary:"获取设备点位列表"`
	DeviceId string `json:"deviceId" in:"path"`
}

// DevicePointsRes 描述设备点位列表响应。
type DevicePointsRes struct {
	Items []commonv1.PointItem `json:"items"`
}

// CreateDevicePointReq 描述新增设备点位请求。
type CreateDevicePointReq struct {
	g.Meta      `path:"/devices/{deviceId}/points" method:"post" tags:"Devices" summary:"新增设备点位"`
	DeviceId    string            `json:"deviceId" in:"path"`
	Id          string            `json:"id"`
	PluginId    string            `json:"pluginId" v:"required#点位插件不能为空"`
	Name        string            `json:"name" v:"required#点位名称不能为空"`
	Description string            `json:"description"`
	Address     string            `json:"address" v:"required#点位地址不能为空"`
	ValueType   string            `json:"valueType" v:"required|in:bool,int,float,string,bytes,datetime,json#值类型不能为空|值类型不支持"`
	Unit        string            `json:"unit"`
	Enabled     bool              `json:"enabled"`
	Tags        map[string]string `json:"tags"`
	Metadata    map[string]any    `json:"metadata"`
}

// CreateDevicePointRes 描述新增设备点位响应。
type CreateDevicePointRes struct {
	Point commonv1.PointItem `json:"point"`
}

// UpdateDevicePointsReq 描述保存设备完整点位表请求。
type UpdateDevicePointsReq struct {
	g.Meta   `path:"/devices/{deviceId}/points" method:"put" tags:"Devices" summary:"保存设备完整点位表"`
	DeviceId string               `json:"deviceId" in:"path"`
	Items    []commonv1.PointItem `json:"items"`
}

// UpdateDevicePointsRes 描述保存设备完整点位表响应。
type UpdateDevicePointsRes struct {
	Items []commonv1.PointItem `json:"items"`
}
