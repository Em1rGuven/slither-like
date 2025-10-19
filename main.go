package main

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"slitherio/cmd"
	"slitherio/services"
	"slitherio/types"
)

type Server struct {
	rdb      *redis.Client
	natsConn *nats.Conn
	ctx      context.Context
	service  *services.Service
}

var (
	rank int
)

func main() {
	server := newServer()
	info := cmd.Execute()
	server.service.JoinAndWaitForFill(info.RoomName, &rank)

	snake, food := types.CreateSnake(info.PlayerName, rank), types.CreateFood()
	go server.service.GameState(info.RoomName)
	services.CreateTheTable(60, 80, snake.Position, food, &types.GameData{})

	go services.LastPressedKey(snake)
	services.Movement(snake, food, server.service, info.RoomName)
	select {}
}
