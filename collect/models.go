package collect

type (
	Game struct {
		ID          string  `json:"id"`
		Name        string  `json:"name"`
		BasePrice   float32 `json:"base_price"`
		ScoreURL    string  `json:"score_url"`
		DurationURL string  `json:"duration_url"`
	}

	Score struct {
		ID    string  `json:"id"`
		Score float32 `json:"score"`
	}

	Duration struct {
		ID       string `json:"id"`
		Duration int    `json:"timeBeat"`
	}

	Price struct {
		ID    string  `json:"id"`
		Price float32 `json:"price"`
		URL   string  `json:"url"`
	}
)
