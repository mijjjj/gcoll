package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) Devices(ctx context.Context, req *v1.DevicesReq) (res *v1.DevicesRes, err error) {
	_ = req

	return c.runtimeSvc.GetDevices(ctx)
}
