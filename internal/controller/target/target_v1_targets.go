package target

import (
	"context"

	"github.com/mijjjj/gcoll/api/target/v1"
)

func (c *ControllerV1) Targets(ctx context.Context, req *v1.TargetsReq) (res *v1.TargetsRes, err error) {
	_ = req
	return c.runtimeSvc.GetTargets(ctx), nil
}
