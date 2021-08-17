package conn

import (
	"fmt"
	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/v2sub/conf"
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/param"
	"github.com/Ericwyn/v2sub/utils/putil"
	"github.com/shirou/gopsutil/process"
	"os"
	"strconv"
	"strings"
)

const v2rayBinPath = "/usr/local/bin/v2ray"

func ParseArgs(args []string) {
	param.AssistParamLength(args, 1)
	switch args[0] {

	case "start": // -conn start 启动 v2ray
		startV2ray()
		fmt.Println("v2ray 已停止")
	case "kill": // -conn stop 停止其他正在运行的 v2ray 和 v2sub
		KillV2Sub()
	default:
		log.E("sub args error")
	}
}

func checkV2ray() {
	vtoo := file.OpenFile(v2rayBinPath)
	if !vtoo.Exits() {
		log.E("can't find v2ray bin in " + v2rayBinPath)
		os.Exit(-1)
	}
}

func startV2ray() {
	log.I("start v2ray ......")

	checkV2ray()

	// 输出当前配置
	runConfig := conf.ServerConfigNow.ServerList[conf.ServerConfigNow.Id]
	conf.SaveDefaultServerConfig(runConfig)

	log.I("use config is :   ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ")
	log.I("========================================================================")
	log.I(
		putil.F("ID", 4),
		putil.F("别名", 50),
		putil.F("地址", 24),
		putil.F("端口", 10),
		putil.F("类型", 5),
	)
	log.I(putil.F(" "+strconv.Itoa(conf.ServerConfigNow.Id), 4),
		putil.F(runConfig.Vmess.Ps, 50),
		putil.F(runConfig.Vmess.Add, 24),
		putil.F(runConfig.Vmess.Port, 10),
		putil.F(runConfig.Vmess.Type, 5))
	log.I("========================================================================")

	log.I("v2ray config path : " + conf.GetV2rayConfigPath())
	fmt.Println()
	fmt.Println()

	err := command.Run(v2rayBinPath, "-config", conf.GetV2rayConfigPath())
	if err != nil {
		log.E("run command error", []string{v2rayBinPath, "-c", conf.GetV2rayConfigPath()})
		log.E(err.Error())
		os.Exit(-1)
	}
}

//
func KillV2Sub() {
	processes, _ := process.Processes()
	fmt.Println()
	pidCurrent := os.Getpid()
	for _, p := range processes {
		cmdline, err := p.Cmdline()
		if err == nil {
			if strings.Index(cmdline, "v2ray") >= 0 ||
				strings.Index(cmdline, "v2sub") >= 0 {
				//fmt.Println(cmdline)
				startCommand := strings.Split(cmdline, " ")[0]
				if strEndWith(startCommand, "v2ray") ||
					strEndWith(startCommand, "v2sub") && strEndWith(startCommand, "-conn") {
					if p.Pid != int32(pidCurrent) {
						fmt.Println("kill pid:", p.Pid, "-->", cmdline)
					}
				}

			}
		}
	}
}

func strEndWith(str string, end string) bool {
	return strings.Index(str, end)+len(end) == len(str)
}
