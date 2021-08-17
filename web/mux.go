package web

import (
	"crypto/rand"
	"fmt"
	"github.com/Ericwyn/GoTools/file"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"math/big"
)

var adminPassword string

func NewMux(pw string) *gin.Engine {
	adminPassword = pw

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	loadStaticPath(router)

	store := cookie.NewStore(GeneralSessionKey())
	router.Use(sessions.Sessions("v2subw-session", store))

	router.Use(gin.Logger())
	//router.LoadHTMLGlob(".assests/*.html")

	initAPI(router)
	return router
}

// 设置 API 路由
func initAPI(router *gin.Engine) {

	// 登录
	router.POST("/login", apiLogin)

	apiV1 := router.Group("/api/v1", AuthMiddleware())
	{
		// v2sub - sub 相关
		// 获取 v2sub subs 配置（v2sub -sub list）
		apiV1.GET("/v2sub/subs/list", apiLogin)
		// 刷新 v2sub subs 配置（v2sub -sub update all）
		apiV1.GET("/v2sub/subs/updateAll", apiLogin)

		// v2sub - ser 相关
		// 获取 v2sub ser 配置 (v2sub -ser list)
		apiV1.GET("/v2sub/ser/list", apiLogin)
		// 设置某个 ser        (v2sub -ser set {id})
		apiV1.GET("/v2sub/ser/set", apiLogin)
		// 设置最快 ser        (v2sub -ser setx)
		apiV1.GET("/v2sub/ser/setx", apiLogin)

		// v2ray 连接配置
		// 获取当前 v2sub 设置 (v2sub -conf list)
		apiV1.GET("/v2sub/conf/list", apiLogin)
		// 设置 http port (v2sub -conf hport {http_port} )
		apiV1.GET("/v2sub/conf/hport", apiLogin)
		// 设置 socks port  (v2sub -conf sport {socks_port} )
		apiV1.GET("/v2sub/conf/sport", apiLogin)
		// 设置局域网连接  (v2sub -conf -conf lconn )
		apiV1.GET("/v2sub/conf/lconn", apiLogin)

		// 获取当前模板 (获取 /etc/v2sub/config_module.json)
		apiV1.GET("/v2sub/config_module", apiLogin)
		// 保存新的模板 (存储新的 config_module.json)
		apiV1.POST("/v2sub/config_module", apiLogin)
	}

}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		//// 设置session数据
		//session.Set("hadLogin", true)
		//session.Set("loginTime", time.Now().Unix())
		//session.Options(sessions.Options{MaxAge: 60 * 60 * 2}) // seconds
		if checkLogin(session) {
			ctx.Next()
		} else {
			ctx.JSON(401, gin.H{
				"code": RestApiAuthorizationError,
				"msg":  "未授权访问!",
			})
			ctx.Abort()
		}
	}
}

func loadStaticPath(router *gin.Engine) {
	staticDirPath := "./.assests"
	staticDir := file.OpenFile(staticDirPath)
	children := staticDir.Children()
	for _, child := range children {
		if child.IsDir() {
			fmt.Println("load static router:", "/"+child.Name(), "->", staticDirPath+"/"+child.Name())
			router.Static(child.Name(), staticDirPath+"/"+child.Name())
			router.Static("static/"+child.Name(), staticDirPath+"/"+child.Name())
		}
	}
}

var keyParisLen = 64

func GeneralSessionKey() []byte {
	return []byte(string(GeneralRandomStr(keyParisLen)))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*("

func GeneralRandomStr(length int) string {
	str := ""
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
		index64 := index.Int64()
		str += letterBytes[int(index64) : int(index64)+1]
	}
	return str
}
