package net

import "github.com/gorilla/websocket"

type ClientConn struct{
	wsConn *websocket.Conn
}

func NewClientConn(wsConn *websocket.Conn) *ClientConn{
	return &ClientConn{
		wsConn: wsConn,
	}
}
