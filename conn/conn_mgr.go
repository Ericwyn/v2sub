package conn

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/v2sub/conf"
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/param"
	"github.com/Ericwyn/v2sub/utils/putil"
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

	err := command.RunSync(v2rayBinPath, "-config", conf.GetV2rayConfigPath())
	if err != nil {
		log.E("run command error", []string{v2rayBinPath, "-c", conf.GetV2rayConfigPath()})
		log.E(err.Error())
		os.Exit(-1)
	}
}

func KillV2Sub() {
	grep := exec.Command("grep", "v2") // 根据v2关键字进行模糊查询进程
	ps := exec.Command("ps", "cax")

	// Get ps's stdout and attach it to grep's stdin.
	pipe, _ := ps.StdoutPipe()
	defer pipe.Close()

	grep.Stdin = pipe

	// Run ps first.
	ps.Start()

	// Run and get the output of grep.
	res, _ := grep.Output()
	resL := string(res)

	processListStr := strings.Split(resL, "\n")
	for _, pc := range processListStr {
		elemList := strings.Split(pc, " ")

		pid := elemList[0]
		pName := elemList[len(elemList)-1]
		if pName != "v2ray" && pName != "v2sub" {
			continue
		}

		_ = command.RunSync("kill", fmt.Sprint(pid))
	}
}
