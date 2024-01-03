package main

import (
	"fmt"
	"loghub/src/hivemq"
	"os"
	"os/signal"
)

func main() {
	var (
		logstr string = "Fatal"
	)

	mt := hivemq.NewMQTTMessage(logstr)

	if token := mt.Client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
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
