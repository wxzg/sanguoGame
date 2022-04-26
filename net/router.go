package net

import (
	"log"
	"strings"
)

// HandlerFunc 我们具体路由的业务处理函数
type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

type group struct {
	prefix string
	handlerMap map[string]HandlerFunc
}

// AddRouter 其实应该叫addHandlerFunc 就是给对应路由添加业务处理函数
func (g *group) AddRouter(name string,handlerFunc HandlerFunc)  {
	g.handlerMap[name] = handlerFunc
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix: prefix,
		handlerMap: make(map[string]HandlerFunc),
	}
	r.group = append(r.group,g)
	return g
}

//exec方法就是在找到对应的路由组别以后，根据name匹配，将匹配到的具体路由的处理函数拿来处理
func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h := g.handlerMap[name]
	if h != nil {
		h(req,rsp)
	}else{

		gateh := g.handlerMap["*"]
		if gateh != nil {
			gateh(req, rsp)
		}else {
			log.Println("路由不匹配")
		}
	}

}

/*
	设计逻辑是完整的路由有很多分组，比如amount、hero，然后每个分组下面又有各种子路由对应相应的业务处理
*/
type Router struct {
	group []*group
}

func NewRouter() *Router  {
	return &Router{}
}

// Run 请求来了我要跑相应的路由，就到这里，先做相应的处理，根据前缀prefix找到对应的组
func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	//req.Body.Name 路径 登录业务 account.login （account组标识）login 路由标识
	strs := strings.Split(req.Body.Name,".")

	prefix := ""
	name := ""

	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}

	for _,g := range r.group{
		if g.prefix == prefix {
			g.exec(name,req,rsp)
		}else if g.prefix == "*"{
			g.exec(name,req,rsp)
		}
	}
}