package server

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaServer struct {
	r *kafka.Reader
}

func NewKafkaServer(brokers []string, topic string, groupID string) *KafkaServer {
	server := &KafkaServer{}
	kafkaConfig := kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
	}
	if groupID != "" {
		kafkaConfig.GroupID = groupID
	}
	server.r = kafka.NewReader(kafkaConfig)

	return server
}

//ReadFromOffset - reads from offset without commit
func (k *KafkaServer) ReadFromOffset(ctx context.Context, offset int64) (*kafka.Message, error) {
	err := k.r.SetOffset(offset)
	if err != nil {
		return nil, err
	}
	msg, err := k.r.FetchMessage(ctx)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (k *KafkaServer) CommitMessage(ctx context.Context, msgs ...kafka.Message) (err error) {
	return k.r.CommitMessages(ctx, msgs...)
}

//ReadMessage - reads message with commit, if group provided
func (k *KafkaServer) ReadMessage(ctx context.Context) (*kafka.Message, error) {
	msg, err := k.r.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (k *KafkaServer) Close() (err error) {
	err = k.r.Close()
	if err != nil {
		return
	}
	return
}
