package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	commonv1 "github.com/mijjjj/gcoll/api/common/v1"
)

// PipelineRulesReq 描述规则过滤列表请求。
type PipelineRulesReq struct {
	g.Meta `path:"/pipeline/rules" method:"get" tags:"Pipeline" summary:"获取规则过滤列表"`
}

// PipelineRulesRes 描述规则过滤列表响应。
type PipelineRulesRes struct {
	Items []commonv1.PipelineRuleItem `json:"items"`
}
