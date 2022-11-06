package my_producer

import (
	"example.com/demo/mailman/producer"
)

func Setup() {
	producer1 := producer.NewProducer("test1")
	err := producer1.Send([]byte("hello, i'm producer1"))
	if err != nil {
		panic(err)
	}
}
