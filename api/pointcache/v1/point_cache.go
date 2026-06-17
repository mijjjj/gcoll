package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

// PointCacheReq 描述最新点位缓存请求。
type PointCacheReq struct {
	g.Meta `path:"/point-cache" method:"get" tags:"PointCache" summary:"获取最新点位缓存"`
}

// PointCacheRes 描述最新点位缓存响应。
type PointCacheRes struct {
	Items []commonv1.PointCacheItem `json:"items"`
}
