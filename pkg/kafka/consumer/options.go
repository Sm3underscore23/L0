package consumer

import (
	"time"

	"github.com/IBM/sarama"
)

type Option func(*Consumer)

// ConnAttempts -.
func ConnAttempts(attempts int) Option {
	return func(c *Consumer) {
		c.connAttempts = attempts
	}
}

// ConnTimeout -.
func ConnTimeout(timeout time.Duration) Option {
	return func(c *Consumer) {
		c.connTimeout = timeout
	}
}

// Strategy -.
func Strategy(assignor string) Option {
	return func(c *Consumer) {
		switch assignor {
		case "range":
			c.cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
		case "sticky":
			c.cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
		case "round-robin":
			c.cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
		}
	}
}

// ShutdownTimeout -.
func ShutdownTimeout(timeout time.Duration) Option {
	return func(c *Consumer) {
		c.shutdownTimeout = timeout
	}
}
