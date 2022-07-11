package net

// 这个是前端传来的请求对应的JSON数据
type ReqBody struct {
	Seq     int64		`json:"seq"` 	//序列号
	Name 	string 		`json:"name"` 	//路由 如 account/login
	Msg		interface{}	`json:"msg"`	//数据
	Proxy	string		`json:"proxy"`	//代理 如ws://127.0.0.1:8003
}
//这个是我们要返回给前端的数据格式
type RspBody struct {
	Seq     int64		`json:"seq"` 	//序列
	Name 	string 		`json:"name"`	//路由
	Code	int			`json:"code"`	//状态码 - 如成功为0 失败为1
	Msg		interface{}	`json:"msg"`	//数据
}

// 将websocket连接和请求体封装一下
type WsMsgReq struct {
	Body	*ReqBody
	Conn	WSConn
}

// 封装一下响应体
type WsMsgRsp struct {
	Body*	RspBody
}

// WSConn 接口 - 理解为 request请求 请求会有参数 请求中放参数 取参数
//前端请求过来时会带有一些参数比如token？，我们需要对这些参数进行处理和操作
type WSConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}

// Handshake 握手
type Handshake struct {
	Key string `json:"key"`
}

//心跳检测
type Heartbeat struct {
	CTime	int64 	`json:"ctime"`
	STime	int64 	`json:"stime"`
}