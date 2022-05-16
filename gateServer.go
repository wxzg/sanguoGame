package main

import (
	"log"
	config "sanguoServer/conf"
	"sanguoServer/net"
	"sanguoServer/server/gate"
)

/**
	网关实现
	1.登录（account.login） - 需要通过网关 - 转发到对应的登录服务器上
	2.网关需要跟websocket服务器（登录服务器）进行交互
	3.网关又要跟游戏客户端进行交互，也就是说网关既是websocket服务端也是客户端
	4.现在就是要实现websocket客户端
	5.网关：代理服务器（代理地址、代理的连接通道） 客户端连接（websocket）
	6.路由：路由 - 接收所有请求 网关的webcocket服务器的功能
	7.还需要实现握手协议，检测第一次建立连接 授信
*/


func main(){
	host := config.File.MustValue("gate_server","host","127.0.0.1")
	port := config.File.MustValue("gate_server","port","8004")
	s := net.NewServer(host+":"+port)
	gate.InitGate()
	s.Router(gate.Router)
	s.Start()
	log.Println("网关服务启动成功")
}