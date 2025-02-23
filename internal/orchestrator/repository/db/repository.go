package db

import (
	"errors"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}
func (r *Repository) CreateExpression(expression string) (int64, error) {
	expr := &models.Expression{
		Expression: expression,
		Status:     "Waiting",
	}
	result := r.DB.Create(expr).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})
	if result.Error != nil {
		return 0, result.Error
	}
	return expr.ID, nil
}
func (r *Repository) GetExpression(id int64) (*models.Expression, error) {
	var expression models.Expression
	result := r.DB.First(&expression, &models.Expression{ID: id})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound
		}
		return nil, result.Error
	}
	return &expression, nil
}
func (r *Repository) GetExpressions(size, page int) ([]*models.Expression, int64, error) {
	var expressions []*models.Expression
	tx := r.DB.Begin()
	result := tx.Model(&models.Expression{}).
		Limit(size).
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
	result = tx.Model(&models.Expression{}).Count(&count)
	if result.Error != nil {
		tx.Rollback()
		return nil, 0, result.Error
	}
	return expressions, count, nil
}
func (r *Repository) UpdateExpression(expression *models.Expression) error {
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
func (r *Repository) SetResult(id int64, result float64) error {
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
