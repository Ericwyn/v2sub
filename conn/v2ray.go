package conn

import (
	"fmt"
	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/v2sub/conf"
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/putil"
	"os"
	"strconv"
)

// 对接 v2ray

const v2rayBinPath = "/usr/local/bin/v2ray"

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

//func runCommand(name string, args ...string) {
//	cmd := exec.Command(name, args...)
//	var stdout io.ReadCloser
//	var err error
//	if stdout, err = cmd.StdoutPipe(); err != nil {     //获取输出对象，可以从该对象中读取输出结果
//		log.E(err)
//	}
//	defer stdout.Close()   // 保证关闭输出流
//
//	if err := cmd.Start(); err != nil {   // 运行命令
//		log.E(err)
//	}
//
//	if opBytes, err := ioutil.ReadAll(stdout); err != nil {  // 读取输出结果
//		log.E(err)
//	} else {
//		fmt.Println(string(opBytes))
//	}
//}

