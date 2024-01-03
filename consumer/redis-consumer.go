package main

import (
	"loghub/src/redis"
	"os"
	"os/signal"
)

func main() {
	var (
		logstr      string = "Fatal"
		ChannelName string = "mylog" //
	)

	rd := redis.NewRedisMessage(logstr)
	// 验证链接是否正常

	if rd.Client == nil {
		rd.InitRedis()
	}
	// 获取消费通道,确保Redis一个一个发送消息

	//logmgr.FailOnError(err, "Rabbitmq Consumer Failure")
	go func() {
		redis.Subscriber(rd.Client, ChannelName)
	}()
	signals := make(chan os.Signal, 1)
	select {
	case <-signals:
		signal.Notify(signals, os.Interrupt)
	}
}
