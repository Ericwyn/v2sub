package main

import (
	"fmt"
	"github.com/Ericwyn/v2sub/conf"
	"github.com/Ericwyn/v2sub/conn"
	"github.com/Ericwyn/v2sub/server"
	"github.com/Ericwyn/v2sub/sub"
	"github.com/Ericwyn/v2sub/utils/log"
	"os"
)

const versionMsg = "Beta 1.0.1"

func main() {

	bootArgs := os.Args

	if len(os.Args) < 2 {
		// 输出 help
		printArgsHelp()
		os.Exit(0)
	}

	conf.LoadLocalConfig()

	// 解析 args，依据不同的 args 进行不同的业务
	// 将运行路径参数删除掉
	parseArg(bootArgs[1:])
}

func parseArg(args []string) {
	switch args[0] {
	case "-h", "--help":
		printArgsHelp()
		os.Exit(0)
	case "-v", "--version", "--v":
		fmt.Println("v2sub", versionMsg)
		fmt.Println("https://github.com/Ericwyn/v2sub")
	case "-sub":
		sub.ParseArgs(args[1:])
	case "-conf": // 设置端口/局域网连接
		conf.ParseArgs(args[1:])
	case "-ser":
		server.ParseArgs(args[1:])
	case "-conn":
		conn.ParseArgs(args[1:])
	default:
		log.E("param error, use -h can get the params help")
		os.Exit(-1)
	}
}

func printArgsHelp() {
	fmt.Println(
		`订阅管理:
    -sub add {name} {url} 
        添加一个订阅，订阅节点自动增加到 ser list
    -sub update {name} 
        更新一个订阅
    -sub update all 
        更新全部订阅结果
    -sub remove {name} 
        删除一个订阅
    -sub list 
        查看当前所有订阅

节点查看:
    -ser list 
        查看所有节点
    -ser set {ser_id} 
        设置默认节点
    -ser setx 
        对节点进行 ping 测速，之后将默认节点设置为最快节点
    -ser setflush
        将当前选择的节点输出到 /etc/v2sub/config.json
    -ser speedtest
        使用 tcping 查看各个节点的连接速度
    
连接配置管理
    -conf sport {socket_port} 
        socket 端口号管理, 默认 1080
    -conf hport {http_port} 
        http 端口号管理， 默认1081
    -conf lconn {true|false} 
        是否允许来自局域网的连接，默认为 false
    -conf list
        展示当前的 port、lconn 配置
  
连接
    -conn start 
        启用 v2ray 连接 server
    -conn kill 
        停止 v2ray （kill 掉其他 v2sub 和 v2ray）

其他
    -v, --version
        查看版本号
    -h, --help
        查看帮助说明
`)
}
