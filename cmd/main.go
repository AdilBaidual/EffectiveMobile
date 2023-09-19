package main

import (
	"EffectiveMobile/internal/kafka"
	"fmt"
)

func main() {
	fmt.Println("server started")
	kafkaConsumer, err := kafka.NewKafkaConsumer("localhost:9092", "FIO")
	if err != nil {
		// Обработка ошибки
		return
	}

	messageProcessor := kafka.NewKafkaMessageProcessor(kafkaConsumer)
	messageProcessor.Start()
}
