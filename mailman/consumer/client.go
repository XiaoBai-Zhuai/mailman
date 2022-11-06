package consumer

import (
	"encoding/json"
	"example.com/demo/mailman/entity"
	"fmt"
	"net"
)

type Consumer interface {
	Register(topic string, handler handler)
}

type handler func(msg []byte)

// Register topic: handler
type Register map[string]*handlerDomain

type handlerDomain struct {
	function handler
}

type consumerDomain struct {
	register Register `json:"register"`
}

// 保证单例
var consumer *consumerDomain

func NewConsumer() *consumerDomain {
	if consumer == nil {
		registers := make(map[string]*handlerDomain)
		consumer = &consumerDomain{
			register: registers,
		}
	}
	return consumer
}

func (c *consumerDomain) Register(topic string, handler handler) {
	if _, ok := c.register[topic]; ok {
		panic("topic is exist")
	} else {
		consumer.register[topic] = &handlerDomain{
			function: handler,
		}
	}
}

func (c *consumerDomain) Start() {
	// 注册一下
	// 监听消息
	c.listenMsg(c.registerConsumers())
}

func (c *consumerDomain) registerConsumers() (address string) {
	// 起一个tcp客户端
	conn, err := net.Dial(entity.NetWork, "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	// 先给服务端发送目前注册的所有topic
	data, _ := json.Marshal(c.register)
	consumerMsg := &entity.ClientMsg{
		Msg: entity.MsgBody{
			Message: data,
		},
		ClientType: entity.Consumer,
	}
	byteList, _ := json.Marshal(consumerMsg)
	fmt.Printf("消费者发送消息: %+v\n", string(byteList))
	_, err = conn.Write(byteList)
	if err != nil {
		panic(err)
	}
	address = conn.LocalAddr().String()
	// 要关闭链接，否则服务端会一直阻塞着
	conn.Close()
	return address
}

func (c *consumerDomain) listenMsg(address string) {
	// 起一个tcp客户端监听消息
	server, err := net.Listen(entity.NetWork, address)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
		}
		var bytes = make([]byte, 255)
		n, err := conn.Read(bytes)
		fmt.Printf("消费者接收消息: %+v\n", string(bytes[:n]))
		if err != nil {
			panic(err)
		}
		var msgBody entity.MsgBody
		_ = json.Unmarshal(bytes[:n], &msgBody)
		fmt.Println(c.register)
		if h, ok := c.register[msgBody.Topic]; ok {
			fmt.Println("handler...")
			// 异步处理，不阻塞处理其他消息
			go h.function(msgBody.Message)
		}
	}
}
