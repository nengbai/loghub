package main

import (
	"fmt"
	"loghub/src/kafka"
	"loghub/src/logmgr"
)

func main() {
	logstr := "Fatal"
	topic := "test-demo"
	delayTopic := "delay-demo"
	k := kafka.NewKafkaMessage(logstr, topic, delayTopic)
	lv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	for {
		if lv < logmgr.ERROR {
			k.Debug("This is Debug log")
			k.Error("This is Errors log")
			k.Info("This is Info log")
			k.Warning("This is Warning log")
		} else if lv <= logmgr.FATAL {
			k.Error("This is Errors log")
			k.Fatal("This is Fatal log")
		}
	}

}
