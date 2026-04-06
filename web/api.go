package web

import (
	"github.com/Ericwyn/v2sub/conn"
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
		session.Set("hadLogin", true)
		session.Set("loginTime", time.Now().Unix())
		session.Options(sessions.Options{MaxAge: 60 * 60 * 2})

		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  "登录成功",
		})
	} else {
		password, b := ctx.GetPostForm("password")
		if b && password == adminPassword {
			session := sessions.Default(ctx)
			session.Set("hadLogin", true)
			session.Set("loginTime", time.Now().Unix())
			session.Options(sessions.Options{MaxAge: 60 * 60 * 2})

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
	ip := ctx.Query("ip")
	hPort := ctx.Query("hPort")
	sPort := ctx.Query("sPort")
	ctx.String(200, "%s", renderPacJs(ip, hPort, sPort))
}

var runLog = make([]string, 0)
var lastStartTimeUnix int64 = 0

func apiConnStart(ctx *gin.Context) {
	status := conn.Status()
	if status.Running {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "v2ray is running now!",
			"data": status,
		})
		return
	}
	go v2subConnStart()

	ctx.JSON(200, gin.H{
		"code": RestApiSuccess,
		"msg":  "start v2ray",
	})
}

func apiConnStatus(ctx *gin.Context) {
	status := conn.Status()
	ctx.JSON(200, gin.H{
		"code": RestApiSuccess,
		"msg":  status.Message,
		"data": status,
	})
}

func apiConnStop(ctx *gin.Context) {
	v2subConnKill()
	runFlag = false
	lastStartTimeUnix = 0
	ctx.JSON(200, gin.H{
		"code": RestApiSuccess,
		"msg":  "stop v2ray",
		"data": conn.Status(),
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
		result, err := v2subSerSet(id)
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

func apiConfSPortSet(ctx *gin.Context) {
	port := ctx.Query("port")
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

func apiConfLConnSet(ctx *gin.Context) {
	enableFlag := ctx.Query("enable")
	result, err := v2subConfLocalConnect(enableFlag)
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
