package web

import (
	"crypto/rand"
	"fmt"
	"github.com/Ericwyn/GoTools/file"
	"github.com/Ericwyn/v2sub/utils/log"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"os"
	"time"
)

var v2subBinPath = "v2sub"
var adminPassword string

var serverJson string
var subsJson string

func StartApiServer(runPort int, pw string, v2subPath string) {
	// 首次载入本地的配置文件
	initLoadConfig()

	// 监听 v2sub 的配置文件
	go startFsNotify()

	//return
	s := &http.Server{
		Addr:           ":" + fmt.Sprint(runPort),
		Handler:        newMux(pw, v2subPath),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}

func newMux(pw string, binPath string) *gin.Engine {
	adminPassword = pw
	v2subBinPath = binPath

	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

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

	// 方便 Android 设备使用, 此处返回 pac.js 文件
	router.GET("/pac.js", apiPacJs)

	var apiV1 *gin.RouterGroup
	if adminPassword != "" {
		apiV1 = router.Group("/api/v1", AuthMiddleware())
	} else {
		apiV1 = router.Group("/api/v1")
	}

	// 放开跨域请求
	apiV1.Use(CorsMiddleware())
	{

		apiV1.GET("/v2sub/conn/start", apiConnStart)
		apiV1.GET("/v2sub/conn/stop", apiConnStop)
		apiV1.GET("/v2sub/conn/restart", apiConnRestart)
		apiV1.GET("/v2sub/conn/status", apiConnStatus)
		apiV1.GET("/v2sub/conn/log", apiConnLog)
		//apiV1.GET("/v2sub/")

		// v2sub - sub 相关
		// 获取 v2sub subs 配置（v2sub -sub list）
		apiV1.GET("/v2sub/subs/list", apiSubsList)
		// 刷新 v2sub subs 配置（v2sub -sub update all）
		apiV1.GET("/v2sub/subs/updateAll", apiSubsUpdateAll)

		// v2sub - ser 相关
		// 获取 v2sub ser 配置 (v2sub -ser list)
		apiV1.GET("/v2sub/ser/list", apiServersList)
		// 设置某个 ser        (v2sub -ser set {id})
		apiV1.GET("/v2sub/ser/set", apiServersSet)
		// 设置最快 ser        (v2sub -ser setx)
		apiV1.GET("/v2sub/ser/setx", apiServersSetX)

		// v2ray 连接配置
		// 获取当前 v2sub 设置 (v2sub -conf list)
		apiV1.GET("/v2sub/conf/list", apiConfList)
		// 设置 http port (v2sub -conf hport {http_port} )
		apiV1.GET("/v2sub/conf/hport/set", apiConfHPortSet)
		// 设置 socks port  (v2sub -conf sport {socks_port} )
		apiV1.GET("/v2sub/conf/sport/set", apiConfSPortSet)
		// 设置局域网连接  (v2sub -conf -conf lconn )
		apiV1.GET("/v2sub/conf/lconn/set", apiConfLConnSet)

		// 获取当前模板 (获取 /etc/v2sub/config_module.json)
		//apiV1.GET("/v2sub/config_module", apiLogin)
		//// 保存新的模板 (存储新的 config_module.json)
		//apiV1.POST("/v2sub/config_module", apiLogin)
	}

}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)

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

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//method := c.Request.Method

		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		defer func() {
			if err := recover(); err != nil {
				log.E("Panic info is: %v", err)
			}
		}()

		c.Next()
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

func initLoadConfig() {
	// 载入配置
	serverConfigFile := file.OpenFile("/etc/v2sub/server.json")
	subConfigFile := file.OpenFile("/etc/v2sub/sub.json")
	if !serverConfigFile.IsFile() || !subConfigFile.IsFile() {
		log.E("can't not find the server.json and sub.json in /etc/v2sub/")
		os.Exit(-1)
	}

	serverJsonBytes, err := serverConfigFile.Read()
	if err != nil {
		log.E("can't read /etc/v2sub/server.json")
		log.E(err.Error())
		os.Exit(-1)
	}
	serverJson = string(serverJsonBytes)

	subJsonBytes, err := subConfigFile.Read()
	if err != nil {
		log.E("can't read /etc/v2sub/sub.json")
		log.E(err.Error())
		os.Exit(-1)
	}
	subsJson = string(subJsonBytes)
}

func fsNotifySyncConfig() {
	// 载入配置
	serverConfigFile := file.OpenFile("/etc/v2sub/server.json")
	subConfigFile := file.OpenFile("/etc/v2sub/sub.json")

	serverJsonBytes, err := serverConfigFile.Read()
	if err == nil {
		serverJson = string(serverJsonBytes)
	} else {
		log.E("read /etc/v2sub/server.json fail")
		log.E(err.Error())
	}

	subJsonBytes, err := subConfigFile.Read()
	if err == nil {
		subsJson = string(subJsonBytes)
	} else {
		log.E("read /etc/v2sub/sub.json fail")
		log.E(err.Error())
	}
}

// 监听 /etc/v2sub/ 下 server.json  sub.json 文件的变化
// 当文件有变化的时候重新载入
func startFsNotify() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.E(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// v2sub
				// 如果是删除了数据的话
				// 延迟 1s 钟后重新监听
				if event.Op == fsnotify.Remove {
					_ = watcher.Remove(event.Name)
					time.AfterFunc(time.Microsecond*500, func() {
						log.I("fsNotify: get config update event from v2sub")
						log.I("fsNotify:update: ", event.Name)
						_ = watcher.Add(event.Name)
						if err != nil {
							log.E(err)
						}
						fsNotifySyncConfig()
					})
				} else {
					log.I("fsNotify: get config update event from file edit")
					log.I("fsNotify:update: ", event.Name)
					fsNotifySyncConfig()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.E("error:", err)
			}
		}
	}()

	// if this is a link, it will follow all the links and watch the file pointed to
	err = watcher.Add("/etc/v2sub/server.json")
	if err != nil {
		log.E(err)
	}
	err = watcher.Add("/etc/v2sub/sub.json")
	if err != nil {
		log.E(err)
	}

	log.I("start watch config files")

	<-done
}
