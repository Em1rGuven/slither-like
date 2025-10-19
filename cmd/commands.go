package cmd

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"slitherio/types"
	"time"
)

var (
	info    *types.Info
	rootCmd = &cobra.Command{Use: "slither.io"}
)

func init() {
	var err error
	rand.Seed(time.Now().UnixNano())

	if err = godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	info = &types.Info{}
	id := uuid.New().String()
	id = id[:5]
	rootCmd.Flags().StringVarP(&info.PlayerName, "player", "p", id, "player name")
	rootCmd.Flags().StringVarP(&info.RoomName, "room", "r", "training room", "room name")
}

func Execute() *types.Info {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
	return info
}
