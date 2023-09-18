package conselog

import (
	"fmt"
	"loghub/src/logmgr"
	"time"
)

type ConsoleLog struct {
	LogLv logmgr.LogLevel
}

func NewConsoleLog(logstr string) *ConsoleLog {
	loglv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &ConsoleLog{LogLv: loglv}
}

func (c ConsoleLog) EnableLog(logLevel logmgr.LogLevel) bool {
	return logLevel >= c.LogLv

}
func (c ConsoleLog) log(lv logmgr.LogLevel, format string, a ...any) {
	if c.EnableLog(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now().Format("2006/01/02 15:04:05")
		funcName, fileName, lineNo := logmgr.GetFileInfo(2)
		fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(c.LogLv), fileName, funcName, lineNo, msg)
	}

}

func (c ConsoleLog) Debug(format string, a ...any) {
	c.log(logmgr.DEBUG, format, a...)

}

func (c ConsoleLog) Trace(format string, a ...any) {
	c.log(logmgr.TRACE, format, a...)
}

func (c ConsoleLog) Info(format string, a ...any) {
	c.log(logmgr.INFO, format, a...)
}

func (c ConsoleLog) Warning(format string, a ...any) {
	c.log(logmgr.WARNING, format, a...)
}

func (c ConsoleLog) Error(format string, a ...any) {
	c.log(logmgr.ERROR, format, a...)
}

func (c ConsoleLog) Fatal(format string, a ...any) {
	c.log(logmgr.FATAL, format, a...)
}
