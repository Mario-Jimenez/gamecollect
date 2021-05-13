package collect

import (
	"encoding/json"

	"github.com/Mario-Jimenez/gamecollect/storage"
	log "github.com/sirupsen/logrus"
)

type ScoreHandler struct {
	storage *storage.InMemory
}

func NewScoreHandler(storage *storage.InMemory) *ScoreHandler {
	return &ScoreHandler{storage}
}

func (h *ScoreHandler) ProcessMessage(message []byte) {
	score := &Score{}
	err := json.Unmarshal(message, score)
	if err != nil {
		log.WithFields(log.Fields{
			"message": string(message),
			"error":   err.Error(),
		}).Warning("An invalid message was received from the broker")
		return
	}

	log.WithFields(log.Fields{
		"score": score,
	}).Info("Inbound message from broker")

	h.storage.AddScore(score.ID, score.Score)
}
