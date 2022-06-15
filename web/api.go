package web

import (
	"fmt"
	"github.com/Ericwyn/v2sub/utils/command"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
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

func apiPacJs(ctx *gin.Context) {
	// 支持手动设定地址
	ip := ctx.Query("ip")
	hPort := ctx.Query("hPort")
	sPort := ctx.Query("sPort")
	ctx.String(200, "%s", renderPacJs(ip, hPort, sPort))
}

var runLog = make([]string, 0)
var lastStartTimeUnix int64 = 0

func apiConnStart(ctx *gin.Context) {
	if runFlag {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "v2ray is running now!",
		})
		return
	}
	// 协程执行 v2sub -conn start
	go v2subConnStart()

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
	v2subConnKill()
	runFlag = false
	lastStartTimeUnix = 0
	ctx.JSON(200, gin.H{
		"code": RestApiSuccess,
		"msg":  "stop v2ray",
	})
}

func apiConnRestart(ctx *gin.Context) {
	v2subConnKill()
	runFlag = false
	lastStartTimeUnix = 0
	apiConnStart(ctx)
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

func apiConnectClearLog(ctx *gin.Context) {
	runLog = make([]string, 0)

	ctx.JSON(200, gin.H{
		"code": RestApiSuccess,
		"msg":  "clear v2ray log",
	})
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
	result, err := v2subSubUpdateAll()
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
	result, err := v2subSerSetX()
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
	result, err := v2subConfList()
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
		result, err := v2subConfHttpPort(port)
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
		result, err := v2subConfSocksPort(port)
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
		result, err := v2subConfLocalConnect(enableStr)
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
