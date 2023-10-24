package main

import (
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
	)
	queueExch := rabbitmq.QueueExchange{
		QuName: queueName,
		RtKey:  routingKey,
		ExName: exchangeName,
		ExType: exchangeType,
	}
	rq := rabbitmq.NewRabbitmqMessage(logstr, &queueExch)

	msgs, err := rq.Channel.Consume(
		rq.QueueName, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	logmgr.FailOnError(err, "Rabbitmq Consumer Failure")

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", string(d.Body))
		}
	}()
	signals := make(chan os.Signal, 1)
	select {
	case <-signals:
		signal.Notify(signals, os.Interrupt)
	}
}
