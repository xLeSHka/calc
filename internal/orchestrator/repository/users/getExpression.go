package repositoryUsers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"gorm.io/gorm"
)

func (r *UsersRepository) GetExpression(id int64, userID uuid.UUID) (*models.Expression, error) {
	var expression models.Expression
	result := r.DB.First(&expression, &models.Expression{ID: id, UserID: userID})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, result.Error
	}
	return &expression, nil
}
