package redisRepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"strconv"
	"time"
)

type RedisRepository struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) repository.RedisRepository {
	return &RedisRepository{rdb: rdb}
}
func (c *RedisRepository) SaveToken(ctx context.Context, userID uuid.UUID, token string) error {
	err := c.rdb.Set(ctx, "jwt:"+userID.String(), token, time.Hour*24).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c *RedisRepository) DeleteToken(ctx context.Context, userID uuid.UUID) error {
	err := c.rdb.Del(ctx, "jwt:"+userID.String()).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c *RedisRepository) SaveID(ctx context.Context, userID uuid.UUID, id int64) error {
	err := c.rdb.Set(ctx, "id:"+userID.String(), strconv.FormatInt(id, 10), time.Hour*24).Err()
	if err != nil {
		return err
	}
	return nil
}
func (c *RedisRepository) GetID(ctx context.Context, userID uuid.UUID) (int64, error) {
	res := c.rdb.Get(ctx, "id:"+userID.String())
	if res.Err() != nil {
		return 0, res.Err()
	}
	return strconv.ParseInt(res.Val(), 10, 64)
}
