package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) DevicePluginConfigPage(ctx context.Context, req *v1.DevicePluginConfigPageReq) (res *v1.DevicePluginConfigPageRes, err error) {
	return c.runtimeSvc.GetDevicePluginConfigPage(ctx, req.DeviceId)
}
