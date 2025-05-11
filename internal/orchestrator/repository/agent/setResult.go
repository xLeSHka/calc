package repositoryAgent

import (
	"errors"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"gorm.io/gorm"
)

func (r *AgentRepository) SetResult(id int64, result float64) error {
	res := r.DB.Model(&models.Expression{}).
		Where("id = ?", id).
		Updates(&models.Expression{Result: &result})
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return repository.ErrNotFound
		}
		return res.Error
	}
	return nil
}
