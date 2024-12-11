package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

func CreateConsumer() *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "example-group",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to create consumer")
	}
	return consumer
}

func consumeStructuredMessages(consumer *kafka.Consumer, topic string) {
	err := consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to subscribe to topic")
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var message KafkaMessage
			err := json.Unmarshal(msg.Value, &message)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Debugln("Failed to unmarshal message")
			} else {
				logrus.WithFields(logrus.Fields{
					"message": string(msg.Value),
				}).Infoln("Received message")
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Debugln("Failed to read message")
		}
	}
}
