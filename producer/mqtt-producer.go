package main

import (
	"fmt"
	"loghub/src/hivemq"
	"loghub/src/logmgr"
)

func main() {
	logstr := "Fatal"
	mt := hivemq.NewMQTTMessage(logstr)
	if mt.Client == nil {
		mt.InitHivemq()
	}
	// 验证链接是否正常
	lv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	for {
		if lv < logmgr.ERROR {
			mt.Debug("This is Debug log")
			mt.Error("This is Errors log")
			mt.Info("This is Info log")
			mt.Warning("This is Warning log")
		} else if lv <= logmgr.FATAL {
			mt.Error("This is Errors log")
			mt.Fatal("This is Fatal log")
		}
	}

}
