package v1

import "github.com/gogf/gf/v2/frame/g"

// OverviewReq 描述工作台总览请求。
type OverviewReq struct {
	g.Meta `path:"/runtime/overview" method:"get" tags:"Runtime" summary:"获取运行时工作台总览"`
}

// OverviewRes 描述工作台总览响应。
type OverviewRes struct {
	Metrics       []MetricItem      `json:"metrics"`
	Runtime       RuntimeStatus     `json:"runtime"`
	DataPlane     DataPlaneStatus   `json:"dataPlane"`
	Tasks         []TaskSummary     `json:"tasks"`
	RecentEvents  []RuntimeEvent    `json:"recentEvents"`
	PluginSummary PluginSummary     `json:"pluginSummary"`
	Network       RuntimeDependency `json:"network"`
}

// PluginsReq 描述插件列表请求。
type PluginsReq struct {
	g.Meta `path:"/plugins" method:"get" tags:"Plugins" summary:"获取插件列表"`
}

// ImportPluginReq 描述插件导入请求。
type ImportPluginReq struct {
	g.Meta      `path:"/plugins/import" method:"post" tags:"Plugins" summary:"导入插件清单"`
	PackagePath string `json:"packagePath" v:"required#插件包路径不能为空"`
}

// ImportPluginRes 描述插件导入响应。
type ImportPluginRes struct {
	Plugin PluginItem `json:"plugin"`
}

// DevicePluginConfigPageReq 描述设备插件配置页请求。
type DevicePluginConfigPageReq struct {
	g.Meta   `path:"/devices/{deviceId}/protocol-config" method:"get" tags:"Devices" summary:"获取设备协议配置页数据"`
	DeviceId string `json:"deviceId" in:"path"`
}

// DevicePluginConfigPageRes 描述设备插件配置页响应。
type DevicePluginConfigPageRes struct {
	Plugin       PluginItem        `json:"plugin"`
	Device       DeviceItem        `json:"device"`
	Config       map[string]any    `json:"config"`
	ConfigSchema map[string]any    `json:"configSchema"`
	Configured   bool              `json:"configured"`
	Points       []PointItem       `json:"points"`
	RecentEvents []RuntimeEvent    `json:"recentEvents"`
	Operations   []PluginOperation `json:"operations"`
	Warnings     []string          `json:"warnings"`
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
	DeviceId string `json:"deviceId" in:"path"`
}

// TestDevicePluginConnectionRes 描述测试设备插件连接响应。
type TestDevicePluginConnectionRes struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	LatencyMs int    `json:"latencyMs"`
	TraceId   string `json:"traceId"`
}

// PluginsRes 描述插件列表响应。
type PluginsRes struct {
	Items []PluginItem `json:"items"`
}

// DevicesReq 描述设备列表请求。
type DevicesReq struct {
	g.Meta `path:"/devices" method:"get" tags:"Devices" summary:"获取设备列表"`
}

// CreateDeviceGroupReq 描述新增设备分组请求。
type CreateDeviceGroupReq struct {
	g.Meta `path:"/device-groups" method:"post" tags:"Devices" summary:"新增设备分组"`
	Id     string `json:"id"`
	Name   string `json:"name" v:"required#设备分组名称不能为空"`
}

// CreateDeviceGroupRes 描述新增设备分组响应。
type CreateDeviceGroupRes struct {
	Group DeviceGroup `json:"group"`
}

// CreateDeviceReq 描述新增设备请求。
type CreateDeviceReq struct {
	g.Meta      `path:"/devices" method:"post" tags:"Devices" summary:"新增设备"`
	Id          string         `json:"id"`
	Name        string         `json:"name" v:"required#设备名称不能为空"`
	Code        string         `json:"code" v:"required#设备编码不能为空"`
	GroupId     string         `json:"groupId" v:"required#设备分组不能为空"`
	PluginId    string         `json:"pluginId" v:"required#设备插件不能为空"`
	Enabled     bool           `json:"enabled"`
	ReportMode  string         `json:"reportMode" v:"required|in:change,all#上报模式不能为空|上报模式只能是 change 或 all"`
	Description string         `json:"description"`
	Config      map[string]any `json:"config"`
}

// CreateDeviceRes 描述新增设备响应。
type CreateDeviceRes struct {
	Device DeviceItem `json:"device"`
}

// MoveDeviceToGroupReq 描述移动设备所属分组请求。
type MoveDeviceToGroupReq struct {
	g.Meta   `path:"/devices/{deviceId}/group" method:"patch" tags:"Devices" summary:"移动设备所属分组"`
	DeviceId string `json:"deviceId" in:"path"`
	GroupId  string `json:"groupId" v:"required#目标设备分组不能为空"`
}

// MoveDeviceToGroupRes 描述移动设备所属分组响应。
type MoveDeviceToGroupRes struct {
	Device DeviceItem `json:"device"`
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

// DevicesRes 描述设备列表响应。
type DevicesRes struct {
	Groups []DeviceGroup `json:"groups"`
	Items  []DeviceItem  `json:"items"`
}

// DevicePointsReq 描述设备点位列表请求。
type DevicePointsReq struct {
	g.Meta   `path:"/devices/{deviceId}/points" method:"get" tags:"Devices" summary:"获取设备点位列表"`
	DeviceId string `json:"deviceId" in:"path"`
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
	Point PointItem `json:"point"`
}

// DevicePointsRes 描述设备点位列表响应。
type DevicePointsRes struct {
	Items []PointItem `json:"items"`
}

// TasksReq 描述采集任务列表请求。
type TasksReq struct {
	g.Meta `path:"/tasks" method:"get" tags:"Tasks" summary:"获取采集任务列表"`
}

// TasksRes 描述采集任务列表响应。
type TasksRes struct {
	Items []TaskSummary `json:"items"`
}

// StartDeviceCollectionTaskReq 描述启动设备默认采集任务请求。
type StartDeviceCollectionTaskReq struct {
	g.Meta   `path:"/devices/{deviceId}/tasks/start" method:"post" tags:"Tasks" summary:"启动设备采集任务"`
	DeviceId string `json:"deviceId" in:"path"`
}

// StartCollectionTaskReq 描述启动采集任务请求。
type StartCollectionTaskReq struct {
	g.Meta `path:"/tasks/{taskId}/start" method:"post" tags:"Tasks" summary:"启动采集任务"`
	TaskId string `json:"taskId" in:"path"`
}

// StopCollectionTaskReq 描述停止采集任务请求。
type StopCollectionTaskReq struct {
	g.Meta `path:"/tasks/{taskId}/stop" method:"post" tags:"Tasks" summary:"停止采集任务"`
	TaskId string `json:"taskId" in:"path"`
}

// CollectionTaskActionRes 描述采集任务启停响应。
type CollectionTaskActionRes struct {
	Task TaskSummary `json:"task"`
}

// PointCacheReq 描述最新点位缓存请求。
type PointCacheReq struct {
	g.Meta `path:"/point-cache" method:"get" tags:"PointCache" summary:"获取最新点位缓存"`
}

// PointCacheRes 描述最新点位缓存响应。
type PointCacheRes struct {
	Items []PointCacheItem `json:"items"`
}

// PipelineRulesReq 描述规则过滤列表请求。
type PipelineRulesReq struct {
	g.Meta `path:"/pipeline/rules" method:"get" tags:"Pipeline" summary:"获取规则过滤列表"`
}

// PipelineRulesRes 描述规则过滤列表响应。
type PipelineRulesRes struct {
	Items []PipelineRuleItem `json:"items"`
}

// TargetsReq 描述北向转发目标列表请求。
type TargetsReq struct {
	g.Meta `path:"/targets" method:"get" tags:"Targets" summary:"获取北向转发目标列表"`
}

// TargetsRes 描述北向转发目标列表响应。
type TargetsRes struct {
	Items []ForwardTargetItem `json:"items"`
}

// LogsReq 描述运行日志列表请求。
type LogsReq struct {
	g.Meta `path:"/logs" method:"get" tags:"Logs" summary:"获取运行日志列表"`
}

// LogsRes 描述运行日志列表响应。
type LogsRes struct {
	Items []RuntimeEvent `json:"items"`
}

// MetricItem 描述工作台指标。
type MetricItem struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
	Hint  string `json:"hint"`
	Tone  string `json:"tone"`
}

// RuntimeStatus 描述核心运行时状态。
type RuntimeStatus struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Mode      string `json:"mode"`
	CheckedAt string `json:"checkedAt"`
	ApiBase   string `json:"apiBase"`
	Database  string `json:"database"`
}

// RuntimeDependency 描述运行时依赖状态。
type RuntimeDependency struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Detail string `json:"detail"`
}

// DataPlaneStatus 描述数据面状态。
type DataPlaneStatus struct {
	QueueUsagePercent int    `json:"queueUsagePercent"`
	RuleHitPercent    int    `json:"ruleHitPercent"`
	ForwardPercent    int    `json:"forwardPercent"`
	Throughput        string `json:"throughput"`
	Latency           string `json:"latency"`
	Backpressure      string `json:"backpressure"`
}

// PluginSummary 描述插件运行汇总。
type PluginSummary struct {
	Running int `json:"running"`
	Total   int `json:"total"`
}

// PluginItem 描述插件列表项。
type PluginItem struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Version     string   `json:"version"`
	Runtime     string   `json:"runtime"`
	Protocol    string   `json:"protocol"`
	Status      string   `json:"status"`
	Permissions []string `json:"permissions"`
	UpdatedAt   string   `json:"updatedAt"`
}

// DeviceGroup 描述设备分组。
type DeviceGroup struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// DeviceItem 描述设备列表项。
type DeviceItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	GroupId     string `json:"groupId"`
	PluginId    string `json:"pluginId"`
	PluginName  string `json:"pluginName"`
	Status      string `json:"status"`
	Enabled     bool   `json:"enabled"`
	PointCount  int    `json:"pointCount"`
	ReportMode  string `json:"reportMode"`
	LastSeenAt  string `json:"lastSeenAt"`
	Description string `json:"description"`
}

// PointItem 描述通用点位表列表项。
type PointItem struct {
	Id          string            `json:"id"`
	DeviceId    string            `json:"deviceId"`
	PluginId    string            `json:"pluginId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Address     string            `json:"address"`
	ValueType   string            `json:"valueType"`
	Unit        string            `json:"unit"`
	Enabled     bool              `json:"enabled"`
	Tags        map[string]string `json:"tags"`
	Metadata    map[string]any    `json:"metadata,omitempty"`
}

// TaskSummary 描述采集任务摘要。
type TaskSummary struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	DeviceId        string `json:"deviceId"`
	DeviceName      string `json:"deviceName"`
	SouthPluginName string `json:"southPluginName"`
	PointCount      int    `json:"pointCount"`
	ReportMode      string `json:"reportMode"`
	Status          string `json:"status"`
	Rate            string `json:"rate"`
	RuleHitRate     string `json:"ruleHitRate"`
	LastCollectedAt string `json:"lastCollectedAt"`
}

// PointCacheItem 描述最新点位缓存项。
type PointCacheItem struct {
	PointId   string `json:"pointId"`
	DeviceId  string `json:"deviceId"`
	PointName string `json:"pointName"`
	Value     string `json:"value"`
	Quality   string `json:"quality"`
	Changed   bool   `json:"changed"`
	UpdatedAt string `json:"updatedAt"`
}

// PipelineRuleItem 描述规则过滤项。
type PipelineRuleItem struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Expression  string `json:"expression"`
	Matched     int    `json:"matched"`
	TargetCount int    `json:"targetCount"`
	UpdatedAt   string `json:"updatedAt"`
}

// ForwardTargetItem 描述北向转发目标。
type ForwardTargetItem struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PluginName string `json:"pluginName"`
	Status     string `json:"status"`
	Endpoint   string `json:"endpoint"`
	LastError  string `json:"lastError"`
	UpdatedAt  string `json:"updatedAt"`
}

// RuntimeEvent 描述运行时事件或日志。
type RuntimeEvent struct {
	Id       string `json:"id"`
	Time     string `json:"time"`
	Level    string `json:"level"`
	Source   string `json:"source"`
	PluginId string `json:"pluginId"`
	DeviceId string `json:"deviceId"`
	TaskId   string `json:"taskId"`
	Message  string `json:"message"`
	TraceId  string `json:"traceId"`
}

// PluginOperation 描述设备插件配置页可执行动作。
type PluginOperation struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
