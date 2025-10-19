package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"slitherio/services"
)

func newServer() *Server {
	var (
		rdb      *redis.Client
		natsConn *nats.Conn
		err      error
	)

	rdb = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	if err = rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}

	natsConn, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	server := &Server{
		rdb:      rdb,
		natsConn: natsConn,
		ctx:      ctx,
		service: &services.Service{
			Rdb: rdb,
			NatsConn: natsConn,
			Ctx: ctx,
		},
	}

	return server
}
