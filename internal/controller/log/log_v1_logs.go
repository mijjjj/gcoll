package log

import (
	"context"

	"github.com/mijjjj/gcoll/api/log/v1"
)

func (c *ControllerV1) Logs(ctx context.Context, req *v1.LogsReq) (res *v1.LogsRes, err error) {
	_ = req
	return c.runtimeSvc.GetLogs(ctx)
}
