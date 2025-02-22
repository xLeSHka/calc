package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	ServerPort         int           `env:"SERVER_PORT" env-default:"9090"`
	ServerHost         string        `env:"SERVER_HOST" env-default:"localhost"`
	ComputingPower     int           `env:"COMPUTING_POWER" env-default:"5"`
	TimeAddiction      time.Duration `env:"TIME_ADDITION_MS"`
	TimeSubtraction    time.Duration `env:"TIME_SUBTRACTION_MS"`
	TimeMultiplication time.Duration `env:"TIME_MULTIPLICATION_MS"`
	TimeDivision       time.Duration `env:"TIME_DIVISION_MS"`
	TimeExponentiation time.Duration `env:"TIME_EXPONENTIATION_MS"`
	TimeUnaryMinus     time.Duration `env:"TIME_UNARY_MINUS_MS"`
	TimeLogarithm      time.Duration `env:"TIME_LOGARITHM_MS"`
	TimeSquareRoot     time.Duration `env:"TIME_SQUARE_ROOT_MS"`
}

// функция конструктор для конфига
func New() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
