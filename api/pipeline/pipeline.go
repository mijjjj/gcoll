// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package pipeline

import (
	"context"

	"github.com/mijjjj/gcoll/api/pipeline/v1"
)

type IPipelineV1 interface {
	PipelineRules(ctx context.Context, req *v1.PipelineRulesReq) (res *v1.PipelineRulesRes, err error)
}
