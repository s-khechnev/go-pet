package queue

import (
	"encoding/json"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"producer/internal/entity"
	"time"
)

type KafkaProducer struct {
	producer     *kafka.Producer
	topic        string
	flushTimeout int
}

func NewKafkaProducer(topic string, config *kafka.ConfigMap, flushTimeout int) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					slog.Error("failed to deliver book", slog.String("error", m.TopicPartition.Error.Error()))
				} else {
					slog.Info("delivered book",
						slog.String("topic", *m.TopicPartition.Topic),
						slog.Int("partition", int(m.TopicPartition.Partition)),
						slog.Int("offset", int(m.TopicPartition.Offset)),
						slog.String("book_id", string(m.Key)))
				}
			case kafka.Error:
				slog.Error("kafka error", slog.String("error", ev.Error()))
			}
		}
	}()

	return &KafkaProducer{
		producer:     producer,
		topic:        topic,
		flushTimeout: flushTimeout,
	}, nil
}

var ErrJsonMarshaling = errors.New("marshaling error")

func (p *KafkaProducer) Produce(book entity.Book) error {
	bookJson, err := json.Marshal(book)
	if err != nil {
		slog.Error("failed marshal book", slog.String("error", err.Error()))
		return ErrJsonMarshaling
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            []byte(book.Id),
		Value:          bookJson,
		Timestamp:      time.Now(),
	}, nil)
	if err != nil {
		slog.Error("failed produce book", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (p *KafkaProducer) Close() {
	num := p.producer.Flush(p.flushTimeout)
	slog.Info("number of outstanding events", slog.Int("num", num))
	p.producer.Close()
}
