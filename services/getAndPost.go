package services

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"slitherio/types"
)

func (s *Service) PublishGameState(position []types.Positions, food *types.Positions, room string, name string) {
	data := types.GameData{
		ID:       name,
		Opponent: position,
		Food:     *food,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to marshal game data: %v", err)
		return
	}

	_ = s.NatsConn.Publish("room:"+room, payload)
}

func (s *Service) SubscribeGameState(name, room string, ch chan types.GameData) {
	_, _ = s.NatsConn.Subscribe("room:"+room, func(m *nats.Msg) {
		var state types.GameData
		if err := json.Unmarshal(m.Data, &state); err != nil {
			return
		}

		if state.ID == name {
			return
		}
		ch <- state
	})
}
