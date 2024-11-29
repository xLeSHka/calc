package config

type Config struct {
	Width  int
	Height int
}

func New() *Config {
	return &Config{Width: 100, Height: 50}
}
