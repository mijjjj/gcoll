package task

import (
	"context"

	"github.com/mijjjj/gcoll/api/task/v1"
)

func (c *ControllerV1) Tasks(ctx context.Context, req *v1.TasksReq) (res *v1.TasksRes, err error) {
	_ = req
	return c.runtimeSvc.GetTasks(ctx)
}
