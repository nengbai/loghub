package main

import (
	"loghub/src/hivemq"
	"os"
	"os/signal"
)

func main() {
	var (
		logstr string = "Fatal"
	)

	mt := hivemq.NewMQTTMessage(logstr)

	if mt.Client == nil {
		mt.InitHivemq()
	}

	// 获取消费通道,确保hivemq一个一个发送消息 qos2
	//logmgr.FailOnError(err, "Rabbitmq Consumer Failure")
	go func() {
		hivemq.Subscribe(mt.Client, mt.Topic, byte(mt.QoS))
	}()
	signals := make(chan os.Signal, 1)
	select {
	case <-signals:
		signal.Notify(signals, os.Interrupt)
	}
}
