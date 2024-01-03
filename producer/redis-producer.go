package main

import (
	"fmt"
	"loghub/src/logmgr"
	"loghub/src/redis"
)

func main() {
	logstr := "Fatal"
	k := redis.NewRedisMessage(logstr)
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
