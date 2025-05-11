package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/models"
)

type RedisRepository interface {
	SaveToken(ctx context.Context, userID uuid.UUID, token string) error
	DeleteToken(ctx context.Context, userID uuid.UUID) error
	SaveID(ctx context.Context, userID uuid.UUID, id int64) error
	GetID(ctx context.Context, userID uuid.UUID) (int64, error)
}
type AgentRepository interface {
	UpdateExpression(expression *models.Expression) error
	SetResult(id int64, result float64) error
}
type UsersRepository interface {
	Register(person *models.User) (*models.User, error)
	Login(login string) (*models.User, error)
	CreateExpression(expression string, userID uuid.UUID) (int64, error)
	GetExpression(id int64, userID uuid.UUID) (*models.Expression, error)
	GetExpressions(size, page int, userID uuid.UUID) ([]*models.Expression, int64, error)
}
