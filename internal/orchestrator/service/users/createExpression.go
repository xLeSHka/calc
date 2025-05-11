package usersService

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/pkg/token"
	"net/http"
)

func (s *UsersService) CreateExpression(expression string, userID uuid.UUID) (int64, *customError.CustomError) {
	expressionNT, err := token.TokenizeExpression(expression)
	if err != nil {
		return 0, customError.New(http.StatusUnprocessableEntity, fmt.Errorf("UsersService.CreateExpression: error: %w", err))
	}

	id, err := s.UsersRepository.CreateExpression(expression, userID)
	if err != nil {
		return 0, customError.New(http.StatusInternalServerError, fmt.Errorf("UsersService.CreateExpression: error: %w", err))
	}
	go s.Calculator.Calc(*expressionNT, id)
	return id, nil
}
