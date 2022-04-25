package web

import (
	"github.com/gin-gonic/gin"
	"sanguoServer/db"
	"sanguoServer/server/web/controller"
	"sanguoServer/server/web/middleware"
)

func Init(router *gin.Engine){
	db.DBInit()
	initRouter(router)
}

func initRouter(router *gin.Engine) {
	router.Use(middleware.Cors())
	router.Any("/account/register", controller.DefaultAccountController.Register)
}
