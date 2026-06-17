package task

import (
	"context"

	"github.com/mijjjj/gcoll/api/task/v1"
)

func (c *ControllerV1) StopCollectionTask(ctx context.Context, req *v1.StopCollectionTaskReq) (res *v1.StopCollectionTaskRes, err error) {
	return c.runtimeSvc.StopCollectionTask(ctx, req.TaskId)
}
