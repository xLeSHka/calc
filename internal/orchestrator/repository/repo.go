package repository

import "github.com/xLeSHka/calc/internal/models"

type Repo interface {
	CreateExpression(expression string) (int64, error)
	GetExpression(id int64) (*models.Expression, error)
	GetExpressions(size, page int) ([]*models.Expression, int64, error)
	UpdateExpression(expression *models.Expression) error
	SetResult(id int64, result float64) error
}
