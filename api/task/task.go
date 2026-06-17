// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package task

import (
	"context"

	"github.com/mijjjj/gcoll/api/task/v1"
)

type ITaskV1 interface {
	Tasks(ctx context.Context, req *v1.TasksReq) (res *v1.TasksRes, err error)
	StartDeviceCollectionTask(ctx context.Context, req *v1.StartDeviceCollectionTaskReq) (res *v1.StartDeviceCollectionTaskRes, err error)
	StartCollectionTask(ctx context.Context, req *v1.StartCollectionTaskReq) (res *v1.StartCollectionTaskRes, err error)
	StopCollectionTask(ctx context.Context, req *v1.StopCollectionTaskReq) (res *v1.StopCollectionTaskRes, err error)
}
