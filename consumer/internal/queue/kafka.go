package queue

import (
	"consumer/internal/entity"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"time"
)

type AnalyticsRepository interface {
	SaveCountBook(book entity.Book) error
}

type Consumer struct {
	analyticsRepo AnalyticsRepository
	consumer      *kafka.Consumer
	topic         string
	context       context.Context
	timeoutOnPoll time.Duration
}

func NewConsumer(
	ctx context.Context,
	repo AnalyticsRepository,
	topic string,
	groupID string,
	config *kafka.ConfigMap,
	timeoutOnPoll time.Duration,
) (*Consumer, error) {
	err := config.SetKey("group.id", groupID)
	if err != nil {
		return nil, err
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	err = c.Subscribe(topic, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:      c,
		topic:         topic,
		analyticsRepo: repo,
		context:       ctx,
		timeoutOnPoll: timeoutOnPoll,
	}, nil
}

func (c *Consumer) Run() {
	slog.Info("consumer started")

	for {
		select {
		case <-c.context.Done():
			return
		default:
		}

		ev := c.consumer.Poll(int(c.timeoutOnPoll))
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			//fmt.Printf("%% Message on %s:\n%s\n",
			//	e.TopicPartition, string(e.Value))

			var book entity.Book
			err := json.Unmarshal(e.Value, &book)
			if err != nil {
				slog.Error("failed unmarshalling book from json", slog.String("error", err.Error()))
			}

			err = c.analyticsRepo.SaveCountBook(book)
			if err != nil {
				slog.Error("failed saving count book", slog.String("error", err.Error()))
			}

			_, err = c.consumer.CommitMessage(e)
			if err != nil {
				slog.Error("failed commit message", slog.String("error", err.Error()))
			}
		case kafka.Error:
			slog.Error("consumer error", slog.String("error", e.Error()))
		default:
		}
	}

}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
