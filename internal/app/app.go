// Package app configures and runs application.
package app

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	lru "github.com/hashicorp/golang-lru/v2"
	"taskL0/internal/config"
	"taskL0/internal/controller/http"
	loggertag "taskL0/internal/entity/logger_tag"
	"taskL0/internal/entity/order"
	orderRepo "taskL0/internal/repository/order"
	cacheloader "taskL0/internal/usecase/cache-loader"
	orderUsecase "taskL0/internal/usecase/order"
	"taskL0/pkg/httpserver"
	"taskL0/pkg/kafka/consumer"
	"taskL0/pkg/logger"
	"taskL0/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg config.Config) {
	ctx := context.Background()

	// Logger
	err := logger.InitLogger(cfg.Log.Level)
	if err != nil {
		logger.Fatal("failed to init logger", err)
	}

	// Validator
	v := validator.New(validator.WithRequiredStructEnabled())

	// Database
	pg, err := postgres.New(cfg.PostgresURl(), cfg.PG.MigrationsDir, postgres.DoMigrations())
	if err != nil {
		logger.Fatal("failed postgres", err)
	}
	defer pg.Close()

	// In memory LRU cache
	cache, err := lru.New[order.OrderUID, order.OrderInfo](cfg.Cache.Limit)
	if err != nil {
		logger.Fatal("failed to init lru cache", err)
	}

	// Repository Order
	orderRepo := orderRepo.New(pg)

	// Usecase CacheLoader
	cacheLoader := cacheloader.New(orderRepo, cache, cfg.Cache.RecoverLimit)

	go func() {
		if err := cacheLoader.WarmUp(ctx); err != nil {
			logger.Fatal("failed to warm up cache", err)
		}
	}()

	// Usecase Order
	orderUsecase := orderUsecase.New(orderRepo, cache, v)

	// Kafka
	consumer, err := consumer.New(cfg.Kafka.BrokerList, cfg.Kafka.ConsumerGroup, cfg.Kafka.Topic, orderUsecase)
	if err != nil {
		logger.Fatal("failed to init consumer", err)
	}

	go func() {
		if err := consumer.Start(context.Background(), cfg.Kafka.WorkersNum); err != nil {
			logger.Fatal("consumer failed", err)
		}
	}()

	// HTTP Server
	router := http.NewRouter(orderUsecase, v)
	httpServer := httpserver.New(
		httpserver.Port(cfg.ServerAddress()),
	)
	httpServer.App.Handler = router

	// Start servers
	httpServer.Start()

	// App stat
	slog.Info(
		"server started SUCCESSFULLY",
		loggertag.ServerAddr, cfg.ServerAddress(),
		slog.Group(loggertag.DBInfo,
			loggertag.DBAddr, net.JoinHostPort(cfg.PG.Host, cfg.PG.Port),
			loggertag.DBName, cfg.PG.Name,
			loggertag.DBUser, cfg.PG.User,
		),
		slog.Group(loggertag.KafkaInfo,
			loggertag.KafkaConsumerGroup, cfg.Kafka.ConsumerGroup,
			loggertag.KafkaBrockerList, cfg.Kafka.BrokerList,
			loggertag.KafkaTopic, cfg.Kafka.Topic,
		),
	)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		slog.Info("app - Run - signal", "signal", s.String())
		slog.Info("shutting down gracefully...")
	case err = <-httpServer.Notify():
		slog.Error("app - Run - httpServer.Notify", loggertag.Error, err)
	case err = <-consumer.Notify():
		slog.Error("app - Run - rmqServer.Notify", loggertag.Error, err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		slog.Error("app - Run - httpServer.Shutdown", loggertag.Error, err)
	}
	slog.Info("server stopped")

	err = consumer.Shutdown()
	if err != nil {
		slog.Error("app - Run - consumer.Shutdown", loggertag.Error, err)
	}
	slog.Info("kafka stopped")
}
