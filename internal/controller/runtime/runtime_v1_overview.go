package runtime

import (
	"context"

	"github.com/mijjjj/gcoll/api/runtime/v1"
)

func (c *ControllerV1) Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error) {
	_ = req
	return c.runtimeSvc.GetOverview(ctx)
}
