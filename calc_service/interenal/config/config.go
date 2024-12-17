package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	GRPCServerPort int `env:"GRPC_SERVER_PORT" env-default:"9090"`
	RestServerPort int `env:"RES_SERVER_PORT" env-default:"9090"`
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
