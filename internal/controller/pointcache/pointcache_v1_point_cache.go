package pointcache

import (
	"context"

	"github.com/mijjjj/gcoll/api/pointcache/v1"
)

func (c *ControllerV1) PointCache(ctx context.Context, req *v1.PointCacheReq) (res *v1.PointCacheRes, err error) {
	_ = req
	return c.runtimeSvc.GetPointCache(ctx), nil
}
