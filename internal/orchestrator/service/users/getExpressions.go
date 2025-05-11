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

func (s *UsersService) GetExpressions(size, page int, userID uuid.UUID) ([]*models.Expression, int64, *customError.CustomError) {
	exprs, total, err := s.UsersRepository.GetExpressions(size, page, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, 0, customError.New(http.StatusNotFound, fmt.Errorf("UsersService.GetExpressions: error: %w", err))
		}
		return nil, 0, customError.New(http.StatusInternalServerError, fmt.Errorf("UsersService.GetExpressions: unknown errorerror: %w", err))
	}
	return exprs, total, nil
}
