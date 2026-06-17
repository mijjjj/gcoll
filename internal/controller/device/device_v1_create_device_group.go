package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) CreateDeviceGroup(ctx context.Context, req *v1.CreateDeviceGroupReq) (res *v1.CreateDeviceGroupRes, err error) {
	return c.runtimeSvc.CreateDeviceGroup(ctx, req)
}
