package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) UpdateDevicePluginConfig(ctx context.Context, req *v1.UpdateDevicePluginConfigReq) (res *v1.UpdateDevicePluginConfigRes, err error) {
	return c.runtimeSvc.UpdateDevicePluginConfig(ctx, req)
}
