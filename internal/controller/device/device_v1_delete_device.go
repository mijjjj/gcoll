package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) DeleteDevice(ctx context.Context, req *v1.DeleteDeviceReq) (res *v1.DeleteDeviceRes, err error) {
	return c.runtimeSvc.DeleteDevice(ctx, req.DeviceId)
}
