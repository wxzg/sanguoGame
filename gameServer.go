package main

import (
	"log"
	config "sanguoServer/conf"
	"sanguoServer/net"
	"sanguoServer/server/game"
)

/**
	1. 登录完成了，创建角色（玩家）
	2. 需要根据用户查询角色，如果没有就创建角色
	3. 初始化
*/
func main(){
	host := config.File.MustValue("game_server","host","127.0.0.1")
	port := config.File.MustValue("game_server","port","8001")
	s := net.NewServer(host+":"+port)
	s.NeedSecret(false)
	game.Init()
	s.Router(game.Router)
	s.Start()
	log.Println("游戏服务启动成功")
}
