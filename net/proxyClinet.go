package net

import (
	"errors"
	"github.com/gorilla/websocket"
	"time"
)

type ProxyClient struct{
	Proxy string
	Conn *ClientConn
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
		if !c.Conn.Start(){
			return errors.New("握手失败")
		}
	}
	return err
}