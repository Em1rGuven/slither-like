package services

import (
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nsf/termbox-go"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func (s *Service) gracefulShutdown(room string) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	log.Println("graceful shutdown...")
	_ = s.Rdb.Decr(s.Ctx, "room:"+room)
	_ = s.Rdb.Close()
	s.NatsConn.Close()
	termbox.Close()
	os.Exit(0)
}

func (s *Service) GetRedisData(room string, rank *int, ch chan struct{}) {
	once := sync.Once{}
	ticker := time.NewTicker(time.Millisecond * 500)
	timeout := time.NewTimer(time.Second * 10)
	defer ticker.Stop()
	defer timeout.Stop()
	for {
		select {
		case <-ticker.C:
			valStr, err := s.Rdb.Get(s.Ctx, "room:"+room).Result()
			if err != nil {
				if errors.Is(err, redis.Nil) {
					continue
				}
				log.Printf("failed to get room value: %v", err)
				continue
			}
			val, _ := strconv.Atoi(valStr)
			if val == 2 {
				if *rank != 1 {
					*rank = 2
				}
				ch <- struct{}{}
				return
			} else if val == 1 {
				*rank = 1
				once.Do(func() {
					fmt.Println("waiting for someone else...")
				})
			} else if val > 2 {
				log.Printf("room %s is busy", room)
				return
			}
		case <-timeout.C:
			log.Println("seems like no one is coming...")
			return
		}
	}
}

func (s *Service) JoinAndWaitForFill(room string, rank *int) {
	ch := make(chan struct{})
	if err := s.Rdb.Incr(s.Ctx, "room:"+room).Err(); err != nil {
		log.Printf("failed to increment room count: %v", err)
		return
	}
	go s.GetRedisData(room, rank, ch)
	go s.gracefulShutdown(room)
	<-ch
	close(ch)
}

func (s *Service) GameState(room string) {
	_, _ = s.NatsConn.Subscribe("room:status:"+room, func(m *nats.Msg) {
		text := string(m.Data)

		switch text {
		case "win":
			termbox.Close()
			fmt.Println("You lost!")
			os.Exit(0)
		case "lose":
			termbox.Close()
			fmt.Println("You win!")
			os.Exit(0)
		case "draw":
			termbox.Close()
			fmt.Println("Draw!")
			os.Exit(0)
		}
	})
}
