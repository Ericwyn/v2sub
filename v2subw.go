package main

import (
	"flag"
	"fmt"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/Ericwyn/v2sub/web"
	"net/http"
	"time"
)

// v2sub web 管理界面后台

var runPort = flag.Int("p", 8886, "run port")
var adminPassword = flag.String("k", "v2sub", "admin key")

func main() {
	flag.Parse()

	log.I("v2sub-w 启动于", ":"+fmt.Sprint(*runPort))

	s := &http.Server{
		Addr:           ":" + fmt.Sprint(*runPort),
		Handler:        web.NewMux(*adminPassword),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()

}
