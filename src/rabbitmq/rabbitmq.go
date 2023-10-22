package rabbitmq

import (
	"fmt"
	"loghub/src/logmgr"
	"sync"
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var mqConn *amqp.Connection
var mqChan *amqp.Channel

// 定义生产者接口
type Producer interface {
	MsgContent() string
}

// 定义接收者接口
type Receiver interface {
	Consumer([]byte) error
}

type RabbitmqMessage struct {
	LogLv        logmgr.LogLevel
	Config       *amqp.Config
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	QueueName    string // 队列名称
	RoutingKey   string // key名称
	ExchangeName string // 交换机名称
	ExchangeType string // 交换机类型
	ProducerList []Producer
	ReceiverList []Receiver
	Mu           sync.RWMutex
}

// 定义队列交换机对象
type QueueExchange struct {
	QuName string // 队列名称
	RtKey  string // key值
	ExName string // 交换机名称
	ExType string // 交换机类型
}

func NewRabbitmqMessage(logstr string, q *QueueExchange) *RabbitmqMessage {
	loglv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fl := &RabbitmqMessage{
		LogLv:        loglv,
		QueueName:    q.QuName,
		RoutingKey:   q.RtKey,
		ExchangeName: q.ExName,
		ExchangeType: q.ExType,
	}
	err = fl.InitRabbitmq()
	if err != nil {
		panic(err)
	}
	return fl
}

func (r *RabbitmqMessage) InitRabbitmq() error {
	// 1.读取配置文件
	var RabbitUrl string
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	user := viper.GetString("rabbitmq.user")
	// RabbitMQ用户的密码
	password := viper.GetString("rabbitmq.password")
	// RabbitMQ Broker 的ip地址
	broker := viper.GetStringSlice("rabbitmq.addrs")
	// RabbitMQ Broker 监听的端口
	//port string = "5672"
	vhost := viper.GetString("rabbitmq.vhost")
	for _, host := range broker {
		RabbitUrl = "amqp://" + user + ":" + password + "@" + host + vhost
		// fmt.Printf("RabbitMQ url is:%s", RabbitUrl)
	}

	// RabbitUrl := fmt.Sprintf("amqp://%s:%s@%s:%d/", "admin", "admin", "Pass@word1", 5672)
	mqConn, err = amqp.Dial(RabbitUrl)
	r.Connection = mqConn // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("MQ打开链接失败:%s \n", err)
		return err
	}
	mqChan, err = mqConn.Channel()
	r.Channel = mqChan // 赋值给RabbitMQ对象
	if err != nil {
		fmt.Printf("MQ打开管道失败:%s \n", err)
		return err
	}
	return nil
}

// 关闭RabbitMQ连接
func (r *RabbitmqMessage) mqClose() {
	// 先关闭管道,再关闭链接
	err := r.Channel.Close()
	if err != nil {
		fmt.Printf("MQ管道关闭失败:%s \n", err)
	}
	err = r.Connection.Close()
	if err != nil {
		fmt.Printf("MQ链接关闭失败:%s \n", err)
	}
}

// 创建一个新的操作对象
func NewQueueExchange(q *QueueExchange) *RabbitmqMessage {
	return &RabbitmqMessage{
		QueueName:    q.QuName,
		RoutingKey:   q.RtKey,
		ExchangeName: q.ExName,
		ExchangeType: q.ExType,
	}
}

// 启动RabbitMQ客户端,并初始化
func (r *RabbitmqMessage) Start() {
	// 开启监听生产者发送任务
	for _, producer := range r.ProducerList {
		go r.listenProducer(producer)
	}
	// 开启监听接收者接收任务
	for _, receiver := range r.ReceiverList {
		go r.listenReceiver(receiver)
	}
	time.Sleep(1 * time.Second)
}

// 注册发送指定队列指定路由的生产者
func (r *RabbitmqMessage) RegisterProducer(producer Producer) {
	r.ProducerList = append(r.ProducerList, producer)
}

// 发送任务
func (r *RabbitmqMessage) listenProducer(producer Producer) {
	// 验证链接是否正常,否则重新链接
	if r.Channel == nil {
		r.InitRabbitmq()
	}
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err := r.Channel.QueueDeclarePassive(r.QueueName, true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.Channel.QueueDeclare(r.QueueName, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return
		}
	}
	// 队列绑定
	err = r.Channel.QueueBind(r.QueueName, r.RoutingKey, r.ExchangeName, true, nil)
	if err != nil {
		fmt.Printf("MQ绑定队列失败:%s \n", err)
		return
	}
	// 用于检查交换机是否存在,已经存在不需要重复声明
	err = r.Channel.ExchangeDeclarePassive(r.ExchangeName, r.ExchangeType, true, false, false, true, nil)
	if err != nil {
		// 注册交换机
		// name:交换机名称,kind:交换机类型,durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;
		// noWait:是否非阻塞, true为是,不等待RMQ返回信息;args:参数,传nil即可; internal:是否为内部
		err = r.Channel.ExchangeDeclare(r.ExchangeName, r.ExchangeType, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册交换机失败:%s \n", err)
			return
		}
	}
	// 发送任务消息
	err = r.Channel.Publish(r.ExchangeName, r.RoutingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(producer.MsgContent()),
	})
	if err != nil {
		fmt.Printf("MQ任务发送失败:%s \n", err)
		return
	}
}

// 注册接收指定队列指定路由的数据接收者
func (r *RabbitmqMessage) RegisterReceiver(receiver Receiver) {
	r.Mu.Lock()
	r.ReceiverList = append(r.ReceiverList, receiver)
	r.Mu.Unlock()
}

// 监听接收者接收任务
func (r *RabbitmqMessage) listenReceiver(receiver Receiver) {
	// 处理结束关闭链接
	defer r.mqClose()
	// 验证链接是否正常
	if r.Channel == nil {
		r.InitRabbitmq()
	}
	// 用于检查队列是否存在,已经存在不需要重复声明
	_, err := r.Channel.QueueDeclarePassive(r.QueueName, true, false, false, true, nil)
	if err != nil {
		// 队列不存在,声明队列
		// name:队列名称;durable:是否持久化,队列存盘,true服务重启后信息不会丢失,影响性能;autoDelete:是否自动删除;noWait:是否非阻塞,
		// true为是,不等待RMQ返回信息;args:参数,传nil即可;exclusive:是否设置排他
		_, err = r.Channel.QueueDeclare(r.QueueName, true, false, false, true, nil)
		if err != nil {
			fmt.Printf("MQ注册队列失败:%s \n", err)
			return
		}
	}
	// 绑定任务
	err = r.Channel.QueueBind(r.QueueName, r.RoutingKey, r.ExchangeName, true, nil)
	if err != nil {
		fmt.Printf("绑定队列失败:%s \n", err)
		return
	}
	// 获取消费通道,确保rabbitMQ一个一个发送消息
	err = r.Channel.Qos(1, 0, true)
	if err != nil {
		fmt.Printf("获取消费通道异常:%s \n", err)
	}
	msgList, err := r.Channel.Consume(r.QueueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Printf("获取消费通道异常:%s \n", err)
		return
	}
	for msg := range msgList {
		// 处理数据
		err := receiver.Consumer(msg.Body)
		if err != nil {
			err = msg.Ack(true)
			if err != nil {
				fmt.Printf("确认消息未完成异常:%s \n", err)
				return
			}
		} else {
			// 确认消息,必须为false
			err = msg.Ack(false)
			if err != nil {
				fmt.Printf("确认消息完成异常:%s \n", err)
				return
			}
			return
		}
	}
}

func (r *RabbitmqMessage) log(lv logmgr.LogLevel, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now().Format("2006/01/02 15:04:05")
	funcName, fileName, lineNo := logmgr.GetFileInfo(3)
	//  发送消息
	// Trap SIGINT to trigger a graceful shutdown.
	// signals := make(chan os.Signal, 1)
	// fmt.Println("Log is:", r.QueueName)
	_, err := r.Channel.QueueDeclare(
		r.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		fmt.Printf("MQ Queuename error:%s \n", err)
	}
	err = r.Channel.ExchangeDeclare(
		r.ExchangeName, // name
		r.ExchangeType, // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		fmt.Printf("MQ Exchange error:%s \n", err)
	}
	err = r.Channel.QueueBind(r.QueueName, r.RoutingKey, r.ExchangeName, false, nil)
	if err != nil {
		fmt.Println("MQ Channel Exchange Binding fail:", err)
	}
	if r.LogLv < logmgr.ERROR {
		message := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(r.LogLv), fileName, funcName, lineNo, msg)
		// fmt.Printf("Error message:%s\n", message)
		rmsg := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		}
		// fmt.Println("Error Log rmsg:", rmsg)
		err := mqChan.Publish(r.ExchangeName, r.RoutingKey, false, false, rmsg)
		if err != nil {
			fmt.Println("Error Log error:", err)
		}
		//debugLoop:

	} else if r.LogLv <= logmgr.FATAL {
		message := fmt.Sprintf("[%s][%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(r.LogLv), fileName, funcName, lineNo, msg)
		// fmt.Printf("Fatal message:%s\n", message)
		rmsg := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		}
		// fmt.Println("Fatal Log rmsg:", rmsg)
		err := mqChan.Publish(r.ExchangeName, r.RoutingKey, false, false, rmsg)
		if err != nil {
			fmt.Println("Error Log error:", err)
		}
	} else {
		fmt.Println("Loglevel set error:")
	}

}

func (r *RabbitmqMessage) Debug(format string, a ...any) {
	r.log(logmgr.DEBUG, format, a...)

}

func (r *RabbitmqMessage) Trace(format string, a ...any) {
	r.log(logmgr.TRACE, format, a...)
}

func (r *RabbitmqMessage) Info(format string, a ...any) {
	r.log(logmgr.INFO, format, a...)
}

func (r *RabbitmqMessage) Warning(format string, a ...any) {
	r.log(logmgr.WARNING, format, a...)
}

func (r *RabbitmqMessage) Error(format string, a ...any) {
	r.log(logmgr.ERROR, format, a...)
}

func (r *RabbitmqMessage) Fatal(format string, a ...any) {
	r.log(logmgr.FATAL, format, a...)
}
