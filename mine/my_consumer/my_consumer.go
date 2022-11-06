package my_consumer

import (
	"example.com/demo/mailman/consumer"
	"fmt"
)

func Setup() {
	consumer1 := consumer.NewConsumer()
	consumer1.Register("test1", func(msg []byte) {
		fmt.Println("i have received! msg = " + string(msg))
	})
	consumer1.Start()
}
