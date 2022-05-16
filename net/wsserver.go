package net

import (
	"encoding/json"
	"errors"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"log"
	"sanguoServer/utils"
	"sync"
)

//websocket服务
type wsServer struct {
	wsConn *websocket.Conn 	//websocket连接
	router *Router			//路由
	outChan chan *WsMsgRsp	//通道 - 用来传输回复的消息
	Seq			int64		//序列号
	property 	map[string]interface{} //用来存储request属性
	propertyLock sync.RWMutex	//存储属性时需要加锁
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

// SetProperty 设置属性
func (w *wsServer) SetProperty(key string, value interface{}){
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

// GetProperty 得到属性
func (w *wsServer) GetProperty(key string) (interface{}, error){
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// RemoveProperty 删除属性
func (w *wsServer) RemoveProperty(key string){
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property,key)
}

// Addr 获取websocket客户端地址
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}

//Push 将消息放到outChan中，然后会走到Write里发送给客户端
func (w *wsServer) Push(name string, data interface{}){
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}

// Start 通道一旦建立 ，那么 收发消息 就得要一直监听才行
func (w *wsServer) Start() {
	//启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) writeMsgLoop() {
	//不断监听是否有数据到了outChan中，一旦有数据过来，就立刻处理然后发送
	for  {
		select {
		case msg := <- w.outChan:
			w.Write(msg)
		}
	}
}

func (w *wsServer) Write(msg *WsMsgRsp) {
	//由于拿到的数据时封装过的wsMsgRep，所以需要先处理一下
	//将Res转换成JSON
	data ,err := json.Marshal(msg.Body)
	if err != nil {
		log.Println(err)
	}
	//取出密钥
	secretKey,err := w.GetProperty("secretKey")
	if err == nil {
		//有加密
		key := secretKey.(string)
		//数据做加密
		data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	}
	//压缩
	if data, err := utils.Zip(data); err == nil{
		//将处理好的数据发送给游戏客户端
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
				//出错后 重新发起握手
				w.Handshake()
			}else{
				//否则d就是解密后的数据
				data = d
			}
		}
		//3. data 转为body ，这里data依旧是JSON数据
		body := &ReqBody{}
		err = json.Unmarshal(data,body)

		if err != nil {
			log.Println("数据格式有误，非法格式:",err)
		}else{
			// 真正获取到前端传递的数据了，拿上这些数据 去具体的业务进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			//拿去给路由处理
			w.router.Run(req,rsp)
			//将rsp数据传给outChan，因为rsp是引用数据，因此那边处理完以后这里会同步更新
			w.outChan <- rsp
		}
	}
	//如果除了循环，则关闭websocket
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
		//如果wsServer中已经有加密key
		secretKey = key.(string)
	}else{
		//如果没有则随机生成一个加密Key
		secretKey = utils.RandSeq(16)
	}
	//新建一个握手，保存对应的加密key
	handshake := &Handshake{Key: secretKey}
	//这里是返给前端的数据，即”handshake“ - 加密Key
	body := &RspBody{Name: HandshakeMsg,Msg: handshake}
	//然后将要返回的数据转成JSON
	if data, err := json.Marshal(body); err == nil {
		//这里保存一下加密Key到wsServer属性中
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			//这里老夫也不懂，我觉得似乎没必要
			w.RemoveProperty("secretKey")
		}
		//将数据打包一下返回给前端
		if data,err := utils.Zip(data); err == nil{
			_ = w.wsConn.WriteMessage(websocket.BinaryMessage,data)
		}
	}
}