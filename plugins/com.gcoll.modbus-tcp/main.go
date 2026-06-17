package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	modbustcp "github.com/mijjjj/gcoll/plugins/com.gcoll.modbus-tcp/runtime"
)

// main 输出插件清单并保持进程运行，后续在此装配 gRPC 服务。
func main() {
	runtime := modbustcp.NewRuntime(nil)
	payload, err := json.MarshalIndent(runtime.Manifest(), "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(payload))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
