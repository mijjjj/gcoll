package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) DevicePoints(ctx context.Context, req *v1.DevicePointsReq) (res *v1.DevicePointsRes, err error) {
	return c.runtimeSvc.GetDevicePoints(ctx, req.DeviceId)
}
