package conn

import (
	"fmt"
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/param"
	"github.com/shirou/gopsutil/process"
	"strings"
)

func ParseArgs(args []string) {
	param.AssistParamLength(args, 1)
	switch args[0] {

	case "start": // -conn start 启动 v2ray
		startV2ray()

	case "kill": // -conn stop 停止其他正在运行的 v2ray 和 v2sub
		KillV2Sub()
	default:
		log.E("sub args error")
	}
}

//
func KillV2Sub() {
	processes, _ := process.Processes()
	fmt.Println()
	for _, p := range processes {
		cmdline, err := p.Cmdline()
		if err == nil {
			if strings.Index(cmdline, "/v2ray") >= 0 ||
				strings.Index(cmdline, "/v2sub") >= 0 {
				//fmt.Println(cmdline)
				startCommand := strings.Split(cmdline, " ")[0]
				if strEndWith(startCommand, "/v2ray") || strEndWith(startCommand, "/v2sub") {
					fmt.Println("kill pid:", p.Pid, "-->", cmdline)
					command.Run("kill", fmt.Sprint(p.Pid))
				}
			}
			//fmt.Println(cmdline)
		}
	}
}

func strEndWith(str string, end string) bool {
	return strings.Index(str, end)+len(end) == len(str)
}
