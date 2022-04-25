package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"log"
	"sanguoServer/utils"
	"sync"
)

//websocket服务
type wsServer struct {
	wsConn *websocket.Conn
	router *Router
	outChan chan *WsMsgRsp
	Seq			int64
	property 	map[string]interface{}
	propertyLock sync.RWMutex
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn: wsConn,
		outChan: make(chan *WsMsgRsp, 1000),
		property: make(map[string]interface{}),
		Seq: 0,
	}
}

func (w *wsServer) Router(router *Router)  {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}){
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

func (w *wsServer) GetProperty(key string) (interface{}, error){
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}
func (w *wsServer) RemoveProperty(key string){
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property,key)
}
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}
func (w *wsServer) Push(name string, data interface{}){
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}

//通道一旦建立 ，那么 收发消息 就得要一直监听才行
func (w *wsServer) Start() {
	//启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) writeMsgLoop() {
	for  {
		select {
		case msg := <- w.outChan:
			w.Write(msg)
		}
	}
}

func (w *wsServer) Write(msg *WsMsgRsp) {
	data ,err := json.Marshal(msg.Body)
	if err != nil {
		log.Println(err)
	}
	secretKey,err := w.GetProperty("secretKey")
	if err == nil {
		//有加密
		key := secretKey.(string)
		//数据做加密
		data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	}
	//压缩
	if data, err := utils.Zip(data); err == nil{
		_ = w.wsConn.WriteMessage(websocket.BinaryMessage,data)
	}
}

func (w *wsServer) readMsgLoop() {
	//先读到客户端 发送过来的数据，然后 进行处理，然后在回消息
	//经过路由 实际处理程序
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	for  {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误:",err)
			break
		}
		//收到消息 解析消息 前端发送过来的消息 就是json格式
		//1. data 解压 unzip
		data,err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，非法格式：",err)
			continue
		}
		//2. 前端的消息 加密消息 进行解密
		secretKey,err := w.GetProperty("secretKey")
		if err == nil {
			//有加密
			key := secretKey.(string)
			//客户端传过来的数据是加密的 需要解密
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式有误，解密失败:",err)
				//出错后 发起握手
				w.Handshake()
			}else{
				data = d
			}
		}
		//3. data 转为body
		body := &ReqBody{}
		err = json.Unmarshal(data,body)
		fmt.Printf("body-Type:%T, value:%v",body,body)
		fmt.Println()
		if err != nil {
			log.Println("数据格式有误，非法格式:",err)
		}else{
			// 获取到前端传递的数据了，拿上这些数据 去具体的业务进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			w.router.Run(req,rsp)
			w.outChan <- rsp
		}
	}
	w.Close()
}

func (w *wsServer) Close()  {
	_ = w.wsConn.Close()
}
const HandshakeMsg = "handshake"

// Handshake 当游戏客户端 发送请求的时候 会先进行握手协议
//后端会发送对应的加密key给客户端
//客户端在发送数据的时候 就会使用此key进行加密处理
func (w *wsServer) Handshake()  {
	secretKey := ""
	key,err := w.GetProperty("secretKey")
	if err == nil {
		secretKey = key.(string)
	}else{
		secretKey = utils.RandSeq(16)
	}
	handshake := &Handshake{Key: secretKey}

	body := &RspBody{Name: HandshakeMsg,Msg: handshake}

	if data, err := json.Marshal(body); err == nil {
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			w.RemoveProperty("secretKey")
		}
		if data,err := utils.Zip(data); err == nil{
			_ = w.wsConn.WriteMessage(websocket.BinaryMessage,data)
		}
	}
}