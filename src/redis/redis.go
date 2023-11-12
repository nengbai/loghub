package redis

import (
	"fmt"
	"loghub/src/logmgr"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

type Producer interface {
	Publish() string
}

// 定义接收者接口
type Receiver interface {
	Subscribe([]byte) error
}
type RedisMessage struct {
	LogLv   logmgr.LogLevel
	Options *redis.Options
	Client  *redis.Client
}

func NewRedisMessage(logstr string) *RedisMessage {
	loglv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fl := &RedisMessage{
		LogLv: loglv,
	}
	err = fl.InitRedis()
	logmgr.FailOnError(err, "Redis初始化失败:")
	return fl
}

func (r *RedisMessage) InitRedis() error {
	// 1.读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	db := viper.GetInt("redis.db")
	// Redis用户的密码
	password := viper.GetString("redis.password")
	PoolSize := viper.GetInt("redis.poolSize")
	MinIdleConns := viper.GetInt("redis.minIdleConn")
	// RabbitMQ Broker 的ip地址
	r.Options = &redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     password,
		DB:           db,
		PoolSize:     PoolSize,
		MinIdleConns: MinIdleConns,
	}
	r.Client = redis.NewClient(r.Options)
	return nil
}

func (r *RedisMessage) log(lv logmgr.LogLevel, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now().Format("2006/01/02 15:04:05")
	funcName, fileName, lineNo := logmgr.GetFileInfo(3)
	//  发送消息
	// Trap SIGINT to trigger a graceful shutdown.
	// signals := make(chan os.Signal, 1)
	// fmt.Println("Log is:", r.QueueName)
	// 用于检查队列是否存在,已经存在不需要重复声明

	if r.LogLv < logmgr.ERROR {
		message := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(r.LogLv), fileName, funcName, lineNo, msg)
		fmt.Printf("Error message:%s\n", message)

		// fmt.Println("Error Log rmsg:", rmsg)

		//debugLoop:

	} else if r.LogLv <= logmgr.FATAL {
		message := fmt.Sprintf("[%s][%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(r.LogLv), fileName, funcName, lineNo, msg)
		// fmt.Printf("Fatal message:%s\n", message)

		fmt.Println("Fatal Log rmsg:", message)

	} else {
		fmt.Println("Loglevel set error:")
	}

}

func (r *RedisMessage) Debug(format string, a ...any) {
	r.log(logmgr.DEBUG, format, a...)

}

func (r *RedisMessage) Trace(format string, a ...any) {
	r.log(logmgr.TRACE, format, a...)
}

func (r *RedisMessage) Info(format string, a ...any) {
	r.log(logmgr.INFO, format, a...)
}

func (r *RedisMessage) Warning(format string, a ...any) {
	r.log(logmgr.WARNING, format, a...)
}

func (r *RedisMessage) Error(format string, a ...any) {
	r.log(logmgr.ERROR, format, a...)
}

func (r *RedisMessage) Fatal(format string, a ...any) {
	r.log(logmgr.FATAL, format, a...)
}

func (r *RedisMessage) Consumer(msg []byte) error {

	return nil

}
