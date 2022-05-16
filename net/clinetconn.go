package net

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

// 这里是客户端的ws请求
type ClientConn struct{
	wsConn *websocket.Conn
	isClosed bool
	property 		map[string]interface{}
	propertyLock  	sync.RWMutex
	Seq				int64
	handshake		bool
	handshakeChan 	chan bool
	onPush func(conn *ClientConn,body *RspBody)
	onClose    		func(conn*ClientConn)
	syncCtxMap map[int64]*syncCtx
	syncCtxLock sync.RWMutex
}

//Start 一直不停的接收消息
func (c *ClientConn) Start() bool {
	//等待握手的消息返回

	c.handshake = false
	go c.wsReadLoop()
	return c.waitHandShake()
}

func (c *ClientConn)waitHandShake() bool {
	//防止超时
	ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	//等待握手的成功 等待握手的消息
	select {
	case _ = <- c.handshakeChan:
		log.Println("握手成功")
		return true
	case  <-ctx.Done():
		log.Println("握手超时了")
		return false
	}
}

func (c *ClientConn) wsReadLoop() {
	for {
		_, data, err := c.wsConn.ReadMessage()
		fmt.Println(data, err)
		//收到握手消息了
		c.handshake = true
		c.handshakeChan <- true
	}
}

type syncCtx struct {
	//Goroutine 的上下文，包含 Goroutine 的运行状态、环境、现场等信息
	ctx context.Context
	cancel context.CancelFunc
	outChan chan *RspBody
}

func NewClientConn(wsConn *websocket.Conn) *ClientConn{
	return &ClientConn{
		wsConn: wsConn,
	}
}
