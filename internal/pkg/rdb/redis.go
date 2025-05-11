package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"go.uber.org/fx"
	"strconv"
)

func New(config config.Config, lc fx.Lifecycle) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisHost + ":" + strconv.Itoa(int(config.RedisPort)),
		DB:   0,
	})

	err := rdb.Ping(context.Background()).Err()

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return rdb.Close()
		},
	})
	return rdb, nil
}
