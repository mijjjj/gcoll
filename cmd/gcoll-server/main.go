package main

import (
	"github.com/gogf/gf/v2/os/gctx"

	internalcmd "github.com/mijjjj/gcoll/internal/cmd"
)

func main() {
	internalcmd.RunServer(gctx.GetInitCtx())
}
