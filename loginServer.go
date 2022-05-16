package main

import (
	"log"
	config "sanguoServer/conf"
	"sanguoServer/net"
	"sanguoServer/server/login"
)

func main(){
	host := config.File.MustValue("login_server","host","127.0.0.1")
	port := config.File.MustValue("login_server","port","8003")
	//新建一个webSocket服务
	s := net.NewServer(host+":"+port)
	//初始化数据库和登录路由
	login.Init()
	//将登录路由放到ws服务中
	s.Router(login.Router)
	//启动服务 ws开始监听
	s.Start()
	log.Println("登录服务启动成功")
}

