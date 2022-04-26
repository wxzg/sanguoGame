package gatecontrollor

import (
	"log"
	"sanguoServer/net"
	"strings"
	"sync"
)

var GateHandler = &Handler{
	proxyMap: make(map[string]map[int]*net.ProxyClient),
}
type Handler struct {
	proxyMutex sync.Mutex
	// 代理地址 -》 客户端连接（游戏客户端id） -》 连接 即当有客户端请求来了以后这里保存一下
	proxyMap map[string]map[int]*net.ProxyClient
	loginProxy string
	gameProxy string
}

func (h *Handler) Router(r *net.Router){
	
	g := r.Group("*")
	g.AddRouter("*", h.all)
}

func (h *Handler) all(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	log.Println("网关处理器")

	////account转发
	//name := req.Body.Name
	//proxyStr := ""
	//
	//if isAccount(name){
	//	proxyStr = h.loginProxy
	//}
	//
	//newProxyClient := &net.ProxyClient{
	//	Proxy: proxyStr,
	//}


}

//判断路由前缀是否是account
func isAccount(name string) bool {
	return strings.HasPrefix(name, "account.")
}
