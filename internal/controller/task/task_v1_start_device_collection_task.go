package task

import (
	"context"

	"github.com/mijjjj/gcoll/api/task/v1"
)

func (c *ControllerV1) StartDeviceCollectionTask(ctx context.Context, req *v1.StartDeviceCollectionTaskReq) (res *v1.StartDeviceCollectionTaskRes, err error) {
	return c.runtimeSvc.StartDeviceCollectionTask(ctx, req.DeviceId)
}
