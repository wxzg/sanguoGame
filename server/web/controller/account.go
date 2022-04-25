package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sanguoServer/server/common"
	"sanguoServer/server/web/logic"
	"sanguoServer/server/web/webmodel"
	"sanguoServer/utils"
)

var DefaultAccountController *AccountController = &AccountController{}

type AccountController struct {

}

func (a *AccountController) Register(ctx *gin.Context){
	/**
		1.获取请求参数
		2.根据用户名 查询数据库中该用户是否存在 如果不存在则注册
		3。告诉前端注册成功
	*/
	rq := &webmodel.RegisterReq{}
	err := ctx.ShouldBind(rq)
	if err != nil {
		log.Println("参数格式不合法", err)
		ctx.JSON(http.StatusOK, common.Error(utils.InvalidParam, "参数不合法"))
		return
	}
	//一般web服务 错误格式 自定义
	err = logic.DefaultAccountLogic.Register(rq)
	if err != nil {
		log.Println("注册业务出错",err)
		//   这里的error其实是一个实现了error接口的结构体，可以用类型断言转换成我们定义的MyError类型
		ctx.JSON(http.StatusOK,common.Error(err.(*common.MyError).Code(),err.Error()))
		return
	}
	ctx.JSON(http.StatusOK,common.Success(utils.OK,nil))
}