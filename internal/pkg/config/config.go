package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	ServerPort         int           `env:"SERVER_PORT" env-default:"9090"`
	ServerHost         string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	ComputingPower     int           `env:"COMPUTING_POWER" env-default:"5"`
	TimeAddiction      time.Duration `env:"TIME_ADDITION_MS" env-default:"2000ms"`
	TimeSubtraction    time.Duration `env:"TIME_SUBTRACTION_MS" env-default:"2000ms"`
	TimeMultiplication time.Duration `env:"TIME_MULTIPLICATION_MS" env-default:"2000ms"`
	TimeDivision       time.Duration `env:"TIME_DIVISION_MS" env-default:"2000ms"`
	TimeExponentiation time.Duration `env:"TIME_EXPONENTIATION_MS" env-default:"2000ms"`
	TimeUnaryMinus     time.Duration `env:"TIME_UNARY_MINUS_MS" env-default:"2000ms"`
	TimeLogarithm      time.Duration `env:"TIME_LOGARITHM_MS" env-default:"2000ms"`
	TimeSquareRoot     time.Duration `env:"TIME_SQUARE_ROOT_MS"env-default:"2000ms"`
	PostgresUser       string        `env:"POSTGRES_USER" env-default:"root"`
	PostgresPassword   string        `env:"POSTGRES_PASSWORD" env-default:"123"`
	PostgresDB         string        `env:"POSTGRES_DB" env-default:"testdatabase"`
	PostgresHost       string        `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort       string        `env:"POSTGRES_PORT" env-default:"5432"`
}

// функция конструктор для конфига
func New() (Config, error) {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
