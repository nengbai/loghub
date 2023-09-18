package logmgr

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type LogLevel uint16
type Mgrloger interface {
	Debug(format string, a ...any)
	Trace(format string, a ...any)
	Info(format string, a ...any)
	Warning(format string, a ...any)
	Error(format string, a ...any)
	Fatal(format string, a ...any)
}

const (
	UNKOWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

func ParseLoglevel(s string) (LogLevel, error) {
	a := strings.ToLower(s)
	switch a {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("this is'not correct log level")
		return UNKOWN, err
	}
}

func GetLogString(lv LogLevel) string {
	switch lv {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKOWN"
	}
}

func GetFileInfo(skip int) (string, string, int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("runtime.caller run failed!")
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file)
	//fmt.Printf("%v,%s,%d\n", funcName, fileName, lineNo)
	return funcName, fileName, lineNo
}

func GetOpenfile(a string) (*os.File, error) {
	fileObj, err := os.OpenFile(a, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer fileObj.Close()
	fmt.Println("fileObj is:", fileObj)
	return fileObj, nil
}
func TimespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}
