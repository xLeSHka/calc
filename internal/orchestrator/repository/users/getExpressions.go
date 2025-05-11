package repositoryUsers

import (
	"errors"
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"gorm.io/gorm"
)

func (r *UsersRepository) GetExpressions(size, page int, userID uuid.UUID) ([]*models.Expression, int64, error) {
	var expressions []*models.Expression
	tx := r.DB.Begin()
	result := tx.Model(&models.Expression{}).
		Limit(size).
		Where("user_id = ?", userID).
		Offset(page * size).
		Order("id desc").
		Find(&expressions)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, 0, repository.ErrNotFound
		}
		tx.Rollback()
		return nil, 0, result.Error
	}
	var count int64
	result = tx.Model(&models.Expression{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		tx.Rollback()
		return nil, 0, result.Error
	}
	return expressions, count, nil
}
