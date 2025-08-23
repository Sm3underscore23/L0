package consumer

import (
	"log"

	"github.com/IBM/sarama"
)

type pool interface {
	Create()
	Handle(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage)
	Wait()
}

type workerPool struct {
	workers chan struct{}
	handle  Handler
}

func NewWP(workersNum int, handle Handler) pool {
	return &workerPool{
		handle:  handle,
		workers: make(chan struct{}, workersNum),
	}
}

func (wp *workerPool) Create() {
	for range cap(wp.workers) {
		wp.workers <- struct{}{}
	}
}

func (wp *workerPool) Handle(session sarama.ConsumerGroupSession, msg *sarama.ConsumerMessage) {
	<-wp.workers
	go func() {
		err := wp.handle.HandleOrder(session.Context(), msg.Value)
		if err != nil {
			log.Printf("failed to handle message: %v", err)
		}
		session.MarkMessage(msg, "")
		wp.workers <- struct{}{}
	}()
}

func (wp *workerPool) Wait() {
	for range len(wp.workers) {
		<-wp.workers
	}
}
