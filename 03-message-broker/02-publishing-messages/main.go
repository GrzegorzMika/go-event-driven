package main

import (
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

// type LoggerAdapter interface {
// 	Error(msg string, err error, fields LogFields)
// 	Info(msg string, fields LogFields)
// 	Debug(msg string, fields LogFields)
// 	Trace(msg string, fields LogFields)
// 	With(fields LogFields) LoggerAdapter
// }

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type SlogWatermillLogger struct{}

func (s SlogWatermillLogger) Error(msg string, err error, fields watermill.LogFields) {
	opLogger := logger
	for k, v := range fields {
		opLogger = opLogger.With(k, v)
	}
	opLogger.With("error", err).Error(msg)
}

func (s SlogWatermillLogger) Info(msg string, fields watermill.LogFields) {
	opLogger := logger
	for k, v := range fields {
		opLogger = opLogger.With(k, v)
	}
	opLogger.Info(msg)
}

func (s SlogWatermillLogger) Debug(msg string, fields watermill.LogFields) {
	opLogger := logger
	for k, v := range fields {
		opLogger = opLogger.With(k, v)
	}
	opLogger.Debug(msg)
}

func (s SlogWatermillLogger) Trace(msg string, fields watermill.LogFields) {
	opLogger := logger
	for k, v := range fields {
		opLogger = opLogger.With(k, v)
	}
	opLogger.Info(msg)
}

func (s SlogWatermillLogger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	opLogger := logger
	for k, v := range fields {
		opLogger = opLogger.With(k, v)
	}
	return s
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, SlogWatermillLogger{})
	if err != nil {
		panic(err)
	}
	defer func() { _ = publisher.Close() }()

	err = publisher.Publish("progress", message.NewMessage(watermill.NewShortUUID(), []byte("50")))
	if err != nil {
		panic(err)
	}
	err = publisher.Publish("progress", message.NewMessage(watermill.NewShortUUID(), []byte("100")))
	if err != nil {
		panic(err)
	}
}
