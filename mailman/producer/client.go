package producer

import (
	"encoding/json"
	"example.com/demo/mailman/entity"
	"fmt"
	"net"
)

type Producer interface {
	Send(msg []byte) error
}

type producerDomain struct {
	topic string
}

func NewProducer(topic string) *producerDomain {
	return &producerDomain{
		topic: topic,
	}
}

func (p *producerDomain) Send(msg []byte) error {
	message := &entity.ClientMsg{
		Msg: entity.MsgBody{
			Message: msg,
			Topic:   p.topic,
		},
		ClientType: entity.Producer,
	}
	// 起一个tcp客户端
	conn, err := net.Dial(entity.NetWork, "127.0.0.1:9999")
	if err != nil {
		return err
	}
	defer conn.Close()
	bytes, err := json.Marshal(message)
	fmt.Printf("生产者发送消息: %+v\n", string(bytes))
	if err != nil {
		return err
	}
	_, err = conn.Write(bytes) // 发送数据
	if err != nil {
		return err
	}
	return nil
}
