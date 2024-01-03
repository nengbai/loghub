package main

import (
	"fmt"
	"loghub/src/logmgr"
	"loghub/src/redis"
)

func main() {
	logstr := "Fatal"
	rd := redis.NewRedisMessage(logstr)
	if rd.Client == nil {
		rd.InitRedis()
	}
	// 验证链接是否正常
	pong, err := rd.Client.Ping().Result() // 检查是否连接
	if err != nil {
		panic(err)
	}
	// 连接成功啦
	fmt.Println(pong)
	lv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	for {
		if lv < logmgr.ERROR {
			rd.Debug("This is Debug log")
			rd.Error("This is Errors log")
			rd.Info("This is Info log")
			rd.Warning("This is Warning log")
		} else if lv <= logmgr.FATAL {
			rd.Error("This is Errors log")
			rd.Fatal("This is Fatal log")
		}
	}

}
