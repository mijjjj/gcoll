package v1

import "github.com/gogf/gf/v2/frame/g"

// HealthReq 描述运行时健康检查请求。
type HealthReq struct {
	g.Meta `path:"/runtime/health" method:"get" tags:"Runtime" summary:"获取运行时健康状态"`
}

// HealthRes 描述运行时健康检查响应。
type HealthRes struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Mode      string `json:"mode"`
	CheckedAt string `json:"checkedAt"`
}
