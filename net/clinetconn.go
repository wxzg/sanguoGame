package net

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"log"
	"sanguoServer/utils"
	"sync"
	"time"
)

type syncCtx struct {
	//Goroutine 的上下文，包含 Goroutine 的运行状态、环境、现场等信息
	ctx context.Context
	cancel context.CancelFunc
	outChan chan *RspBody
}

func NewSyncCtx() *syncCtx  {
	ctx,cancel := context.WithTimeout(context.Background(),15*time.Second)
	return &syncCtx{
		ctx:ctx,
		cancel: cancel,
		outChan: make(chan *RspBody),
	}
}

func (s *syncCtx) wait () *RspBody {
	select {
	case msg := <- s.outChan:
		return msg
	case <- s.ctx.Done():
		log.Println("代理服务响应消息超时了...")
		return nil
	}
}


// 这里是客户端的ws请求
type ClientConn struct{
	wsConn 			*websocket.Conn							//websocket连接，连接网关和客户端的
	isClosed 		bool									//连接是否关闭
	property 		map[string]interface{}					//ws连接上的属性
	propertyLock  	sync.RWMutex							//属性锁
	Seq				int64									//序列号
	handshake		bool									//是否握手成功
	handshakeChan 	chan bool								//握手信号，用来告诉握手是否成功的
	onPush 			func(conn *ClientConn,body *RspBody)    //返回消息给客户端
	onClose    		func(conn *ClientConn)					//关闭ws连接
	syncCtxMap 		map[int64]*syncCtx						//goroutine的上下文
	syncCtxLock 	sync.RWMutex							//上下文锁
}

//Start 一直不停的接收消息
func (c *ClientConn) Start() bool {
	//等待握手的消息返回
	c.handshake = false
	go c.wsReadLoop()
	return c.waitHandShake()
}

func (c *ClientConn)waitHandShake() bool {
	if !c.handshake{

		//防止超时
		ctx, cancel := context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		//等待握手的成功 等待握手的消息
		select {
		case _ = <- c.handshakeChan:
			log.Println("握手成功")
			return true
		case  <- ctx.Done():
			log.Println("握手超时了")
			return false
		}
	}
	return true
}

func (c *ClientConn) wsReadLoop() {

	//注意这里的逻辑都是和服务端的mysql交互的

	//处理一下panic
	defer func() {
		if err := recover(); err != nil{
			log.Println("clientConn ")
			c.Close()
		}
	}()

	for {
		// 准备接收消息
		wc := c.wsConn
		//接收消息
		_, data, err  := wc.ReadMessage()
		if err != nil {
			log.Println("clientconn - 消息接收失败：", err)
			break
		}

		//收到消息以后是json格式的压缩数据，先解压
		data, err = utils.UnZip(data)
		if err != nil{
			log.Println("解压数据出错，非法格式：",err)
			continue
		}

		//解密
		secretkey, err := c.GetProperty("secretKey")
		if err == nil{
			//有加密
			key := secretkey.(string)
			//数据解密
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式有误，解密失败:",err)
			} else {
				data = d
			}
		}

		//将JSON转化为结构体
		body := &RspBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误，非法格式:", err)
		} else {
			//握手 别的一些请求
			//如果是第一个请求，也就是seq为0，注意这个地方是和服务端交互的，是收到的服务端握手信息
			if body.Seq == 0 {
				//并且是握手信息
				if body.Name == HandshakeMsg {
					//获取秘钥
					hs := &Handshake{}
					mapstructure.Decode(body.Msg,hs)
					if hs.Key != "" {
						c.SetProperty("secretKey",hs.Key)
					}else{
						c.RemoveProperty("secretKey")
					}
					//握手成功
					c.handshake = true
					c.handshakeChan <- true
				}else{
					if c.onPush != nil {
						//这里相当于把服务端返还的信息push给前端
						c.onPush(c,body)
					}
				}
			}else{
				c.syncCtxLock.RLock()
				ctx,ok := c.syncCtxMap[body.Seq]
				c.syncCtxLock.RUnlock()
				if ok {
					ctx.outChan <- body
				}else{
					log.Println("no seq syncCtx find")
				}
			}
		}
	}
	c.Close()
}

func (c *ClientConn) write(body interface{}) error{
	data ,err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return err
	}
	//发送给服务端的就不加密了
	//secretKey,err := c.GetProperty("secretKey")
	//if err == nil {
	//	//有加密
	//	key := secretKey.(string)
	//	//数据做加密
	//	data, err = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	//	if err != nil {
	//		log.Println("加密失败",err)
	//		return err
	//	}
	//}
	//压缩
	if data, err := utils.Zip(data); err == nil{
		err := c.wsConn.WriteMessage(websocket.BinaryMessage,data)
		if err != nil {
			log.Println("写数据失败",err)
			return err
		}
	}else {
		log.Println("压缩数据失败",err)
		return err
	}
	return nil
}

func (c *ClientConn) Close(){
	_ = c.wsConn.Close()
}

func NewClientConn(wsConn *websocket.Conn) *ClientConn{
	return &ClientConn{
		wsConn: wsConn,
		handshakeChan: make(chan bool),
		Seq: 0,
		isClosed: false,
		property: make(map[string]interface{}),
		syncCtxMap: map[int64]*syncCtx{},
	}
}

func (c *ClientConn) SetProperty(key string, value interface{}){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

func (c *ClientConn) GetProperty(key string) (interface{}, error){
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (c *ClientConn) RemoveProperty(key string){
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property,key)
}

func (c *ClientConn) Addr() string {
	return c.wsConn.RemoteAddr().String()
}

// 这个函数可以把响应的消息返还给服务端
func (c *ClientConn) Push(name string, data interface{}){
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	c.write(rsp.Body)
}

func (c *ClientConn) SetOnPush(hook func(conn *ClientConn,body *RspBody)) {
	c.onPush = hook
}

func (c *ClientConn) Send(name string, msg interface{}) (*RspBody, error) {
	//把请求 发送给 代理服务器 登录服务器，等待返回
	c.Seq += 1
	seq := c.Seq
	//新建一个上下文，表示这一次连接的上下文环境
	sc := NewSyncCtx()
	c.syncCtxLock.Lock()
	//把上下文存起来以便下次同样的连接还可以用
	c.syncCtxMap[seq] = sc
	c.syncCtxLock.Unlock()
	rsp := &RspBody{Name: name,Seq: seq,Code: utils.OK}
	//req请求
	req := &ReqBody{Seq: seq,Name: name,Msg: msg}
	log.Println(req)
	//将请求发给服务端
	err := c.write(req)

	if err != nil {
		sc.cancel()
	}else{
		r := sc.wait()
		if r == nil {
			rsp.Code = utils.ProxyConnectError
		}else{
			rsp = r
		}
	}
	c.syncCtxLock.Lock()
	delete(c.syncCtxMap,seq)
	c.syncCtxLock.Unlock()
	return rsp, nil
}
