package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	ServerPort         int           `env:"SERVER_PORT" env-default:"9090"`
	ServerHost         string        `env:"SERVER_HOST" env-default:"0.0.0.0"`
	ComputingPower     int           `env:"COMPUTING_POWER" env-default:"5"`
	TimeAddiction      time.Duration `env:"TIME_ADDITION_MS" env-default:"1ms"`
	TimeSubtraction    time.Duration `env:"TIME_SUBTRACTION_MS" env-default:"1ms"`
	TimeMultiplication time.Duration `env:"TIME_MULTIPLICATION_MS" env-default:"1ms"`
	TimeDivision       time.Duration `env:"TIME_DIVISION_MS" env-default:"1ms"`
	TimeExponentiation time.Duration `env:"TIME_EXPONENTIATION_MS" env-default:"1ms"`
	TimeUnaryMinus     time.Duration `env:"TIME_UNARY_MINUS_MS" env-default:"1ms"`
	TimeLogarithm      time.Duration `env:"TIME_LOGARITHM_MS" env-default:"1ms"`
	TimeSquareRoot     time.Duration `env:"TIME_SQUARE_ROOT_MS"env-default:"1ms"`
	PostgresUser       string        `env:"POSTGRES_USER" env-default:"root"`
	PostgresPassword   string        `env:"POSTGRES_PASSWORD" env-default:"123"`
	PostgresDB         string        `env:"POSTGRES_DB" env-default:"testdatabase"`
	PostgresHost       string        `env:"POSTGRES_HOST" env-default:"localhost"`
	PostgresPort       string        `env:"POSTGRES_PORT" env-default:"5432"`
	RedisHost          string        `env:"REDIS_HOST" env-default:"localhost"`
	RedisPort          int           `env:"REDIS_PORT" env-default:"6379"`
	RandomSecret       string        `env:"RANDOM_SECRET" env-default:"secret"`
	CryptoKey          string        `env:"CRYPTO_KEY" default:"12345678901234567890123456789012"`
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
