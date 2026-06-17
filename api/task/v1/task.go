package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

// TasksReq 描述采集任务列表请求。
type TasksReq struct {
	g.Meta `path:"/tasks" method:"get" tags:"Tasks" summary:"获取采集任务列表"`
}

// TasksRes 描述采集任务列表响应。
type TasksRes struct {
	Items []commonv1.TaskSummary `json:"items"`
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
	Task commonv1.TaskSummary `json:"task"`
}

// StartDeviceCollectionTaskRes 描述启动设备默认采集任务响应。
type StartDeviceCollectionTaskRes = CollectionTaskActionRes

// StartCollectionTaskRes 描述启动采集任务响应。
type StartCollectionTaskRes = CollectionTaskActionRes

// StopCollectionTaskRes 描述停止采集任务响应。
type StopCollectionTaskRes = CollectionTaskActionRes
