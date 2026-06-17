package runtime

import (
	"context"

	"github.com/mijjjj/gcoll/api/runtime/v1"
)

func (c *ControllerV1) Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error) {
	_ = req
	return c.runtimeSvc.GetHealth(ctx), nil
}
