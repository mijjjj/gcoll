package plugin

import (
	"context"

	"github.com/mijjjj/gcoll/api/plugin/v1"
)

func (c *ControllerV1) Plugins(ctx context.Context, req *v1.PluginsReq) (res *v1.PluginsRes, err error) {
	_ = req
	return c.runtimeSvc.GetPlugins(ctx)
}
