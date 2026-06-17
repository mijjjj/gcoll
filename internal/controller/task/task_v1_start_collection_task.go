package task

import (
	"context"

	"github.com/mijjjj/gcoll/api/task/v1"
)

func (c *ControllerV1) StartCollectionTask(ctx context.Context, req *v1.StartCollectionTaskReq) (res *v1.StartCollectionTaskRes, err error) {
	return c.runtimeSvc.StartCollectionTask(ctx, req.TaskId)
}
