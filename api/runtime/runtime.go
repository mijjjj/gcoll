// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package runtime

import (
	"context"

	"github.com/mijjjj/gcoll/api/runtime/v1"
)

type IRuntimeV1 interface {
	Health(ctx context.Context, req *v1.HealthReq) (res *v1.HealthRes, err error)
	Overview(ctx context.Context, req *v1.OverviewReq) (res *v1.OverviewRes, err error)
}
