package main

import (
	"example.com/demo/mine/my_consumer"
	"example.com/demo/mine/my_producer"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		my_consumer.Setup()
	}()
	time.Sleep(3 * time.Second)
	go func() {
		defer wg.Done()
		my_producer.Setup()
	}()
	wg.Wait()

}
