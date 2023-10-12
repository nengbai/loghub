package kafka

import (
	"fmt"
	"loghub/src/logmgr"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

type KafkaMessage struct {
	LogLv       logmgr.LogLevel
	Topic       string
	KafkaClient []string
	DelayTopic  string
	Config      *sarama.Config
}

type KafkaConfig struct {
}

func NewKafkaMessage(logstr, topic, delayTopic string) *KafkaMessage {
	loglv, err := logmgr.ParseLoglevel(logstr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fl := &KafkaMessage{
		LogLv:      loglv,
		Topic:      topic,
		DelayTopic: delayTopic,
	}
	err = fl.InitKafkaLog()
	if err != nil {
		panic(err)
	}
	return fl
}

func (k *KafkaMessage) InitKafkaLog() error {
	// 1.读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	brokers := viper.GetStringSlice("kafka.addrs")
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	// sasl认证
	//  config.Net.SASL.Enable = true
	//  config.Net.SASL.User = "admin"
	//  config.Net.SASL.Password = "admin"

	// 2.获取Kafka Brocker IP and Port
	//brokers := []string{"192.168.101.9:19092", "192.168.101.9:29092", "192.168.101.9:39092"}
	k.KafkaClient = brokers
	k.Config = config

	return nil
}

func (k *KafkaMessage) log(lv logmgr.LogLevel, format string, a ...any) {

	msg := fmt.Sprintf(format, a...)
	now := time.Now().Format("2006/01/02 15:04:05")
	funcName, fileName, lineNo := logmgr.GetFileInfo(3)
	producer, err := sarama.NewAsyncProducer(k.KafkaClient, k.Config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// 4. 发送消息
	var (
		enqueued, producerErrors int //
	)

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)

	// timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	//fmt.Println("lv=:,k.LogLv:=", lv, k.LogLv)
	if k.LogLv < logmgr.ERROR {
		//fmt.Printf("newfile:%v\n", f.FileObj)
		message := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(k.LogLv), fileName, funcName, lineNo, msg)
		kafkaMsg := &sarama.ProducerMessage{
			Topic:     k.Topic,
			Partition: int32(-1),
			Value:     sarama.StringEncoder(message),
			Timestamp: time.Now().UTC(),
		}
		//debugLoop:
		select {
		case producer.Input() <- kafkaMsg:
			//	enqueued++
			//producer.Input() <- kafkaMsg
			enqueued++
		case err := <-producer.Errors():
			fmt.Println("Failed to produce message", err)
			producerErrors++
		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			signal.Notify(signals, os.Interrupt)
			break
		}
	} else if k.LogLv <= logmgr.FATAL {
		message := fmt.Sprintf("[%s][%s] [%s:%s:%d] %s\n", now, logmgr.GetLogString(k.LogLv), fileName, funcName, lineNo, msg)
		kafkaMsg := &sarama.ProducerMessage{
			Topic:     k.Topic,
			Partition: int32(-1),
			Key:       sarama.StringEncoder("key"),
			Value:     sarama.StringEncoder(message),
			Timestamp: time.Now().UTC(),
		}
		select {
		case producer.Input() <- kafkaMsg:
			enqueued++
		case err := <-producer.Errors():
			fmt.Println("Failed to produce message", err)
			producerErrors++
		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			signal.Notify(signals, os.Interrupt)
			break
		}
	} else {
		fmt.Println("Loglevel set error")
	}
	fmt.Printf("Successfully produced: %d; errors: %d\n", enqueued, producerErrors)

}

func (k *KafkaMessage) Debug(format string, a ...any) {
	k.log(logmgr.DEBUG, format, a...)

}

func (k *KafkaMessage) Trace(format string, a ...any) {
	k.log(logmgr.TRACE, format, a...)
}

func (k *KafkaMessage) Info(format string, a ...any) {
	k.log(logmgr.INFO, format, a...)
}

func (k *KafkaMessage) Warning(format string, a ...any) {
	k.log(logmgr.WARNING, format, a...)
}

func (k *KafkaMessage) Error(format string, a ...any) {
	k.log(logmgr.ERROR, format, a...)
}

func (k *KafkaMessage) Fatal(format string, a ...any) {
	k.log(logmgr.FATAL, format, a...)
}

func (k *KafkaMessage) Close() {

}
