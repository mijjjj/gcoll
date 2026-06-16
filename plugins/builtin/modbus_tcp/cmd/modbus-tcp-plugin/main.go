package main

import (
	"encoding/json"
	"fmt"

	modbustcp "github.com/mijjjj/gcoll/plugins/builtin/modbus_tcp"
)

// main 当前只输出插件清单，后续由生成的 gRPC 服务装配进程入口。
func main() {
	runtime := modbustcp.NewRuntime(nil)
	payload, err := json.MarshalIndent(runtime.Manifest(), "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(payload))
}
