package kafka

type Consumer interface {
	Consume() error
	Close() error
}

type KafkaMessageProcessor struct {
	consumer Consumer
}

func NewKafkaMessageProcessor(consumer Consumer) *KafkaMessageProcessor {
	return &KafkaMessageProcessor{consumer: consumer}
}

func (kmp *KafkaMessageProcessor) Start() {
	go kmp.consumer.Consume()
}
