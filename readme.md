# What's Loghub 
Loghub is an open-source structured logger for Go (golang) which completely API compatible with the standard library logger  and used by thousands of companies for high-performance log output, message intergration, and mission-critical applications.
Loghub is easy to set debug-mode and errors-mode for deveplopment and product. I believe Loghub' biggest contribution is to have played a part in today's widespread use of structured logging in Golang. Trere doesn't seem to be a reason to do a major, breaking iteration into Loghub V2, since the fantastic Go community has built those independently. Many fantastic alternatives have sprung up. Loghub would look like those, had it been re-designed with what we know about structured logging in Go today. Check out, for example,console log,file log, and message hub which is eseay to interited with kafka,RabbitMQ,Redis and HiveMQ. You only use it like log output. please give a thumbs up to github.com/nenbai/loghub. 

## 1. licenses
The Loghub licenses to distribute software and documentation, and to accept regular contributions from individuals and corporations and larger grants of existing software products.
These licenses help us achieve our goal of providing reliable and long-lived software products through collaborative, open-source software development. In all cases, contributors retain full rights to use their original contributions for any other purpose outside of Loghub and its projects the right to distribute and build upon their work within Loghub.

## 2. What's scenario for Loghub?
### 2.1  log output set debug-mode and errors-mode for deveplopment and product.
It's easy to set log mode and only set paramater when you create log struct instance.
for example:mylog := conselog.NewConsoleLog("ERROR")

### 2.2. How to interaged with Kafka?
First, you should have kafka cluster, if you don't have kafka cluster, you can use below scripts to deploy a kafka cluster asap. 
for example: deployment kafka(kraft mode) with docker-compose 
1. Install Docker and docker-compose enviroment.
```
sudo yum install docker -y
sudo curl -SL https://github.com/docker/compose/releases/download/v2.15.1/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
sudo chmod +x /usr/bin/docker-compose
docker-compose --version
```

2. configure for docker-compose 
```
version: "3.6"
services:
  kafka1:
    container_name: kafka1
    image: 'bitnami/kafka:3.3.1'
    user: root
    ports:
      - '19092:9092'
      - '19093:9093'
    environment:
      # enalable Kraft mode
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      # Define kafka brocker internat ip and port 
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      # Set security protocal
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      # Define extra access IP and port
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://1.1.1.1:19092
      - KAFKA_BROKER_ID=1
      - KAFKA_KRAFT_CLUSTER_ID=iZWRiSqjZAlYwlKEqHFQWI
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.23.0.11:9093,2@172.23.0.12:9093,3@172.23.0.13:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      # Set inital memory and maximum memory for kafka broker 
      - KAFKA_HEAP_OPTS=-Xmx512M -Xms256M
	# Set kafka broker data path
    volumes:
      - /opt/volume/kafka/broker01:/bitnami/kafka:rw
	# Set kafka broker network ip
    networks:
      netkafka:
        ipv4_address: 172.23.0.11
  kafka2:
    container_name: kafka2
    image: 'bitnami/kafka:3.3.1'
    user: root
    ports:
      - '29092:9092'
      - '29093:9093'
    environment:
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://1.1.1.1:29092  #修改宿主机ip
      - KAFKA_BROKER_ID=2
      - KAFKA_KRAFT_CLUSTER_ID=iZWRiSqjZAlYwlKEqHFQWI #哪一，三个节点保持一致
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.23.0.11:9093,2@172.23.0.12:9093,3@172.23.0.13:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_HEAP_OPTS=-Xmx512M -Xms256M
    volumes:
      - /opt/volume/kafka/broker02:/bitnami/kafka:rw
    networks:
      netkafka:
        ipv4_address: 172.23.0.12
  kafka3:
    container_name: kafka3
    image: 'bitnami/kafka:3.3.1'
    user: root
    ports:
      - '39092:9092'
      - '39093:9093'
    environment:
      - KAFKA_ENABLE_KRAFT=yes
      - KAFKA_CFG_PROCESS_ROLES=broker,controller
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://1.1.1.1:39092  #修改宿主机ip
      - KAFKA_BROKER_ID=3
      - KAFKA_KRAFT_CLUSTER_ID=iZWRiSqjZAlYwlKEqHFQWI
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=1@172.23.0.11:9093,2@172.23.0.12:9093,3@172.23.0.13:9093
      - ALLOW_PLAINTEXT_LISTENER=yes
      - KAFKA_HEAP_OPTS=-Xmx512M -Xms256M
    volumes:
      - /opt/volume/kafka/broker03:/bitnami/kafka:rw
    networks:
      netkafka:
        ipv4_address: 172.23.0.13
networks:
  name:
  netkafka:
    driver: bridge
    name: netkafka
    ipam:
      driver: default
      config:
        - subnet: 172.23.0.0/25
          gateway: 172.23.0.1
```
and than, you only add kafka broker IP and port in config.yaml.

## 3. How to use Loghub?

### 3.1. Log output example
```
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

```

### 3.2. Messages Exchange Hub example for Kafka

Edit config/config.yaml add kafka configure as below:
```
kafka:
  addrs: ["192.168.101.9:19092","192.168.101.9:29092","192.168.101.9:39092"]
```

kafka producer for example:

```
package main

import (
	"fmt"
	"loghub/src/kafka"
	"loghub/src/logmgr"
)

func main() {
	logstr := "Fatal"
	topic := "test-demo"
	delayTopic := "delay-demo"
	k := kafka.NewKafkaMessage(logstr, topic, delayTopic)
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

```

kafka consumer for example:

```
package main

import (
	"fmt"
	"os"
	"os/signal"
  "github.com/spf13/viper"
	"github.com/IBM/sarama"
)

func main() {
  // 1.读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
  brokers := viper.GetStringSlice("kafka.addrs")
  fmt.Println(" brokers:", brokers)
	// 连接Kafka集群
	config := sarama.NewConfig()
	// 控制每次从 Kafka 中获取的最大记录数 1000
	config.Consumer.Fetch.Max = 10
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	topic := "test-demo"
	// 订阅主题
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println("consumer.partitions error:", err.Error())
		panic(err)
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			panic(err)
		}
		defer pc.AsyncClose()

		// 处理消息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition: %d, Offset: %d, Key: %s, Value: %s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
			}
		}(pc)
	}
	signals := make(chan os.Signal, 1)
	select {
	case <-signals:
		signal.Notify(signals, os.Interrupt)
	}
}

```
### 3.3. Messages Exchange Hub example for RabbitMQ
Edit config/config.yaml add RabbitMQ configure as below:
```
rabbitmq:
  user: guest
  password: **********
  addrs: 
  - "192.168.0.84:5672"
  vhost: /
```

RabbitMQ Producer code for example：
```
package main

import (
	"fmt"
	"loghub/src/logmgr"
	"loghub/src/rabbitmq"

	"github.com/streadway/amqp"
)

func main() {
	var (
		logstr       string = "Fatal"
		queueName    string = "TopicDemo"
		routingKey   string = "Rabbit"
		exchangeName string = "TopicExchange" // RabbitExchange
		exchangeType string = "topic"         //"",fanout,direct,topic
		reliable     bool   = false           //
	)
	queueExch := rabbitmq.QueueExchange{
		QuName: queueName,
		RtKey:  routingKey,
		ExName: exchangeName,
		ExType: exchangeType,
	}
	rq := rabbitmq.NewRabbitmqMessage(logstr, &queueExch)
	lv, err := logmgr.ParseLoglevel(logstr)
	logmgr.FailOnError(err, "Failed to Parse Log Level")
	// Reliable publisher confirms require confirm.select support from the connection.
	if reliable {
		fmt.Println("enabling publishing confirms.")
		if err := rq.Channel.Confirm(false); err != nil {
			fmt.Printf("Channel could not be put into confirm mode: %s\n", err)
		}
		confirms := rq.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))
		defer rabbitmq.ConfirmOne(confirms)
	}
	for {
		if lv < logmgr.ERROR {
			rq.Debug("This is Debug log")
			rq.Error("This is Errors log")
			rq.Info("This is Info log")
			rq.Warning("This is Warning log")
		} else if lv <= logmgr.FATAL {
			rq.Error("This is Errors log")
			rq.Fatal("This is Fatal log")
		}
	}
}

```


RabbitMQ Consumer code for example：
```
package main

import (
	"log"
	"loghub/src/logmgr"
	"loghub/src/rabbitmq"
	"os"
	"os/signal"
)

func main() {
	var (
		logstr       string = "Fatal"
		queueName    string = "TopicDemo"
		routingKey   string = "Rabbit"
		exchangeName string = "TopicExchange" // RabbitExchange
		exchangeType string = "topic"         //"",fanout,direct,topic
	)
	queueExch := rabbitmq.QueueExchange{
		QuName: queueName,
		RtKey:  routingKey,
		ExName: exchangeName,
		ExType: exchangeType,
	}
	rq := rabbitmq.NewRabbitmqMessage(logstr, &queueExch)

	msgs, err := rq.Channel.Consume(
		rq.QueueName, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	logmgr.FailOnError(err, "Rabbitmq Consumer Failure")

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", string(d.Body))
		}
	}()
	signals := make(chan os.Signal, 1)
	select {
	case <-signals:
		signal.Notify(signals, os.Interrupt)
	}
}

```

### 3.4. Messages Exchange Hub example for Redis
Edit config/config.yaml add RabbitMQ configure as below:

```
redis:
  addrs: 
  - "190.2.28.1:26379"
  password: *ra@****.****
  db: 0
  poolSize: 1000
  minIdleConn: 100
```

Redis Producer puscribe code for example：
```
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

```
Redis Consumer subcribe code for example：
```
package main

import (
	"loghub/src/redis"
	"os"
	"os/signal"
)

func main() {
	var (
		logstr      string = "Fatal"
		ChannelName string = "mylog" //
	)

	rd := redis.NewRedisMessage(logstr)
	// 验证链接是否正常

	if rd.Client == nil {
		rd.InitRedis()
	}
	// 获取消费通道,确保Redis一个一个发送消息

	//logmgr.FailOnError(err, "Rabbitmq Consumer Failure")
	go func() {
		redis.Subscriber(rd.Client, ChannelName)
	}()
	signals := make(chan os.Signal, 1)
	select {
	case <-signals:
		signal.Notify(signals, os.Interrupt)
	}
}

```

### 3.5. Messages Exchange Hub example for mqtt
Edit config/config.yaml add hiveMQ configure as below:

```
hivemq:
  addrs: 
  - "190.2.28.1:1883"
  password: O**@******om
  username: admin
  topic: mylog
  QoS0: 0    // 至多一次
  QoS1: 1    // 至少一次
  QoS2: 2    // 确保只有一次
```

MQTT Producer puscribe code for example：
```
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
```

MQTT Consumer subcribe code for example：
```
package main

import (
	"loghub/src/hivemq"
	"os"
	"os/signal"
)

func main() {
	var (
		logstr string = "Fatal"
	)

	mt := hivemq.NewMQTTMessage(logstr)

	if mt.Client == nil {
		mt.InitHivemq()
	}

	// 获取消费通道,确保hivemq一个一个发送消息 qos2
	//logmgr.FailOnError(err, "Rabbitmq Consumer Failure")
	go func() {
		hivemq.Subscribe(mt.Client, mt.Topic, byte(mt.QoS))
	}()
	signals := make(chan os.Signal, 1)
	select {
	case <-signals:
		signal.Notify(signals, os.Interrupt)
	}
}
```