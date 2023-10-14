package main

import (
	"fmt"
	"loghub/src/logmgr"
	"loghub/src/rabbitmq"
)

func main() {
	logstr := "Fatal"
	topic := "test-demo"
	delayTopic := "delay-demo"
	var (
		queueName    string = ""
		routingKey   string = ""
		exchangeName string = ""
		exchangeType string = ""
	)
	queueExch := rabbitmq.QueueExchange{
		QuName: queueName,
		RtKey:  routingKey,
		ExName: exchangeName,
		ExType: exchangeType,
	}
	rq := rabbitmq.NewRabbitmqMessage(logstr, topic, delayTopic, &queueExch)
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
