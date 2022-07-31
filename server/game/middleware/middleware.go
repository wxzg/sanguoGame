package middleware

import (
	"fmt"
	"log"
	"sanguoServer/net"
	"sanguoServer/utils"
)

func CheckRole() net.MiddlewareFunc  {
	return func(next net.HandlerFunc) net.HandlerFunc {
		return func(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
			_, err := req.Conn.GetProperty("role")
			if err != nil {
				rsp.Body.Code = utils.RoleNotInConnect
				return
			}
			log.Println("角色检测通过")
			next(req, rsp)
		}
	}
}


func Log() net.MiddlewareFunc {
	return func(next net.HandlerFunc) net.HandlerFunc {
		return func(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
			log.Println("请求路由",req.Body.Name)
			log.Println("请求数据",fmt.Sprintf("%v",req.Body.Msg))
			next(req, rsp)
		}
	}
}
