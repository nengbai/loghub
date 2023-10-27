package main

import (
	"fmt"
	"log"
	"loghub/src/logmgr"
	"loghub/src/rabbitmq"
	"os"
	"os/signal"
)

func main() {
	var (
		logstr       string = "Fatal"
		queueName    string = "TopicDemo"
		routingKey   string = "Rabbit"
		exchangeName string = "TopicExchange" // RabbitExchange
		exchangeType string = "topic"         //"",fanout,direct,topic
		consumerName string = "consumer-01"   //
	)
	queueExch := rabbitmq.QueueExchange{
		QuName: queueName,
		RtKey:  routingKey,
		ExName: exchangeName,
		ExType: exchangeType,
	}
	rq := rabbitmq.NewRabbitmqMessage(logstr, &queueExch)
	// 验证链接是否正常
	defer rq.MQClose()
	if rq.Channel == nil {
		rq.InitRabbitmq()
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err := rq.Channel.Qos(1, 0, true)
	if err != nil {
		fmt.Printf("获取消费通道异常:%s \n", err)
	}
	msgList, err := rq.Channel.Consume(
		rq.QueueName, // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	logmgr.FailOnError(err, "Rabbitmq Consumer Failure")
	go func() {
		for msg := range msgList {
			// 处理数据
			log.Printf(" [x] %s", string(msg.Body))
			//err = rabbitmq.Receiver.Consumer(msg.Body)
			/**err = msg.Ack(true)
			if err != nil {
				fmt.Printf("确认消息未完成异常:%s \n", err)
				return
			} else {
				// 确认消息,必须为false
				fmt.Printf("确认消息完成异常:%s \n", err)
				return
			}
			**/
		}
	}()
	signals := make(chan os.Signal, 1)
	select {
	case <-signals:
		signal.Notify(signals, os.Interrupt)
	}
}
