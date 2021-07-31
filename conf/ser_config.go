package conf

import (
	"encoding/json"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/utils/storage"
	"strconv"
	"strings"
)

// 记录一个服务器， 主要记录 vmess 连接和对应的 config.json
type ServerConfig struct {
	Id         int
	ServerList []VServer
}

type VServer struct {
	//Note       string
	SubName    string
	Vmess      VmessJson
	Source     string // 原始 vmess:// 连接
	ConfigJson string // 解析得到的 v2ray 配置文件
}

var ServerConfigNow ServerConfig

type VmessJson struct {
	Ps   string `json:"ps"`
	Add  string `json:"add"`
	Port string `json:"port"`
	ID   string `json:"id"`
	Aid  int    `json:"aid"`
	Net  string `json:"net"`
	Type string `json:"type"`
	TLS  string `json:"tls"`
}

// 返回 Source 对象 + v2ray 的 json 配置文件
func ParseVmessConfigToConfigJson(vmessJsonStr string) (*VmessJson, string) {
	var vmess VmessJson
	err := json.Unmarshal([]byte(vmessJsonStr), &vmess)
	if err != nil {
		log.E("parse vmess config fail : " + vmessJsonStr)
		return nil, ""
	} else {
		return &vmess, parseVmessJson(vmess)
	}
}

func SaveDefaultServerConfig(server VServer) {
	if server.ConfigJson != "" {
		storage.WriteConfigFileLocal(server.ConfigJson, "config.json")
	}
}

func GetV2rayConfigPath() string {
	return storage.GetConfigDirPath() + "/config.json"
}

func parseVmessJson(vmess VmessJson) string {
	module := storage.LoadV2ConfigModule()

	module = strings.Replace(module, "{Add}", vmess.Add, 1)
	module = strings.Replace(module, "{ID}", vmess.ID, 1)
	module = strings.Replace(module, "{Aid}", strconv.Itoa(vmess.Aid), 1)
	module = strings.Replace(module, "{Port}", vmess.Port, 1)
	module = strings.Replace(module, "{Net}", vmess.Net, 1)
	//module = strings.Replace(module, "{}", vmess.Type, 1)
	//module = strings.Replace(module, "{}", vmess.TLS, 1)

	return module
}
