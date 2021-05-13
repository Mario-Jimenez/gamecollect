package storage

type Game struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	BasePrice   float32 `json:"base_price"`
	Score       float32 `json:"score"`
	Duration    int     `json:"duration"`
	Price       float32 `json:"price"`
	ScoreURL    string  `json:"score_url"`
	DurationURL string  `json:"duration_url"`
	PriceURL    string  `json:"price_url"`
}
