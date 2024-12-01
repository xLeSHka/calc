package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	ServerPort int `env:"SERVER_PORT" env-default:"9090"`
}
//функция конструктор для конфига
func New() *Config {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil
	}
	return &cfg
}
