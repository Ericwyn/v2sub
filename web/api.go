package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

const RestApiParamError = "4000"
const RestApiAuthorizationError = "4001"
const RestApiServerError = "4003"
const RestApiSuccess = "1000"

func apiLogin(ctx *gin.Context) {
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
