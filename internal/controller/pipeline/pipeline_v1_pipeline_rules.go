package pipeline

import (
	"context"

	"github.com/mijjjj/gcoll/api/pipeline/v1"
)

func (c *ControllerV1) PipelineRules(ctx context.Context, req *v1.PipelineRulesReq) (res *v1.PipelineRulesRes, err error) {
	_ = req
	return c.runtimeSvc.GetPipelineRules(ctx), nil
}
