package config

type Config struct {
	Port string
	DSN  string
}

func GetConfig() (*Config, error) {

	cfg := new(Config)
	cfg.Port = "localhost:8080"
	cfg.DSN = "file:resources\\market.db"

	return cfg, nil
}
