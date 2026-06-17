package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) CreateDevice(ctx context.Context, req *v1.CreateDeviceReq) (res *v1.CreateDeviceRes, err error) {
	return c.runtimeSvc.CreateDevice(ctx, req)
}
