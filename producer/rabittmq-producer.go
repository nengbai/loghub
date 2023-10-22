package main

import (
	"fmt"
	"loghub/src/logmgr"
	"loghub/src/rabbitmq"
)

func main() {
	logstr := "Fatal"
	var (
		queueName    string = "RabbtDemo"
		routingKey   string = "Rabbit"
		exchangeName string = "RabbitExchange"
		exchangeType string = "fanout"
	)
	queueExch := rabbitmq.QueueExchange{
		QuName: queueName,
		RtKey:  routingKey,
		ExName: exchangeName,
		ExType: exchangeType,
	}
	rq := rabbitmq.NewRabbitmqMessage(logstr, &queueExch)
	lv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
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
