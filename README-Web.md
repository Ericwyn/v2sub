# v2sub-web

v2subw (v2sub-web) 基于 v2sub 提供的一套 http api 接口

方便当我们把 v2sub 部署到路由器或者服务器上面的时候可以对 v2sub 进行远程控制

具体实现实际上是通过调用了本地 `v2sub` 命令来实现的功能

可自己搭配前端客户端实现远程控制功能（免去了你手动 shell 登录之后手动敲命令~）

## 运行参数和方法
```
  -b string
        binary, v2sub 的位置，默认使用 v2sub，可自定义为 v2sub 程序所在路径如("/opt/v2/v2sub")
  -k string
        admin key, 默认为 null，代表不需要密码即可访问 api 接口
  -p int
        run port，运行的端口，默认为 8886
```


## 接口列表
```
// 控制相关
/api/v1/v2sub/conn/start
/api/v1/v2sub/conn/stop
/api/v1/v2sub/conn/status
/api/v1/v2sub/conn/log

// v2sub - sub 相关
// 获取 v2sub subs 配置（v2sub -sub list）
/v2sub/subs/list

// 刷新 v2sub subs 配置（v2sub -sub update all）
/v2sub/subs/updateAll

// v2sub - ser 相关
// 获取 v2sub ser 配置 (v2sub -ser list)
/api/v1/v2sub/ser/list

// 设置某个 ser        (v2sub -ser set {id})
/api/v1/v2sub/ser/set

// 设置最快 ser        (v2sub -ser setx)
/api/v1/v2sub/ser/setx

// v2ray 连接配置
// 获取当前 v2sub 设置 (v2sub -conf list)
/api/v1/v2sub/conf/list

// 设置 http port (v2sub -conf hport {http_port} )
/api/v1/v2sub/conf/hport/set

// 设置 socks port  (v2sub -conf sport {socks_port} )
/api/v1/v2sub/conf/sport/set

// 设置局域网连接  (v2sub -conf -conf lconn )
/api/v1/v2sub/conf/lconn/set
```