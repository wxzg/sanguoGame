package gatecontrollor

import (
	"log"
	config "sanguoServer/conf"
	"sanguoServer/net"
	"sanguoServer/utils"
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
	//创建路由组别
	g := r.Group("*")
	//  */* 匹配一切路由，全都放到h.all中处理
	g.AddRouter("*", h.all)
}

func (h *Handler) all(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//先拿到路由名
	name := req.Body.Name
	proxyStr := ""
	//判断是不是登录服务
	if isAccount(name){
		proxyStr = h.loginProxy
	}

	if proxyStr == "" {
		rsp.Body.Code = utils.ProxyNotInConnect
		return
	}

	h.proxyMutex.Lock()
	_, ok := h.proxyMap[proxyStr]
	if !ok {
		h.proxyMap[proxyStr] = make(map[int64]*net.ProxyClient)
	}
	h.proxyMutex.Unlock()

	cidValue, err := req.Conn.GetProperty("cid")
	if err != nil {
		log.Println("cid未取到", err)
		rsp.Body.Code = utils.InvalidParam
		return
	}
	cid := cidValue.(int64)
	proxy, ok := h.proxyMap[proxyStr][cid]
	if !ok {
		proxy = net.NewProxyClient(proxyStr)
		err := proxy.Connect()
		if err != nil{
			h.proxyMutex.Lock()
			delete(h.proxyMap[proxyStr],cid)
			h.proxyMutex.Unlock()
			rsp.Body.Code = utils.ProxyConnectError
			return
		}
		//把连接存起来
		h.proxyMutex.Lock()
		h.proxyMap[proxyStr][cid] = proxy
		h.proxyMutex.Unlock()

		proxy.SetProperty("cid", cid)
		proxy.SetProperty("gateConn", req.Conn)
		proxy.SetProperty("proxy", proxyStr)
		proxy.SetOnPush(h.onPush)
	}

	rsp.Body.Seq = req.Body.Seq
	rsp.Body.Name = req.Body.Name
	log.Println(req)
	r, err := proxy.Send(req.Body.Name, req.Body.Msg)
	if err == nil {
		rsp.Body.Code = r.Code
		rsp.Body.Msg = r.Msg
	} else {
		rsp.Body.Code = utils.ProxyConnectError
		rsp.Body.Msg = nil
	}
}

func (h *Handler) onPush(conn *net.ClientConn, body *net.RspBody) {
	gateConn, err := conn.GetProperty("gateConn")
	if err != nil {
		return
	}
	wc := gateConn.(net.WSConn)
	wc.Push(body.Name, body.Msg)
}

//判断路由前缀是否是account
func isAccount(name string) bool {
	return strings.HasPrefix(name, "account.")
}

