package gatecontrollor

import (
	"log"
	config "sanguoServer/conf"
	"sanguoServer/net"
	"strings"
	"sync"
)

var GateHandler = &Handler{
	proxyMap: make(map[string]map[int64]*net.ProxyClient),
}
type Handler struct {
	proxyMutex sync.Mutex
	// 代理地址 -》 客户端连接（游戏客户端id） -》 连接 即当有客户端请求来了以后这里保存一下
	proxyMap map[string]map[int64]*net.ProxyClient
	loginProxy string
	gameProxy string
}

func (h *Handler) Router(r *net.Router){
	//登录服务器代理
	h.loginProxy = config.File.MustValue("gate_server", "login_proxy", "ws://127.0.0.1:8003")
	//游戏服务器代理
	h.gameProxy = config.File.MustValue("gate_server", "game_proxy", "ws://127.0.0.1:8001")
	g := r.Group("*")
	g.AddRouter("*", h.all)
}

func (h *Handler) all(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	log.Println("网关处理器")

	//account转发
	name := req.Body.Name
	proxyStr := ""
	//判断是不是登录服务
	if isAccount(name){
		proxyStr = h.loginProxy
	}

	loginProxy := net.NewProxyClient(proxyStr)

	loginProxy.Connect()
}

//判断路由前缀是否是account
func isAccount(name string) bool {
	return strings.HasPrefix(name, "account.")
}
