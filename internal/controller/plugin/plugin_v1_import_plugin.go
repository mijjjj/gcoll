package plugin

import (
	"context"

	"github.com/mijjjj/gcoll/api/plugin/v1"
)

func (c *ControllerV1) ImportPlugin(ctx context.Context, req *v1.ImportPluginReq) (res *v1.ImportPluginRes, err error) {
	return c.runtimeSvc.ImportPlugin(ctx, req.PackagePath)
}
