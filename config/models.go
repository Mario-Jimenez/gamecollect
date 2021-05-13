package config

// Values stores the app configuration values
type Values struct {
	LogLevel        string   `json:"logLevel"`
	Port            string   `json:"port"`
	KafkaConnection []string `json:"kafkaConnection"`
	NumberOfThreads int      `json:"numberOfThreads"`
}
