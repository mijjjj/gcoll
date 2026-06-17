package v1

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

// PluginCustomConfigPage 描述插件自带设备配置页资源。
type PluginCustomConfigPage struct {
	Enabled bool   `json:"enabled"`
	Entry   string `json:"entry"`
	Script  string `json:"script"`
	Html    string `json:"html"`
	Js      string `json:"js"`
}

// PluginCustomPointPage 描述插件自带点位配置页资源。
type PluginCustomPointPage struct {
	Enabled bool   `json:"enabled"`
	Entry   string `json:"entry"`
	Script  string `json:"script"`
	Html    string `json:"html"`
	Js      string `json:"js"`
}
