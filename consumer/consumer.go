package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

func main() {
	// 连接Kafka集群
	brokers := []string{"192.168.101.9:19092", "192.168.101.9:29092", "192.168.101.9:39092"}
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
