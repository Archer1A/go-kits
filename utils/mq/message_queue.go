package mq

import "context"

type MessageQ interface {
	SendMessage(ctx context.Context, topic string, message []byte) string
	Subscribe(ctx context.Context, topic string) ([]byte, error)
}
