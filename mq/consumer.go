package mq

import "log"

var done chan bool

// 开始监听队列，获取信息
func StartConsume(qName, cName string, callBack func(msg []byte) bool) {
	// 1. 用过channel.Consume获取消息信道
	msgs, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 2. 循环获取队列的消息

	go func() {
		for msg := range msgs {
			// 3.调用callback方法处理新的消息
			procssSuc := callBack(msg.Body)
			if !procssSuc {
				// TODO：将任务写到另外一个队列，用于异常情况的重试
			}
		}
	}()

	// done没有新的消息过来，会一直发生阻塞
	<-done

	// 关闭rabbitmq
	channel.Close()

}
