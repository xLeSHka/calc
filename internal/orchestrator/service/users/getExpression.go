package usersService

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"net/http"
)

func (s *UsersService) GetExpression(id int64, userID uuid.UUID) (*models.Expression, *customError.CustomError) {
	expr, err := s.UsersRepository.GetExpression(id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, customError.New(http.StatusNotFound, fmt.Errorf("UsersService.GetExpression: error: %w", err))
		}
		return nil, customError.New(http.StatusInternalServerError, fmt.Errorf("UsersService.GetExpression: unknown errorerror: %w", err))
	}
	return expr, nil
}
