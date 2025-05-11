package service

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/pkg/customError"
)

type UsersService interface {
	GetExpressions(size, page int, userID uuid.UUID) ([]*models.Expression, int64, *customError.CustomError)
	GetExpression(id int64, userID uuid.UUID) (*models.Expression, *customError.CustomError)
	CreateExpression(expression string, userID uuid.UUID) (int64, *customError.CustomError)
	Register(user *models.User) (string, *customError.CustomError)
	Login(login, password string) (string, *customError.CustomError)
}
type AgentService interface {
	SetResult(id, expressionID int64, result float64, error *string) *customError.CustomError
	GetTask() (*models.Task, *customError.CustomError)
}
