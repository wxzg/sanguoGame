package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	config "sanguoServer/conf"
	"sanguoServer/server/web"
	"time"
)

func main(){
	host := config.File.MustValue("web_server","host","127.0.0.1")
	port := config.File.MustValue("web_server","port","8088")

	//启动一个gin的router
	router := gin.Default()
	//路由
	web.Init(router)
	s := &http.Server{
		Addr:		host+":"+port,
		Handler: 	router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil{
		log.Println("web服务器启动失败",err)
	}
}