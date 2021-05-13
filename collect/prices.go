package collect

import (
	"encoding/json"

	"github.com/Mario-Jimenez/gamecollect/storage"
	log "github.com/sirupsen/logrus"
)

type PriceHandler struct {
	storage *storage.InMemory
}

func NewPricesHandler(storage *storage.InMemory) *PriceHandler {
	return &PriceHandler{storage}
}

func (h *PriceHandler) ProcessMessage(message []byte) {
	price := &Price{}
	err := json.Unmarshal(message, price)
	if err != nil {
		log.WithFields(log.Fields{
			"message": string(message),
			"error":   err.Error(),
		}).Warning("An invalid message was received from the broker")
		return
	}

	log.WithFields(log.Fields{
		"price": price,
	}).Info("Inbound message from broker")

	h.storage.AddPrice(price.ID, price.URL, price.Price)
}
