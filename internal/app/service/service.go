package service

import (
	"errors"
	"fmt"
	"github.com/xLeSHka/calc/internal/app/repository"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/pkg/calculator"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/pkg/token"
	"go.uber.org/zap"
	"net/http"
)

type Service struct {
	Log        *zap.Logger
	Repository repository.Repo
	Calculator *calculator.Calculator
}

func New(
	repo repository.Repo,
	log *zap.Logger,
	calculator *calculator.Calculator,
) *Service {
	return &Service{
		Log:        log,
		Repository: repo,
		Calculator: calculator,
	}
}

func (s *Service) GetExpressions(size, page int) ([]*models.Expression, int64, *customError.CustomError) {
	exprs, total, err := s.Repository.GetExpressions(size, page)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, 0, customError.New(http.StatusNotFound, fmt.Errorf("Service.GetExpressions: error: %w", err))
		}
		return nil, 0, customError.New(http.StatusInternalServerError, fmt.Errorf("Service.GetExpressions: unknown errorerror: %w", err))
	}
	return exprs, total, nil
}
func (s *Service) GetExpression(id int64) (*models.Expression, *customError.CustomError) {
	expr, err := s.Repository.GetExpression(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, customError.New(http.StatusNotFound, fmt.Errorf("Service.GetExpression: error: %w", err))
		}
		return nil, customError.New(http.StatusInternalServerError, fmt.Errorf("Service.GetExpression: unknown errorerror: %w", err))
	}
	return expr, nil
}
func (s *Service) CreateExpression(expression string) (int64, *customError.CustomError) {
	expressionNT, err := token.TokenizeExpression(expression)
	if err != nil {
		return 0, customError.New(http.StatusUnprocessableEntity, fmt.Errorf("Service.CreateExpression: error: %w", err))
	}
	id, err := s.Repository.CreateExpression(expression)
	if err != nil {
		return 0, customError.New(http.StatusInternalServerError, fmt.Errorf("Service.CreateExpression: error: %w", err))
	}
	go s.Calculator.Calc(*expressionNT, id)
	return id, nil
}
func (s *Service) SetResult(id, expressionID int64, result float64) *customError.CustomError {
	task := &models.Task{
		ID:           id,
		ExpressionID: expressionID,
		Result:       &result,
	}
	s.Calculator.Results <- task
	return nil
}
func (s *Service) GetTask() (*models.Task, *customError.CustomError) {
	select {
	case task, ok := <-s.Calculator.Tasks:
		if !ok {
			return nil, customError.New(http.StatusInternalServerError, fmt.Errorf("Service.GetTask: task channel closed"))
		}
		return task, nil
	default:
		return nil, customError.New(http.StatusNotFound, fmt.Errorf("Service.GetTask: task not found"))
	}
}
