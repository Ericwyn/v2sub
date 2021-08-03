package conf

import (
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/param"
	"strconv"
	"strings"
)

func ParseArgs(args []string) {
	param.AssistParamLength(args, 1)
	switch args[0] {
	case "sport": // -conf sport 1080 设置 socks 代理端口
		param.AssistParamLength(args, 2)
		setSocksPort(args[1])
		break
	case "hport": // -conf sport 1081 设置 hport 代理端口
		param.AssistParamLength(args, 2)
		setHttpPort(args[1])
		break
	case "lconn": // -conf lconn true 允许来自局域网的连接
		param.AssistParamLength(args, 2)
		setLocalConnEnable(args[1])
		break
	case "list": // -conf list 允许来自局域网的连接
		param.AssistParamLength(args, 1)
		log.I("Config SocksPort:       ", ServerConfigNow.SocksPort)
		log.I("Config HttpPort:        ", ServerConfigNow.HttpPort)
		log.I("Config AllLocalConnect: ", ServerConfigNow.AllowLocalConnect)

		break
	default:
		log.E("sub args error")
	}
}

func setSocksPort(port string) {
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum <= 0 || portNum >= 65534 {
		log.E("socks port error, port : " + port)
		return
	}

	log.I("set SocksPort to : ", portNum)
	ServerConfigNow.SocksPort = portNum
	FlushConfig()
}

func setHttpPort(port string) {
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum <= 0 || portNum >= 65534 {
		log.E("http port error, port : " + port)
		return
	}

	log.I("set HttpPort to : ", portNum)
	ServerConfigNow.HttpPort = portNum
	FlushConfig()
}

func setLocalConnEnable(enableFlag string) {
	enableFlag = strings.ToLower(enableFlag)
	enable := false
	if enableFlag == "t" || enableFlag == "true" || enableFlag == "1" {
		enable = true
	} else if enableFlag == "f" || enableFlag == "false" || enableFlag == "0" {
		enable = false
	} else {
		log.E("enable flag error, plase input true|false, flag now : " + enableFlag)
		return
	}

	log.I("set AllowLocalConnect to : ", enable)
	ServerConfigNow.AllowLocalConnect = enable
	FlushConfig()
}
