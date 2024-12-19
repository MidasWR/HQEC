package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
	"time"
)

type KafkaMessage struct {
	Id      int       `json:"id"`
	UserId  int       `json:"user_id"`
	Balance float32   `json:"balance"`
	Date    time.Time `json:"date"`
}

func CreateProducer() (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func ProduceMessage(producer *kafka.Producer, topic string, message KafkaMessage) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to marshal message")
	}
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonMessage,
	}, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to produce message")
	}
	e := <-producer.Events()
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		logrus.WithFields(logrus.Fields{
			"error": msg.TopicPartition.Error,
		}).Debugln("Failed to produce message")
	} else {
		fmt.Printf("Message delivered to topic %s [%d] at offset %v\n",
			*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
	}
}
