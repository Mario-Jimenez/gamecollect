package collect

import (
	"encoding/json"

	"github.com/Mario-Jimenez/gamecollect/storage"
	log "github.com/sirupsen/logrus"
)

type GameHandler struct {
	storage *storage.InMemory
}

func NewGameHandler(storage *storage.InMemory) *GameHandler {
	return &GameHandler{storage}
}

func (h *GameHandler) ProcessMessage(message []byte) {
	game := &Game{}
	err := json.Unmarshal(message, game)
	if err != nil {
		log.WithFields(log.Fields{
			"message": string(message),
			"error":   err.Error(),
		}).Warning("An invalid message was received from the broker")
		return
	}

	log.WithFields(log.Fields{
		"game": game,
	}).Info("Inbound message from broker")

	h.storage.AddGame(game.ID, game.Name, game.ScoreURL, game.DurationURL, game.BasePrice)
}
