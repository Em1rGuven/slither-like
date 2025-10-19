package types

type (
	Info struct {
		PlayerName string `json:"playerName"`
		RoomName   string `json:"roomName"`
	}

	Snake struct {
		Name          string      `json:"name"`
		LastDirection byte        `json:"lastPressedKey"`
		Position      []Positions `json:"position"`
	}

	Positions struct {
		X int `json:"x"`
		Y int `json:"y"`
	}

	GameData struct {
		ID       string      `json:"id"`
		Opponent []Positions `json:"opponent"`
		Food     Positions   `json:"food"`
	}
)

func CreateSnake(playerName string, rank int) *Snake {
	var (
		xCord int
		key   byte
	)
	switch rank {
	case 1:
		xCord = 10
		key = 'd'
	case 2:
		xCord = 50
		key = 'a'
	}

	return &Snake{
		Name:          playerName,
		LastDirection: key,
		Position: []Positions{
			{xCord, 20},
			{xCord + 1, 20},
			{xCord + 2, 20},
			{xCord + 3, 20},
			{xCord + 4, 20},
		},
	}
}

func CreateFood() *Positions {
	return &Positions{
		X: -1,
		Y: -1,
	}
}
