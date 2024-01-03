package hivemq

import (
	"fmt"
	"loghub/src/logmgr"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

type Producer interface {
	Publish() string
}

// 定义接收者接口
type Receiver interface {
	Subscribe([]byte) error
}
type MqttMessage struct {
	LogLv  logmgr.LogLevel
	Client mqtt.Client
	Topic  string
	QoS    int8
}

func NewMQTTMessage(logstr string) *MqttMessage {
	loglv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fl := &MqttMessage{
		LogLv: loglv,
	}
	err = fl.InitHivemq()
	logmgr.FailOnError(err, "Redis初始化失败:")
	return fl
}

func (r *MqttMessage) InitHivemq() error {
	// 1.读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	brokers := viper.GetStringSlice("hivemq.addrs")
	// Redis用户的密码
	password := viper.GetString("hivemq.password")
	Username := viper.GetString("hivemq.username")
	Qos := viper.GetInt("hivemq.QoS2")

	Topic := viper.GetString("hivemq.topic")

	// Redis Broker 的ip地址
	/** NewClientOptions() default value
	Port: 1883
	CleanSession: True
	Order: True (note: it is recommended that this be set to FALSE unless order is important)
	KeepAlive: 30 (seconds)
	ConnectTimeout: 30 (seconds)
	MaxReconnectInterval 10 (minutes)
	AutoReconnect: True
		**/
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokers[0])
	// 设置用户名
	opts.SetUsername(Username)
	// 设置密码
	opts.SetPassword(password)
	// 使用连接信息进行连接
	MqttClient := mqtt.NewClient(opts)
	if token := MqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("订阅 MQTT 失败")
		panic(token.Error())
	}
	r.Client = MqttClient
	r.QoS = int8(Qos)
	r.Topic = Topic
	return nil
}

func Subscribe(client mqtt.Client, topic string, qos byte) {
	// 使用模式订阅
	client.Subscribe(topic, qos, subCallBackFunc)

	// 处理订阅接收到的消息
}

// subCallBackFunc
/**
 *  @Description: 回调函数
 *  @param client
 *  @param msg
 */
func subCallBackFunc(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("订阅: 当前话题是 [%s]; 信息是 [%s] \n", msg.Topic(), string(msg.Payload()))
}

func Publisher(client mqtt.Client, topic string, qos byte, message string) {
	for {
		// 发布消息到频道
		token := client.Publish(topic, qos, true, message)
		if token != nil {
			panic(token)
		}
		token.Wait()

		client.Disconnect(250)
		fmt.Println("Message published.")
		time.Sleep(1 * time.Second)
	}
}
func (r *MqttMessage) log(lv logmgr.LogLevel, format string, a ...any) {
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
		Publisher(r.Client, r.Topic, byte(r.QoS), message)

	} else if r.LogLv <= logmgr.FATAL {
		message := fmt.Sprintf("[%s][%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(r.LogLv), fileName, funcName, lineNo, msg)
		// fmt.Printf("Fatal message:%s\n", message)

		Publisher(r.Client, r.Topic, byte(r.QoS), message)

	} else {
		fmt.Println("Loglevel set error:")
	}

}

func (r *MqttMessage) Debug(format string, a ...any) {
	r.log(logmgr.DEBUG, format, a...)

}

func (r *MqttMessage) Trace(format string, a ...any) {
	r.log(logmgr.TRACE, format, a...)
}

func (r *MqttMessage) Info(format string, a ...any) {
	r.log(logmgr.INFO, format, a...)
}

func (r *MqttMessage) Warning(format string, a ...any) {
	r.log(logmgr.WARNING, format, a...)
}

func (r *MqttMessage) Error(format string, a ...any) {
	r.log(logmgr.ERROR, format, a...)
}

func (r *MqttMessage) Fatal(format string, a ...any) {
	r.log(logmgr.FATAL, format, a...)
}

func (r *MqttMessage) Consumer(msg []byte) error {
	return nil
}
