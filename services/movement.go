package services

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"slitherio/types"
	"time"
)

func LastPressedKey(snake *types.Snake) {
	for {
		ev := termbox.PollEvent()
		if ev.Type != termbox.EventKey {
			continue
		}

		switch ev.Ch {
		case 'w':
			if snake.LastDirection == 's' {
				continue
			}
			snake.LastDirection = 'w'
		case 'a':
			if snake.LastDirection == 'd' {
				continue
			}
			snake.LastDirection = 'a'
		case 's':
			if snake.LastDirection == 'w' {
				continue
			}
			snake.LastDirection = 's'
		case 'd':
			if snake.LastDirection == 'a' {
				continue
			}
			snake.LastDirection = 'd'
		}
	}
}

func Movement(snake *types.Snake, food *types.Positions, s *Service, room string) {
	ch := make(chan types.GameData)
	go s.SubscribeGameState(snake.Name, room, ch)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		for i := len(snake.Position) - 1; i > 0; i-- {
			snake.Position[i] = snake.Position[i-1]
		}

		switch snake.LastDirection {
		case 'w':
			snake.Position[0].Y--
		case 'a':
			snake.Position[0].X--
		case 's':
			snake.Position[0].Y++
		case 'd':
			snake.Position[0].X++
		}

		if snake.Position[0].X <= 0 || snake.Position[0].X >= 79 || snake.Position[0].Y <= 0 || snake.Position[0].Y >= 59 {
			_ = s.NatsConn.Publish("room:status:"+room, []byte("lose"))
			termbox.Close()
			fmt.Println(snake.Name + " hit the wall!")
			os.Exit(0)
		}

		if snake.Position[0].X == food.X && snake.Position[0].Y == food.Y {
			snake.Position = append(snake.Position, types.Positions{
				X: snake.Position[len(snake.Position)-1].X,
				Y: snake.Position[len(snake.Position)-1].Y,
			})
			food.X = -1
			food.Y = -1
		}

		s.PublishGameState(snake.Position, food, room, snake.Name)
		var opponent types.GameData
		select {
		case gameData := <-ch:
			opponent = gameData
		default:
		}

		if len(opponent.Opponent) != 0 && snake.Position[0] == opponent.Opponent[0] { // draw
			termbox.Close()
			fmt.Println("Draw!")
			_ = s.NatsConn.Publish("room:status:"+room, []byte("draw"))
			os.Exit(0)
		} else if len(opponent.Opponent) != 0 { // win - lose
			if state, status := IsAKill(snake.Position, opponent.Opponent); state {
				_ = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				_ = termbox.Flush()
				termbox.Close()
				if status {
					fmt.Println(snake.Name + " won!")
					_ = s.NatsConn.Publish("room:status:"+room, []byte("lose"))
				} else {
					fmt.Println(opponent.ID + " won!")
					_ = s.NatsConn.Publish("room:status:"+room, []byte("win"))
				}
				os.Exit(0)
			}
		}

		CreateTheTable(60, 80, snake.Position, food, &opponent)
	}
}

func IsAKill(snake []types.Positions, opponent []types.Positions) (bool, bool) {
	for _, j := range snake {
		if opponent[0] == j {
			return true, false
		}
	}

	for _, j := range opponent {
		if snake[0] == j {
			return true, true
		}
	}

	return false, false
}
