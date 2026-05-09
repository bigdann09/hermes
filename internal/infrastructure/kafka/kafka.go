package kafka

type Kafka struct {
	brokers []string
}

func NewKafka(brokers []string) *Kafka {
	return &Kafka{brokers: brokers}
}
