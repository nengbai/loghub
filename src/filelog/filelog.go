package filelog

import (
	"fmt"
	"loghub/src/logmgr"
	"os"
	"path"
	"syscall"
	"time"
)

type FileLog struct {
	LogLv          logmgr.LogLevel
	FilePath       string
	FileName       string
	FileObj        *os.File
	ErrFileName    string
	ErrFileObj     *os.File
	FilemaxSize    int64
	FileCreateTime time.Time
}

func NewFileLog(logstr, filePath, fileName, errFileName string, fileMaxSize int64) *FileLog {
	loglv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fl := &FileLog{
		LogLv:       loglv,
		FilePath:    filePath,
		FileName:    fileName,
		ErrFileName: errFileName,
		FilemaxSize: fileMaxSize,
	}
	err = fl.InitFile()
	if err != nil {
		panic(err)
	}
	return fl
}

func (f *FileLog) InitFile() error {
	fullFilename := path.Join(f.FilePath, f.FileName)
	fileObj, err := os.OpenFile(fullFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open File Failed:%v\n", err.Error())
		return err
	}
	errFullFilename := path.Join(f.FilePath, f.ErrFileName)
	errFileObj, err := os.OpenFile(errFullFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Open File Failed:%v\n", err.Error())
		return err
	}
	f.FileObj = fileObj
	f.ErrFileObj = errFileObj
	fileCreateTime := time.Now()
	f.FileCreateTime = fileCreateTime
	return nil
}

func (f *FileLog) EnableLog(logLevel logmgr.LogLevel) bool {
	//fmt.Printf("logLevel:%d,f.LogLv:%d\n", logLevel, f.LogLv)
	return logLevel >= f.LogLv

}
func (f *FileLog) CheckFileSize(file *os.File) bool {
	fileObj, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	fileSize := fileObj.Size()
	return fileSize >= f.FilemaxSize
}
func (f *FileLog) CheckTimeStamp(file *os.File) bool {

	fileObj, err := os.Stat(file.Name())
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	stat_t := fileObj.Sys().(*syscall.Stat_t)
	//fmt.Println(timespecToTime(stat_t.Ctimespec))
	fileCTime := logmgr.TimespecToTime(stat_t.Birthtimespec)
	fileDurationTime := int(time.Since(fileCTime).Hours())
	fmt.Println("fileDurationTime:", fileDurationTime)
	return fileDurationTime > 24
}
func (f *FileLog) FileSplit(fo *os.File) (*os.File, error) {
	timestamp := time.Now().Format("20060102150405000")
	fmt.Printf("Get spit file:%v\n", fo)
	fileInfo, err := fo.Stat()
	if err != nil {
		fmt.Println("Open new redo log failed:", err.Error())
		return nil, err
	}
	file := path.Join(f.FilePath, fileInfo.Name())
	bakFile := file + "-bak-" + timestamp
	// 1. 关闭原文件
	fo.Close()
	// 2. 备份源文件,生成新文件
	os.Rename(file, bakFile)
	// 3. 重新打开新日志文件
	fileObj, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Open new redo log failed:", err.Error())
		return nil, err
	}
	// 4. return FileObj
	return fileObj, nil

}

func (f *FileLog) log(lv logmgr.LogLevel, format string, a ...any) {
	if f.EnableLog(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now().Format("2006/01/02 15:04:05")
		funcName, fileName, lineNo := logmgr.GetFileInfo(3)
		if lv < logmgr.ERROR {
			Tinfo := f.CheckTimeStamp(f.FileObj)
			fmt.Printf("Check file timestamp:%v\n", Tinfo)
			if f.CheckFileSize(f.FileObj) {
				newFile, err := f.FileSplit(f.FileObj)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				f.FileObj = newFile

			}
			//fmt.Printf("newfile:%v\n", f.FileObj)
			fmt.Fprintf(f.FileObj, "[%s] [%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(f.LogLv), fileName, funcName, lineNo, msg)

		} else if lv >= logmgr.ERROR {
			if f.CheckFileSize(f.ErrFileObj) {
				newFile, err := f.FileSplit(f.ErrFileObj)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				f.ErrFileObj = newFile
			}
			//fmt.Printf("new err file:%v\n", f.ErrFileObj)
			fmt.Fprintf(f.ErrFileObj, "[%s][%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(f.LogLv), fileName, funcName, lineNo, msg)
		}
	}

}

func (f FileLog) Debug(format string, a ...any) {
	f.log(logmgr.DEBUG, format, a...)

}

func (f FileLog) Trace(format string, a ...any) {
	f.log(logmgr.TRACE, format, a...)
}

func (f FileLog) Info(format string, a ...any) {
	f.log(logmgr.INFO, format, a...)
}

func (f FileLog) Warning(format string, a ...any) {
	f.log(logmgr.WARNING, format, a...)
}

func (f FileLog) Error(format string, a ...any) {
	f.log(logmgr.ERROR, format, a...)
}

func (f FileLog) Fatal(format string, a ...any) {
	f.log(logmgr.FATAL, format, a...)
}

func (f *FileLog) Close() {
	f.FileObj.Close()
	f.ErrFileObj.Close()
}
