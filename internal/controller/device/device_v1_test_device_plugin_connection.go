package device

import (
	"context"

	"github.com/mijjjj/gcoll/api/device/v1"
)

func (c *ControllerV1) TestDevicePluginConnection(ctx context.Context, req *v1.TestDevicePluginConnectionReq) (res *v1.TestDevicePluginConnectionRes, err error) {
	return c.runtimeSvc.TestDevicePluginConnection(ctx, req)
}
