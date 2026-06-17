// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package pointcache

import (
	"context"

	"github.com/mijjjj/gcoll/api/pointcache/v1"
)

type IPointcacheV1 interface {
	PointCache(ctx context.Context, req *v1.PointCacheReq) (res *v1.PointCacheRes, err error)
}
