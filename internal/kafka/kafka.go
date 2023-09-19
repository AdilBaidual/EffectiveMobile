package kafka

import (
	"EffectiveMobile/internal/service"
)

type Consumer interface {
	Consume()
	Close() error
}

type KafkaMessageProcessor struct {
	Consumer
}

func NewKafkaMessageProcessor(brokerAddr, topic string, service *service.Service) (*KafkaMessageProcessor, error) {
	concumer, err := NewKafkaConsumer(brokerAddr, topic, service)
	if err != nil {
		return nil, err
	}
	return &KafkaMessageProcessor{Consumer: concumer}, nil
}

func (kmp *KafkaMessageProcessor) Start() {
	go kmp.Consume()
}
