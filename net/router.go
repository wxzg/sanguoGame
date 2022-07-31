package net

import (
	"log"
	"strings"
)

// HandlerFunc 我们具体路由的业务处理函数
type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)
type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc
type group struct {
	prefix string
	handlerMap map[string]HandlerFunc
	middlewareMap map[string][]MiddlewareFunc
	milldewares []MiddlewareFunc
}

// AddRouter 其实应该叫addHandlerFunc 就是给对应路由添加业务处理函数
func (g *group) AddRouter(name string,handlerFunc HandlerFunc, middlewares... MiddlewareFunc)  {
	g.handlerMap[name] = handlerFunc
	g.middlewareMap[name] = middlewares
}

//Group 根据前缀给路由新建一个路由组，返回一个Group，但这是router上的方法
func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix: prefix,
		handlerMap: make(map[string]HandlerFunc),
		middlewareMap: map[string][]MiddlewareFunc{},
		milldewares: make([]MiddlewareFunc,0),
	}
	r.group = append(r.group,g)
	return g
}

func (g *group) Use (middlewares... MiddlewareFunc) {
	g.milldewares = append(g.milldewares, middlewares...)
}

//exec方法就是在找到对应的路由组别以后，根据name匹配，将匹配到的具体路由的处理函数拿来处理
func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	//根据具体路由（如login）匹配业务处理函数
	h, ok := g.handlerMap[name]
	if ok {
		//先把中间件走一遍,这更是全局中间件比如account
		for i := 0; i < len(g.milldewares); i++ {
			h = g.milldewares[i](h)
		}
		//这个是具体业务的中间件比如login
		m, ok := g.middlewareMap[name]
		if ok {
			for i := 0; i < len(m); i++ {
				h = m[i](h)
			}
		}
		//匹配到了就执行
		h(req,rsp)
	}else{
		//如果是匹配不到再试试匹配“*”
		gateh, ok := g.handlerMap["*"]
		if ok {
			//如果匹配到了表示网关那边有处理函数，执行
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

// NewRouter 新建一个路由
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
	//这里主要是看看分组能不能匹配上，比如这里就是要匹配到account组，然后将login传过去
	for _,g := range r.group{
		if g.prefix == prefix {
			// 这里进行具体的业务处理
			g.exec(name,req,rsp)
		}else if g.prefix == "*"{
			//如果是网关“*”来了也要能通过
			g.exec(name,req,rsp)
		}
	}
}