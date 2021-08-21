package web

import (
	"fmt"
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

const RestApiParamError = "4000"
const RestApiAuthorizationError = "4001"
const RestApiServerError = "4003"
const RestApiSuccess = "1000"

func apiLogin(ctx *gin.Context) {
	if adminPassword == "" {
		session := sessions.Default(ctx)
		// 设置session数据
		session.Set("hadLogin", true)
		session.Set("loginTime", time.Now().Unix())
		session.Options(sessions.Options{MaxAge: 60 * 60 * 2}) // seconds

		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  "登录成功",
		})
	} else {
		password, b := ctx.GetPostForm("password")
		if b && password == adminPassword {
			session := sessions.Default(ctx)
			// 设置session数据
			session.Set("hadLogin", true)
			session.Set("loginTime", time.Now().Unix())
			session.Options(sessions.Options{MaxAge: 60 * 60 * 2}) // seconds

			ctx.JSON(200, gin.H{
				"code": RestApiSuccess,
				"msg":  "登录成功",
			})
		} else {
			ctx.JSON(200, gin.H{
				"code": RestApiAuthorizationError,
				"msg":  "密码错误!",
			})
		}
	}
}

var runLog = make([]string, 0)
var runFlag = false
var lastStartTimeUnix int64 = 0

func apiConnStart(ctx *gin.Context) {
	if runFlag {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "v2ray is running now!",
		})
		return
	}
	// 协程启动 v2sub
	go func() {
		runFlag = true
		lastStartTimeUnix = time.Now().Unix()
		_ = command.RunSyncForResultCb(func(s string) {
			//fmt.Print(s)
			s = strings.Replace(s, "\u0000", "", -1)
			s = strings.Replace(s, "\t", "", -1)
			s = strings.Replace(s, "\r", "", -1)
			//s = strings.Split(s, "\n")[0]
			runLog = append(runLog, s)
		}, v2subBinPath, "-conn", "kill")
		_ = command.RunSyncForResultCb(func(s string) {
			//fmt.Print(s)
			s = strings.Replace(s, "\u0000", "", -1)
			s = strings.Replace(s, "\t", "", -1)
			s = strings.Replace(s, "\r", "", -1)
			//s = strings.Split(s, "\n")[0]
			runLog = append(runLog, s)
		}, v2subBinPath, "-conn", "start")

		// 阻塞住了
		runFlag = false
	}()
	ctx.JSON(200, gin.H{
		"code": RestApiSuccess,
		"msg":  "start v2ray",
	})
}

func apiConnStatus(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"code": RestApiSuccess,
		"msg":  "",
		"data": map[string]interface{}{
			"running":           fmt.Sprint(runFlag),
			"lastStartTimeUnix": lastStartTimeUnix,
		},
	})
}

func apiConnStop(ctx *gin.Context) {
	_ = command.RunSyncForResultCb(func(s string) {
		log.I(s)
	}, v2subBinPath, "-conn", "kill")
	runFlag = false
	lastStartTimeUnix = 0
	ctx.JSON(200, gin.H{
		"code": RestApiSuccess,
		"msg":  "stop v2ray",
	})
}

func apiConnLog(ctx *gin.Context) {
	maxLogResLen := 4000
	var resLog []string
	if len(runLog) > maxLogResLen {
		resLog = runLog[len(runLog)-maxLogResLen:]
	} else {
		resLog = runLog
	}
	res := ""
	for _, s := range resLog {
		res += s
	}
	ctx.String(200, "%s", res)
}

func apiSubsList(ctx *gin.Context) {
	//ctx.JSON(200,gin.H {
	//	"code": RestApiSuccess,
	//	"data": conf.SubConfigNow,
	//})

	if subsJson == "" {
		ctx.JSON(200, gin.H{
			"code": RestApiServerError,
			"msg":  "read fail",
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  subsJson,
		})
	}

}

func apiSubsUpdateAll(ctx *gin.Context) {
	//sub.ParseArgs([]string{"-sub", "updateall"})
	//ctx.JSON(200, "update success")
	result, err := command.RunResult(v2subBinPath + " -sub updateall")
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": RestApiServerError,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  result,
		})
	}
}

func apiServersList(ctx *gin.Context) {
	//ctx.JSON(200, gin.H {
	//	"code": RestApiSuccess,
	//	"data": conf.ServerConfigNow.ServerList,
	//})
	if serverJson == "" {
		ctx.JSON(200, gin.H{
			"code": RestApiServerError,
			"msg":  "read fail",
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  serverJson,
		})
	}
}

func apiServersSet(ctx *gin.Context) {
	id := ctx.Query("id")
	log.I("[v2sub-w] set server id :" + id)

	index, err := strconv.Atoi(id)
	if err != nil || index < 0 {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "id error",
		})
	} else {
		result, err := command.RunResult(v2subBinPath + " -ser set " + id)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": RestApiServerError,
				"msg":  err.Error(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"code": RestApiSuccess,
				"msg":  result,
			})
		}
	}
}

func apiServersSetX(ctx *gin.Context) {
	//speedSorts := server.SortBySpeedTest(conf.ServerConfigNow.ServerList)
	//server.SaveDefaultConfig(strconv.Itoa(speedSorts[0].Index))
	//ctx.JSON(200, gin.H{
	//	"code": RestApiSuccess,
	//	"msg": "set config index to : " + strconv.Itoa(speedSorts[0].Index),
	//	"data": speedSorts,
	//})
	result, err := command.RunResult(v2subBinPath + " -ser setx")
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": RestApiServerError,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  result,
		})
	}
}

func apiConfList(ctx *gin.Context) {
	result, err := command.RunResult(v2subBinPath + " -conf list")
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": RestApiServerError,
			"msg":  err.Error(),
		})
	} else {
		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  result,
		})
	}
}

func apiConfHPortSet(ctx *gin.Context) {
	port := ctx.Query("port")
	log.I("[v2sub-w] set http port :" + port)

	index, err := strconv.Atoi(port)
	if err != nil || index < 0 || index > 65535 {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "port error",
		})
	} else {
		result, err := command.RunResult(v2subBinPath + " -conf hport " + port)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": RestApiServerError,
				"msg":  err.Error(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"code": RestApiSuccess,
				"msg":  result,
			})
		}
	}
}

func apiConfSPortSet(ctx *gin.Context) {
	port := ctx.Query("port")
	log.I("[v2sub-w] set http port :" + port)

	index, err := strconv.Atoi(port)
	if err != nil || index < 0 || index > 65535 {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "port error",
		})
	} else {
		result, err := command.RunResult(v2subBinPath + " -conf sport " + port)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": RestApiServerError,
				"msg":  err.Error(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"code": RestApiSuccess,
				"msg":  result,
			})
		}
	}
}

func apiConfLConnSet(ctx *gin.Context) {
	enable := ctx.Query("enable")
	log.I("[v2sub-w] set lconn enable :" + enable)

	if enable != "1" && enable != "0" && enable != "false" && enable != "true" {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "enable param error",
		})
	} else {
		enableStr := "false"
		if enable == "1" || enable == "true" {
			enableStr = "true"
		}
		result, err := command.RunResult(v2subBinPath + " -conf lconn " + enableStr)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": RestApiServerError,
				"msg":  err.Error(),
			})
		} else {
			ctx.JSON(200, gin.H{
				"code": RestApiSuccess,
				"msg":  result,
			})
		}
	}
}

func checkLogin(session sessions.Session) bool {
	if session != nil &&
		session.Get("hadLogin") != nil &&
		session.Get("hadLogin").(bool) && session.Get("loginTime") != nil {

		lastLoginTime := session.Get("loginTime").(int64)
		// 3600s * 2 = 2小时
		if time.Now().Unix()-lastLoginTime > (60 * 60 * 2) {
			session.Delete("hadLogin")
			session.Delete("loginTime")
			return false
		}
		return true
	} else {
		return false
	}
}
