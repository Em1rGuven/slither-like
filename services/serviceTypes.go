package services

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Rdb *redis.Client
	NatsConn *nats.Conn
	Ctx      context.Context
}
