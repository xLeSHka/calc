package repositoryAgent

import (
	"errors"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"gorm.io/gorm"
)

func (r *AgentRepository) UpdateExpression(expression *models.Expression) error {
	result := r.DB.Model(&models.Expression{}).
		Where("id = ?", expression.ID).
		Updates(&models.Expression{Status: expression.Status, Result: expression.Result})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return result.Error
	}
	return nil
}
