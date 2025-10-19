package services

import (
	"github.com/nsf/termbox-go"
	"log"
	"math/rand"
	"slitherio/types"
	"sync"
)

var once = sync.Once{}

func CreateTheTable(rows, cols int, positions []types.Positions, food *types.Positions, opponent *types.GameData) {
	once.Do(func() {
		rawTerminal()
	})

	_ = termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	_ = termbox.Flush()

	for x := 0; x < cols; x++ {
		termbox.SetCell(x, 0, '#', termbox.ColorGreen, termbox.ColorDefault)
		termbox.SetCell(x, rows-1, '#', termbox.ColorGreen, termbox.ColorDefault)
	}
	for y := 0; y < rows; y++ {
		termbox.SetCell(0, y, '#', termbox.ColorGreen, termbox.ColorDefault)
		termbox.SetCell(cols-1, y, '#', termbox.ColorGreen, termbox.ColorDefault)
	}
	for _, pos := range positions {
		termbox.SetCell(pos.X, pos.Y, 'X', termbox.ColorYellow, termbox.ColorDefault)
	}

	if food.X == -1 || food.Y == -1 {
		for {
			food.X = rand.Intn(cols-2) + 1
			food.Y = rand.Intn(rows-2) + 1
			if !isOnSnake(food, positions) {
				break
			}
		}
	}
	termbox.SetCell(food.X, food.Y, 'O', termbox.ColorRed, termbox.ColorDefault)

	for _, pos := range opponent.Opponent {
		termbox.SetCell(pos.X, pos.Y, 'X', termbox.ColorBlue, termbox.ColorDefault)
	}
	termbox.SetCell(opponent.Food.X, opponent.Food.Y, 'O', termbox.ColorRed, termbox.ColorDefault)

	_ = termbox.Flush()
}

func isOnSnake(food *types.Positions, positions []types.Positions) bool {
	for _, pos := range positions {
		if pos.X == food.X && pos.Y == food.Y {
			return true
		}
	}
	return false
}

func rawTerminal() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
}
