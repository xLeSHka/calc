package repositoryUsers

import (
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/models"
	"gorm.io/gorm/clause"
)

func (r *UsersRepository) CreateExpression(expression string, userID uuid.UUID) (int64, error) {
	expr := &models.Expression{
		Expression: expression,
		UserID:     userID,
		Status:     "Waiting",
	}
	result := r.DB.Create(expr).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})
	if result.Error != nil {
		return 0, result.Error
	}
	return expr.ID, nil
}
