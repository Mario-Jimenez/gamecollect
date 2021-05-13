package storage

import (
	"fmt"
	"sync"
)

type InMemory struct {
	// concurrency safe
	mu    sync.Mutex
	games map[string]*Game
}

func NewInMemoryHandler() *InMemory {
	return &InMemory{
		games: map[string]*Game{},
	}
}

func (m *InMemory) AddGame(id, name, scoreURL, durationURL string, basePrice float32) {
	game := &Game{
		ID:          id,
		Name:        name,
		BasePrice:   basePrice,
		ScoreURL:    scoreURL,
		DurationURL: durationURL,
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.games[id] = game
}

func (m *InMemory) AddScore(id string, score float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.games[id]; ok {
		m.games[id].Score = score
	}
}

func (m *InMemory) AddDuration(id string, duration int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.games[id]; ok {
		m.games[id].Duration = duration
	}
}

func (m *InMemory) AddPrice(id, url string, price float32) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.games[id]; ok {
		m.games[id].Price = price
		m.games[id].PriceURL = url
	}
}

func (m *InMemory) Get() []Game {
	m.mu.Lock()
	defer m.mu.Unlock()
	gameList := []Game{}
	for i := 0; i <= 99; i++ {
		if g, ok := m.games[fmt.Sprintf("%d", i)]; ok {
			gameList = append(gameList, *g)
		}
	}

	return gameList
}
