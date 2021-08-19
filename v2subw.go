package main

import (
	"flag"
	"fmt"
	"github.com/Ericwyn/v2sub/conf"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/web"
	"net/http"
	"time"
)

// v2sub web 管理界面后台

var runPort = flag.Int("p", 8886, "run port，运行的端口")
var adminPassword = flag.String("k", "", "admin key, 默认为 null，代表不需要密码即可访问 api 接口")
var v2subPath = flag.String("b", "v2sub", "binary, v2sub 的位置，默认使用 v2sub，可自定义为 v2sub 程序所在路径")

func main() {
	flag.Parse()

	conf.LoadLocalConfig()

	log.I("[v2sub-w] 启动于", ":"+fmt.Sprint(*runPort))

	s := &http.Server{
		Addr:           ":" + fmt.Sprint(*runPort),
		Handler:        web.NewMux(*adminPassword, *v2subPath),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()

}
