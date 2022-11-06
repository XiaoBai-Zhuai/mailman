package main

import (
	"bufio"
	"encoding/json"
	"example.com/demo/mailman/consumer"
	"example.com/demo/mailman/entity"
	"fmt"
	"net"
	"sync"
)

// topic - clients(ip:port)
var consumers = make(map[string][]string)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listen.Accept() // 监听客户端的连接请求
		if err != nil {
			fmt.Println("Accept() failed, err: ", err)
			continue
		}
		go process(conn) // 启动一个goroutine来处理客户端的连接请求
	}
}

// 处理转发生产者生产消息、消费者注册
func process(conn net.Conn) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	reader := bufio.NewReader(conn)
	var bytes = make([]byte, 255)
	n, err := reader.Read(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("收到消息: %+v\n len(bytes) = %+v\n", string(bytes[:n]), len(bytes))
	var clientMsg entity.ClientMsg
	_ = json.Unmarshal(bytes[:n], &clientMsg)
	switch clientMsg.ClientType {
	case entity.Consumer:
		fmt.Println("processClientRegister...")
		// 处理消费者注册
		err = processClientRegister(clientMsg.Msg.Message, conn.RemoteAddr().String())
		if err != nil {
			panic(err)
		}
	case entity.Producer:
		fmt.Println("processProducerMsg...")
		// 转发生产者消息到消费者
		err = processProducerMsg(clientMsg.Msg)
	}

}

// 处理消费者客户端的注册
func processClientRegister(bytes []byte, address string) error {
	consumerList := make(consumer.Register)
	err := json.Unmarshal(bytes, &consumerList)
	if err != nil {
		return err
	}
	fmt.Println("address = " + address)
	for topic, _ := range consumerList {
		consumers[topic] = append(consumers[topic], address)
	}
	return nil
}

func processProducerMsg(msg entity.MsgBody) error {
	topic := msg.Topic
	bytes, _ := json.Marshal(msg)
	consumerAddressList := consumers[topic]
	fmt.Println(consumerAddressList)
	wg := sync.WaitGroup{}
	wg.Add(len(consumerAddressList))
	// todo 如果有些消费者进程停止了，即连接不上了，需要从消费者列表里清除
	for _, address := range consumerAddressList {
		go func() {
			conn, err := net.Dial(entity.NetWork, address)
			defer func() {
				wg.Done()
				conn.Close()
			}()
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = conn.Write(bytes)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()
	}
	return nil
}
