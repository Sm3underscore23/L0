package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

const (
	_defaultConnAttempts    = 10
	_defaultConnTimeout     = time.Second
	_defaultShutdownTimeout = 3 * time.Second
)

type Handler interface {
	HandleOrder(ctx context.Context, value []byte) error
}

// Kafka Consumer
type Consumer struct {
	cfg       *sarama.Config
	group     sarama.ConsumerGroup
	topicName string
	handler   Handler
	error     chan error

	connAttempts int
	connTimeout  time.Duration

	shutdownTimeout time.Duration
}

// New -.
func New(brokers []string, group string, topicNmae string,
	handler Handler, opts ...Option,
) (*Consumer, error) {
	consumer := Consumer{
		cfg:       sarama.NewConfig(),
		topicName: topicNmae,
		handler:   handler,
		error:     make(chan error),

		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,

		shutdownTimeout: _defaultShutdownTimeout,
	}

	consumer.cfg.Version = sarama.V4_0_0_0
	consumer.cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer.cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	for _, opt := range opts {
		opt(&consumer)
	}

	var err error
	for consumer.connAttempts > 0 {
		consumer.group, err = sarama.NewConsumerGroup(brokers, group, consumer.cfg)
		if err == nil {
			return &consumer, nil
		}

		log.Printf("Trying connect to Kafka, attempts left: %d", consumer.connAttempts)
		consumer.connAttempts--
		time.Sleep(consumer.connTimeout)
	}

	return nil, fmt.Errorf("connAttempts == 0: %w", err)
}

// Start -.
func (c *Consumer) Start(ctx context.Context, workersNum int) error {
	h := consumerGroupHandler{handler: c.handler, p: NewWP(workersNum, c.handler)}
	if err := c.group.Consume(ctx, []string{c.topicName}, h); err != nil {
		return err
	}
	return nil
}

// Notify -.
func (c *Consumer) Notify() <-chan error {
	return c.error
}

// Shutdown -.
func (c *Consumer) Shutdown() error {
	select {
	case <-c.error:
		return nil
	default:
	}

	time.Sleep(c.shutdownTimeout)

	err := c.group.Close()
	if err != nil {
		return err
	}

	return nil
}
