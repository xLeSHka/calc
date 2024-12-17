package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	RestServerPort int `env:"SERVER_PORT" env-default:"9090"`
}

//функция конструктор для конфига
func New() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
