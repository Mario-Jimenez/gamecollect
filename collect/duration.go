package collect

import (
	"encoding/json"

	"github.com/Mario-Jimenez/gamecollect/storage"
	log "github.com/sirupsen/logrus"
)

type DurationHandler struct {
	storage *storage.InMemory
}

func NewDurationHandler(storage *storage.InMemory) *DurationHandler {
	return &DurationHandler{storage}
}

func (h *DurationHandler) ProcessMessage(message []byte) {
	duration := &Duration{}
	err := json.Unmarshal(message, duration)
	if err != nil {
		log.WithFields(log.Fields{
			"message": string(message),
			"error":   err.Error(),
		}).Warning("An invalid message was received from the broker")
		return
	}

	log.WithFields(log.Fields{
		"duration": duration,
	}).Info("Inbound message from broker")

	h.storage.AddDuration(duration.ID, duration.Duration)
}
