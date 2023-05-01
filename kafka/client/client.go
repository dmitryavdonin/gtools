package client

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	w *kafka.Writer
}

func NewKafkaClient(brokers []string, topic string, groupID string) *KafkaClient {
	client := &KafkaClient{}
	client.w = kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	return client
}

func (k *KafkaClient) SendMessage(ctx context.Context, key string, msg []byte) (err error) {
	return k.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: msg,
	})
}

func (k *KafkaClient) Close() (err error) {
	err = k.w.Close()
	if err != nil {
		return
	}
	return
}
