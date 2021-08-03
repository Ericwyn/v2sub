package server

import (
	"fmt"
	"github.com/Ericwyn/v2sub/conf"
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/decode"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/param"
	"github.com/Ericwyn/v2sub/utils/putil"
	"os"
	"strconv"
	"strings"
	"sync"
)

func ParseArgs(args []string) {
	param.AssistParamLength(args, 1)
	switch args[0] {
	case "setflush": // -ser setflush 将当前选中的节点输出到 /etc/v2sub/config.json
		param.AssistParamLength(args, 1)
		SaveDefaultConfig(strconv.Itoa(conf.ServerConfigNow.Id))
		break
	case "set": // 设置某个配置作为 v2ray 启动配置
		param.AssistParamLength(args, 2)
		SaveDefaultConfig(args[1])
		break
	case "setx": // -ser setx 测试 + 设置节点为最快节点
		SpeedTestAll(true)
		break
	case "speedtest": // -ser speedtest 测试
		SpeedTestAll(false)
		break
	case "list":
		ListServer()
		os.Exit(0)
	default:
		log.E("sub args error")
	}
}

func SpeedTestAll(setDefaultConfigFlag bool) {
	fmt.Println("=======================================================")
	fmt.Println(
		putil.F("ID", 4),
		putil.F("别名", 50),
		putil.F("地址", 24),
		putil.F("端口", 10),
		putil.F("类型", 5),
		putil.F("测速", 5),
	)
	var wg sync.WaitGroup
	var speedTestResultServer []conf.VServer = make([]conf.VServer, 0)
	for i, config := range conf.ServerConfigNow.ServerList {
		wg.Add(1)
		go SpeedTestFun(i, config, &wg, &speedTestResultServer)
	}
	wg.Wait()
	fmt.Println("=======================================================")

	if len(speedTestResultServer) > 0 {
		fastServer := speedTestResultServer[0]
		for i, server := range conf.ServerConfigNow.ServerList {
			if fastServer.Vmess.Ps == server.Vmess.Ps && fastServer.Vmess.Port == server.Vmess.Port &&
				fastServer.Vmess.Add == server.Vmess.Add {

				fmt.Println("最快节点为")
				if i == conf.ServerConfigNow.Id {
					fmt.Println(
						putil.F("["+strconv.Itoa(i)+"]", 4),
						putil.F(server.Vmess.Ps, 50),
						putil.F(server.Vmess.Add, 24),
						putil.F(server.Vmess.Port, 10),
						putil.F(server.Vmess.Net, 5),
					)
				} else {
					fmt.Println(
						putil.F(" "+strconv.Itoa(i), 4),
						putil.F(server.Vmess.Ps, 50),
						putil.F(server.Vmess.Add, 24),
						putil.F(server.Vmess.Port, 10),
						putil.F(server.Vmess.Net, 5),
					)
				}

				if setDefaultConfigFlag {
					log.I()
					log.I("set default server id : " + strconv.Itoa(i))
					conf.ServerConfigNow.Id = i
					conf.FlushConfig()
				}

				return
			}
		}
	}

}

// 返回一个
func SpeedTestFun(i int, server conf.VServer, wg *sync.WaitGroup, serverList *[]conf.VServer) {
	//执行 ping -c 4 baidu.com | grep '^rtt' | awk -F"/" '{print $5F}'
	result, err := command.RunResult("ping -c 3 " + server.Vmess.Add + " | grep '^rtt' | awk -F\"/\" '{print $5F}'")
	if err == nil {
		result = strings.Replace(result, " ", "", -1)
		result = strings.Replace(result, "\n", "", -1)
		result = strings.Replace(result, "\r", "", -1)
		speedMs, err := strconv.ParseFloat(result, 64)
		*serverList = append(*serverList, server)
		if err != nil {
			//return -1
			log.E(err)
		}

		if i == conf.ServerConfigNow.Id {
			fmt.Println(
				putil.F("["+strconv.Itoa(i)+"]", 4),
				putil.F(server.Vmess.Ps, 50),
				putil.F(server.Vmess.Add, 24),
				putil.F(server.Vmess.Port, 10),
				putil.F(server.Vmess.Net, 5),
				putil.F(fmt.Sprint(speedMs)+" ms", 5),
			)
		} else {
			fmt.Println(
				putil.F(" "+strconv.Itoa(i), 4),
				putil.F(server.Vmess.Ps, 50),
				putil.F(server.Vmess.Add, 24),
				putil.F(server.Vmess.Port, 10),
				putil.F(server.Vmess.Net, 5),
				putil.F(fmt.Sprint(speedMs)+" ms", 5),
			)
		}
	} else {
		log.E(err)
	}
	wg.Done()
}

func SaveDefaultConfig(id string) {
	fmt.Println("save v2ray config to /etc/v2sub/config.json")
	index, _ := strconv.Atoi(id)
	if index >= len(conf.ServerConfigNow.ServerList) || index < 0 {
		log.E("config id error")
		os.Exit(-1)
	}
	conf.ServerConfigNow.Id = index

	// 刷新配置
	vmess, configJson := ParseVmessLink(conf.ServerConfigNow.ServerList[index].Source)
	conf.ServerConfigNow.ServerList[index].Vmess = *vmess
	conf.ServerConfigNow.ServerList[index].ConfigJson = configJson

	conf.FlushConfig()
	conf.SaveDefaultServerConfig(conf.ServerConfigNow.ServerList[index])

	fmt.Println("save success")
}

func ListServer() {
	fmt.Println("=======================================================")
	fmt.Println(
		putil.F("ID", 4),
		putil.F("别名", 50),
		putil.F("地址", 24),
		putil.F("端口", 10),
		putil.F("类型", 5),
	)
	for i, config := range conf.ServerConfigNow.ServerList {
		if i == conf.ServerConfigNow.Id {
			fmt.Println(
				putil.F("["+strconv.Itoa(i)+"]", 4),
				putil.F(config.Vmess.Ps, 50),
				putil.F(config.Vmess.Add, 24),
				putil.F(config.Vmess.Port, 10),
				putil.F(config.Vmess.Type, 5),
			)
		} else {
			fmt.Println(
				putil.F(" "+strconv.Itoa(i), 4),
				putil.F(config.Vmess.Ps, 50),
				putil.F(config.Vmess.Add, 24),
				putil.F(config.Vmess.Port, 10),
				putil.F(config.Vmess.Type, 5),
			)
		}
	}
	fmt.Println("=======================================================")

}

// 解析 vmess 链接， 得到具体 vmess 配置信息以及
// 以及一个 v2ray 的 config json
func ParseVmessLink(vmessStr string) (*conf.VmessJson, string) {
	if strings.Index(vmessStr, "vmess://") == 0 {
		vmessBase64 := vmessStr[8:len(vmessStr)] // 去除前缀
		vmessJson := decode.VmessBase64Decode(vmessBase64)
		log.I("get vmess json: ", vmessJson)
		// 通过 vmess 链接来获取 config.VServer 对象
		vmessJsonObj, configJson := conf.ParseVmessConfigToConfigJson(vmessJson)

		// 设置 http 和 socks 的链接
		configJson = strings.Replace(configJson, "{sPort}",
			strconv.Itoa(conf.ServerConfigNow.SocksPort),
			1)

		configJson = strings.Replace(configJson, "{hPort}",
			strconv.Itoa(conf.ServerConfigNow.HttpPort),
			1)

		// 设置是否允许来自局域网的连接
		bindAddr := "127.0.0.1"
		if conf.ServerConfigNow.AllowLocalConnect {
			bindAddr = "0.0.0.0"
		}
		configJson = strings.Replace(configJson, "{bindAddr}", bindAddr, -1)

		return vmessJsonObj, configJson
	} else {
		return nil, ""
	}
}
