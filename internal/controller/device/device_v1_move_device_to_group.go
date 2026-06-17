package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) MoveDeviceToGroup(ctx context.Context, req *v1.MoveDeviceToGroupReq) (res *v1.MoveDeviceToGroupRes, err error) {
	return c.runtimeSvc.MoveDeviceToGroup(ctx, req)
}
