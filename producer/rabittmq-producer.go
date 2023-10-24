package main

import (
	"fmt"
	"loghub/src/logmgr"
	"loghub/src/rabbitmq"

	"github.com/streadway/amqp"
)

func main() {
	logstr := "Fatal"
	var (
		queueName    string = "TopicDemo"
		routingKey   string = "Rabbit"
		exchangeName string = "TopicExchange" // RabbitExchange
		exchangeType string = "topic"         //"",fanout,direct,topic
		reliable     bool   = false           //
	)
	queueExch := rabbitmq.QueueExchange{
		QuName: queueName,
		RtKey:  routingKey,
		ExName: exchangeName,
		ExType: exchangeType,
	}
	rq := rabbitmq.NewRabbitmqMessage(logstr, &queueExch)
	lv, err := logmgr.ParseLoglevel(logstr)
	logmgr.FailOnError(err, "Failed to Parse Log Level")
	// Reliable publisher confirms require confirm.select support from the connection.
	if reliable {
		fmt.Println("enabling publishing confirms.")
		if err := rq.Channel.Confirm(false); err != nil {
			fmt.Printf("Channel could not be put into confirm mode: %s\n", err)
		}
		confirms := rq.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer rabbitmq.ConfirmOne(confirms)
	}
	for {
		if lv < logmgr.ERROR {
			rq.Debug("This is Debug log")
			rq.Error("This is Errors log")
			rq.Info("This is Info log")
			rq.Warning("This is Warning log")
		} else if lv <= logmgr.FATAL {
			rq.Error("This is Errors log")
			rq.Fatal("This is Fatal log")
		}
	}
}
