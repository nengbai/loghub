package main

import (
	"loghub/src/filelog"
	"loghub/src/logmgr"
	"time"
)

var mylog logmgr.Mgrloger

func main() {
	filePath := "./"
	fileName := "demo.log"
	errFileName := "demo.log.err"

	mylog = filelog.NewFileLog("DEBUG", filePath, fileName, errFileName, 10*1024*1024)
	//mylog := conselog.NewConsoleLog("ERROR")
	for {
		//fmt.Println("------------------")
		mylog.Debug("This is Debug log:%s,%s", errFileName, fileName)
		mylog.Trace("This is Trace log:%s,%s", errFileName, fileName)
		mylog.Info("This is Info log:%s,%s", errFileName, fileName)
		mylog.Warning("This is warning log:%s,%s", errFileName, fileName)
		mylog.Error("This is Error log:%s,%s", errFileName, fileName)
		mylog.Fatal("This is Fatal log:%s,%s", errFileName, fileName)
		time.Sleep(time.Second)
	}
}
