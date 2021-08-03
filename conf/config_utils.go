package conf

import (
	"encoding/json"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/storage"
)

// 载入 config
// 写 config
// 创建 config 文件夹
// 底层对接 storage

var initFlag = false
var serverConfigName = "server.json"
var subConfigName = "sub.json"

// 从本地读取配置文件
func LoadLocalConfig() {
	if !initFlag {
		subConfigBytes := storage.ReadConfigFileLocal(subConfigName)
		if string(subConfigBytes) != "" {
			err := json.Unmarshal(subConfigBytes, &SubConfigNow)
			if err != nil {
				log.E("parse sub config file to json error")
			}
		}

		serverConfigBytes := storage.ReadConfigFileLocal(serverConfigName)
		if string(serverConfigBytes) != "" {
			err := json.Unmarshal(serverConfigBytes, &ServerConfigNow)
			if err != nil {
				log.E("parse server config file to json error")
			}
			if ServerConfigNow.SocksPort == 0 {
				ServerConfigNow.SocksPort = 1080
			}
			if ServerConfigNow.HttpPort == 0 {
				ServerConfigNow.HttpPort = 1081
			}
		}

		if ServerConfigNow.ServerList == nil {
			ServerConfigNow.ServerList = make([]VServer, 0)
		}

		initFlag = true
	}
}

// 将配置文件保存到本地
func FlushConfig() {
	writeLocalConfig(SubConfigNow, ServerConfigNow)
}

// 将配置输出到本地文件中
func writeLocalConfig(subMap map[string]VSub, serverList ServerConfig) {
	mapJson, err := json.MarshalIndent(subMap, "", "    ")
	if err != nil {
		log.E("general sub map json error")
		panic(err)
	} else {
		storage.WriteConfigFileLocal(string(mapJson), subConfigName)
	}

	serversJson, err := json.MarshalIndent(serverList, "", "    ")
	if err != nil {
		log.E("general sub map json error")
		panic(err)
	} else {
		storage.WriteConfigFileLocal(string(serversJson), serverConfigName)
	}
}
