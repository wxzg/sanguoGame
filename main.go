package main

import (
	"log"
	config "sanguoServer/conf"
	"sanguoServer/login"
	"sanguoServer/net"
)

func main(){
	host := config.File.MustValue("login_server","host","127.0.0.1")
	port := config.File.MustValue("login_server","port","8003")
	s := net.NewServer(host+":"+port)
	login.Init()
	s.Router(login.Router)
	s.Start()
	log.Println("登录服务启动成功")
}
