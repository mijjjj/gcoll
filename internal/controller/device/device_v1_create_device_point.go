package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) CreateDevicePoint(ctx context.Context, req *v1.CreateDevicePointReq) (res *v1.CreateDevicePointRes, err error) {
	return c.runtimeSvc.CreateDevicePoint(ctx, req)
}
