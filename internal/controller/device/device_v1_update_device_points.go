package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) UpdateDevicePoints(ctx context.Context, req *v1.UpdateDevicePointsReq) (res *v1.UpdateDevicePointsRes, err error) {
	return c.runtimeSvc.UpdateDevicePoints(ctx, req)
}
