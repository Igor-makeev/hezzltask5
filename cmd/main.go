package main

import (
	"context"
	"hezzltask5/internal/config"
	"hezzltask5/internal/handler"
	"hezzltask5/internal/repository"
	rediscache "hezzltask5/internal/repository/cache/redis"
	logclickhouse "hezzltask5/internal/repository/log/clickhouse"
	"hezzltask5/internal/repository/storage/postgres"
	"hezzltask5/internal/server"
	"hezzltask5/internal/service"
	"hezzltask5/internal/service/queue"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		logrus.Fatal(err)
	}
	ctx := context.Background()
	natsconn, err := queue.NewNatsconn(config.NatsAddres)
	if err != nil {
		logrus.Fatalf("Unable to connect to queue: %v\n", err)
	}
	defer natsconn.Close()

	clickhouseconn, err := logclickhouse.NewClickHouseClient(config.LogDNS)
	if err != nil {
		logrus.Fatalf("Unable to connect to log: %v\n", err)
	}
	defer clickhouseconn.Close()

	postgresconn, err := postgres.NewPostgresClient(config.DBDSN)
	if err != nil {
		logrus.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer postgresconn.Close(ctx)

	rediscon, err := rediscache.NewRedisClient(config.RedisAddress)
	if err != nil {
		logrus.Fatalf("Could not connect to Redis: %v", err)
	}
	defer rediscon.Close()

	log := logclickhouse.NewLogClickHouse(clickhouseconn)
	queue := queue.NewNatsQueue(natsconn, log)
	queue.Run()
	cache := rediscache.NewRedisCache(rediscon)
	storage := postgres.NewPostgresStorage(postgresconn)

	repo := repository.NewRepository(storage, cache)
	service := service.NewService(repo, queue)
	handler := handler.NewHandler(service)
	server := server.NewServer(config.ServerAddress, handler)
	serverErrChan := server.Run()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-signals:

		logrus.Info("main: got shutdown signal. Shutting down...")
		server.Shutdown()

	case <-serverErrChan:
		logrus.Info("main: got server err signal. Shutting down...")
		server.Shutdown()

	}

}
