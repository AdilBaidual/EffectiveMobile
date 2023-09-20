package kafka

import (
	"EffectiveMobile/entity"
	"EffectiveMobile/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type FioResponse struct {
	entity.Fio
	Message string `json:"message"`
}

type KafkaConsumer struct {
	consumer *kafka.Consumer
	service  service.User
}

func NewKafkaConsumer(brokerAddr, topic string, service service.User) (*KafkaConsumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokerAddr,
		"group.id":          "my-consumer-group",
		"auto.offset.reset": "earliest",
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer, service: service}, nil
}

func (kc *KafkaConsumer) Consume() {
	for {
		msg, err := kc.consumer.ReadMessage(time.Second)
		if err == nil {
			var data entity.Fio
			var flag bool

			err := json.Unmarshal(msg.Value, &data)
			if err != nil {
				flag = true
				logrus.Error("Unmarshal error: ", err)
			}

			if (data.Name == "" || data.Surname == "") && !flag {
				flag = true
				err = errors.New("Incorrect message")
				logrus.Error("Incorrect message")
			}

			if flag {
				if errKafka := produceToFailedTopic(data, err.Error()); errKafka != nil {
					logrus.Error("Kafka producer error: ", errKafka)
				}
				continue
			}
			fmt.Println("Успешно ", data)
			kc.service.Test()
		}
	}
	kc.Close()
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}

func produceToFailedTopic(data entity.Fio, errorMessage string) error {
	config := &kafka.ConfigMap{
		"bootstrap.servers": viper.GetString("brokerAddr"),
	}

	producer, err := kafka.NewProducer(config)
	if err != nil {
		return err
	}
	defer producer.Close()

	topic := "FIO_FAILED"
	response := FioResponse{Fio: data, Message: errorMessage}
	value, _ := json.Marshal(response)

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	return producer.Produce(message, nil)
}
