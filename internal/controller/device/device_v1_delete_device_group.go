package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) DeleteDeviceGroup(ctx context.Context, req *v1.DeleteDeviceGroupReq) (res *v1.DeleteDeviceGroupRes, err error) {
	return c.runtimeSvc.DeleteDeviceGroup(ctx, req.GroupId)
}
