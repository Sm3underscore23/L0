package consumer

import (
	"github.com/IBM/sarama"
)

type consumerGroupHandler struct {
	p       pool
	handler Handler
}

func (h consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h consumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	h.p.Create()

	for msg := range claim.Messages() {
		h.p.Handle(session, msg)
	}

	h.p.Wait()

	return nil
}
