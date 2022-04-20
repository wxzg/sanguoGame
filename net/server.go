package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

/*
	1. 先定义 server 结构体来表示server服务，只要包含路由以及地址（127.0.0.1:8003)
*/
type server struct {
	addr string
	router *Router
}
/*
	2. 暴露对外定义的接口，返回一个*server
*/
func NewServer(addr string) *server  {
	return &server{
		addr: addr,
	}
}

func (s *server) Router(router *Router) {
	s.router = router
}

/*
	3. 启动服务，包括http监听，然后主要的逻辑处理在wsHandler里面
*/
func (s *server) Start()  {
	http.HandleFunc("/",s.wsHandler)
	err := http.ListenAndServe(s.addr,nil)
	if err != nil {
		panic(err)
	}
}

// http升级到websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request)  {
	//思考 websocket
	//1. http协议升级为websocket协议
	wsConn,err := wsUpgrader.Upgrade(w,r,nil)

	if err != nil {
		//打印日志 同时 退出应用程序
		log.Println("websocket服务连接出错",err)
	}
	//websocket通道建立之后 不管是客户端还是服务端 都可以收发消息
	//发消息的时候 把消息当做路由 来去处理 消息是有格式的，先定义消息的格式
	//客户端 发消息的时候 {Name:"account.login"} 收到之后 进行解析 认为想要处理登录逻辑

	wsServer := NewWsServer(wsConn)
	wsServer.Router(s.router)
	wsServer.Start()
	wsServer.Handshake()
}