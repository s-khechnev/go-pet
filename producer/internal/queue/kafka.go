package queue

import (
	"errors"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log/slog"
	"producer/internal/entity"
	"time"
)

type KafkaProducer struct {
	producer        *kafka.Producer
	topic           string
	flushTimeout    int
	deliveryTimeout time.Duration

	deliveryChan chan kafka.Event
}

func NewKafkaProducer(topic string, config *kafka.ConfigMap, flushTimeout int, deliveryTimeout time.Duration) (*KafkaProducer, error) {
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer:     producer,
		topic:        topic,
		flushTimeout: flushTimeout,

		deliveryTimeout: deliveryTimeout,
		deliveryChan:    make(chan kafka.Event),
	}, nil
}

var DeliveryTimeoutError = errors.New("delivery timeout")

func (p *KafkaProducer) Produce(message entity.Message) error {
	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            []byte(message.IP),
		Value:          []byte(message.Payload),
		Timestamp:      time.Now(),
	}, p.deliveryChan)
	if err != nil {
		slog.Error("failed produce message", slog.String("error", err.Error()))
		return err
	}

	select {
	case e := <-p.deliveryChan:
		switch ev := e.(type) {
		case *kafka.Message:
			m := ev
			if m.TopicPartition.Error != nil {
				slog.Error("failed to deliver message", slog.String("error", m.TopicPartition.Error.Error()))
				return m.TopicPartition.Error
			} else {
				slog.Info("delivered message",
					slog.String("topic", *m.TopicPartition.Topic),
					slog.Int("partition", int(m.TopicPartition.Partition)),
					slog.Int("offset", int(m.TopicPartition.Offset)))
				return nil
			}
		case kafka.Error:
			slog.Error("kafka error", slog.String("error", ev.Error()))
			return ev
		}
		return nil
	case <-time.After(p.deliveryTimeout):
		slog.Error("delivery timeout", slog.String("message_key", message.IP))
		return DeliveryTimeoutError
	}
}

func (p *KafkaProducer) Close() {
	num := p.producer.Flush(p.flushTimeout)
	slog.Info("number of outstanding events", slog.Int("num", num))
	p.producer.Close()
}
