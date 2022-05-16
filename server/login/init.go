package login

import (
	"sanguoServer/db"
	"sanguoServer/net"
	"sanguoServer/server/login/controller"
)

// 新建一个login路由
var Router = net.NewRouter()

func Init()  {
	//初始化数据库
	db.DBInit()
	//还有别的初始化方法
	initRouter()
}

func initRouter()  {
	controller.DefaultAccount.Router(Router)
}