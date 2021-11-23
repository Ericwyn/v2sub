package sub

import (
	"errors"
	"fmt"
	"github.com/Ericwyn/v2sub/ajax"
	"github.com/Ericwyn/v2sub/conf"
	"github.com/Ericwyn/v2sub/server"
	"github.com/Ericwyn/v2sub/utils/decode"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/param"
	"github.com/Ericwyn/v2sub/utils/putil"
	"os"
	"strings"
)

// sub 订阅管理， 订阅的增加/删除/修改

func ParseArgs(args []string) {
	param.AssistParamLength(args, 1)
	switch args[0] {
	case "a", "add": // -sub a {sub_name} / -sub add {sub_name}
		param.AssistParamLength(args, 3)
		if args[1] != "" && args[2] != "" {
			AddSub(args[1], args[2])
		} else {
			log.E("sub args error")
		}
	case "u", "update": // -sub u {sub_name} / -sub update {sub_name} 更新某一个订阅
		param.AssistParamLength(args, 2)
		UpdateSub(args[1])
	case "c", "customer": // -sub c {sub_name} {customer_result} 更新某一个订阅
		param.AssistParamLength(args, 3)
		UpdateSubCustomer(args[1], args[2])
	case "ua", "updateall": // -sub ua / -sub updateall 更新全部订阅
		param.AssistParamLength(args, 1)
		log.I("update all sub msgs")
		for name, _ := range conf.SubConfigNow {
			UpdateSub(name)
		}
	case "r", "remove": // -sub remove {sub_name}
		param.AssistParamLength(args, 2)
		RemoveSubByName(args[1])
	case "l", "list": //  -sub l / -sub list  展示当前订阅设置
		ListSubs()
		os.Exit(0)
	default:
		log.E("sub args error")
	}
}

// 将 sub 地址添加到 config
// 将
func AddSub(subName string, subPath string) {
	log.I("start add sub")
	if _, value := conf.SubConfigNow[subName]; value {
		log.E("sub name exits")
		return
	}

	sub := conf.VSub{
		SubUrl:  subPath,
		SubName: subName,
	}

	log.D("sub url: ", sub.SubUrl)
	// 对 sub 地址进行校验
	// 获取 sub 地址的 JSON
	serverList := checkSub(sub)
	if serverList == nil || len(serverList) == 0 {
		log.E("can't parse server msg from sub path")
		return
	}
	log.I("parse sub url success, get ", len(serverList), " server configs")

	// 保存数据
	conf.ServerConfigNow.ServerList = append(conf.ServerConfigNow.ServerList, serverList...)
	conf.SubConfigNow[subName] = sub

	conf.FlushConfig()
}

func UpdateSub(subName string) {
	log.I("start update sub")
	if _, value := conf.SubConfigNow[subName]; !value {
		log.E("sub name not exits")
		return
	}
	sub := conf.SubConfigNow[subName]

	log.D("sub url: ", sub.SubUrl)
	// 对 sub 地址进行校验
	// 获取 sub 地址的 JSON
	serverList := checkSub(sub)
	if serverList == nil || len(serverList) == 0 {
		log.E("can't parse server msg from sub path")
		return
	}
	log.I("parse sub url success, get ", len(serverList), " server configs")

	log.I("delete old sub msgs")
	// 删除旧的数据
	RemoveSubByName(subName)

	log.I("save new sub msgs")
	// 保存数据 ServerList 和 SubConfig
	conf.ServerConfigNow.ServerList = append(conf.ServerConfigNow.ServerList, serverList...)
	conf.SubConfigNow[subName] = sub

	// 判断配置 id 是否越界
	if conf.ServerConfigNow.Id >= len(conf.ServerConfigNow.ServerList) {
		conf.ServerConfigNow.Id = 0
	}

	log.I("save sub msgs success")

	conf.FlushConfig()
}

// 手动将请求的 result 设置到 subName
// 解决某些时候访问不到订阅地址的问题
func UpdateSubCustomer(subName string, subResult string) {
	log.I("start update sub by customer result")
	if _, value := conf.SubConfigNow[subName]; !value {
		log.E("sub name not exits")
		return
	}

	sub := conf.SubConfigNow[subName]

	serverList, err := parseSubMsg(subResult, subName)
	if err != nil {
		log.E("parse sub result error")
		return
	}

	if serverList == nil || len(serverList) == 0 {
		log.E("can't parse server msg from sub path")
		return
	}
	log.I("parse sub url success, get ", len(serverList), " server configs")

	log.I("delete old sub msgs")
	// 删除旧的数据
	RemoveSubByName(subName)

	log.I("save new sub msgs")
	// 保存数据 ServerList 和 SubConfig
	conf.ServerConfigNow.ServerList = append(conf.ServerConfigNow.ServerList, serverList...)
	conf.SubConfigNow[subName] = sub

	// 判断配置 id 是否越界
	if conf.ServerConfigNow.Id >= len(conf.ServerConfigNow.ServerList) {
		conf.ServerConfigNow.Id = 0
	}

	log.I("save sub msgs success")

	conf.FlushConfig()
}

func RemoveSubByUrl(subUrl string) {
	log.I("start delete sub")
	for subName, sub := range conf.SubConfigNow {
		if sub.SubUrl == subUrl {
			// 清除订阅设置
			delete(conf.SubConfigNow, subName)

			// 清除服务器配置
			newServerList := make([]conf.VServer, 0)
			for _, server := range conf.ServerConfigNow.ServerList {
				if server.SubName != subName {
					newServerList = append(newServerList, server)
				}
			}
			conf.ServerConfigNow.ServerList = newServerList
			conf.FlushConfig()
		}
	}
}

func RemoveSubByName(subNameInput string) {
	log.I("start delete sub")
	for subName, sub := range conf.SubConfigNow {
		if sub.SubName == subNameInput {
			// 清除订阅设置
			delete(conf.SubConfigNow, subName)

			// 清除服务器配置
			newServerList := make([]conf.VServer, 0)
			for _, server := range conf.ServerConfigNow.ServerList {
				if server.SubName != subName {
					newServerList = append(newServerList, server)
				}
			}

			conf.ServerConfigNow.ServerList = newServerList
			conf.FlushConfig()
			return
		}
	}
	log.E("can't find sub by name : " + subNameInput)
}

func ListSubs() {
	fmt.Println("=======================================================")
	fmt.Println(putil.F("name", 10), "url")
	for name, sub := range conf.SubConfigNow {
		fmt.Println(putil.F(name, 10), sub.SubUrl)
	}
	fmt.Println("=======================================================")
}

func checkSub(sub conf.VSub) []conf.VServer {
	var res []conf.VServer
	ajax.Send(ajax.Request{
		Url:    sub.SubUrl,
		Method: ajax.GET,
		//Header:  nil,
		Success: func(response *ajax.Response) {
			servers, err := parseSubMsg(response.Body, sub.SubName)
			if err != nil {
				log.E("parse sub msg error, subMsg: " + response.Body)
				panic(err)
			} else {
				res = servers
			}
		},
		Fail: func(status int, errMsg string) {
			err := errors.New(errMsg)
			log.E("get server msgs from sub url fail")
			panic(err)
		},
		Always: nil,
	})

	return res
}

// 解析订阅地址返回的 base64
func parseSubMsg(subResponse string, subName string) ([]conf.VServer, error) {
	serverList := make([]conf.VServer, 0)

	subMsg := decode.Base64Decode(subResponse)
	if subMsg == "" {
		err := errors.New("decode base64 fail while parse sub msg")
		return nil, err
	}
	subLinks := strings.Split(subMsg, "\n")
	vmessLinks := make([]string, 0)
	for _, link := range subLinks {
		if strings.Index(link, "vmess://") == 0 {
			//serverTemp := conf.VServer{
			//	Source:      link,
			//}
			//serverList = append(serverList, serverTemp)
			vmessLinks = append(vmessLinks, link)
		}
	}
	if len(vmessLinks) == 0 {
		err := errors.New("can't find the vmess link in sub msgs")
		return nil, err
	}
	// 解析 vmess 链接
	for _, vmLink := range vmessLinks {
		vmess, configJson := server.ParseVmessLink(vmLink)
		serverList = append(serverList, conf.VServer{
			SubName:    subName,
			Source:     vmLink,
			ConfigJson: configJson,
			Vmess:      *vmess,
		})
	}
	return serverList, nil
}
