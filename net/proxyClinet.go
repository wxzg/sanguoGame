package net

import (
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

type ProxyClient struct{
	Proxy string			//服务端地址
	Conn *ClientConn		//客户端连接
}

func NewProxyClient(proxy string) *ProxyClient  {
	return &ProxyClient{
		Proxy: proxy,
	}
}

//Connect 去连接websocket服务端
func (c *ProxyClient) Connect() error{
	//通过Dialer连接websocket服务器
	var dialer = websocket.Dialer{
		Subprotocols:     []string{"p1", "p2"},
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}
	ws, _, err := dialer.Dial(c.Proxy, nil)
	if err == nil{
		//如果没有错误就将websocket客户端丢给定义好的客户端连接
		c.Conn = NewClientConn(ws)
		//启动客户端，相当于把快递站搭起来开始运行了
		if !c.Conn.Start(){
			return errors.New("握手失败")
		}
	}
	return err
}

func (c *ProxyClient) SetOnPush(hook func(conn *ClientConn,body *RspBody)){
	//如果有连接，就把钩子函数发给客户端连接
	if c.Conn != nil{
		c.Conn.SetOnPush(hook)
	}
}

func (c *ProxyClient) Send(name string, msg interface{}) (*RspBody, error){
	if c.Conn != nil {
		return c.Conn.Send(name, msg)
	}

	return nil, errors.New("proxy连接未建立")
}

func (c *ProxyClient) SetProperty (key string, value interface{}){
	if c.Conn != nil {
		c.Conn.SetProperty(key, value)
	}
}
