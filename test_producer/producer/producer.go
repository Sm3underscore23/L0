package producer

import (
	"time"

	"github.com/IBM/sarama"
)

type Producer interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
}

func New(brokerList []string) (Producer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = time.Millisecond * 250
	config.Producer.Return.Successes = true
	_ = config.Producer.Partitioner
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
