package usersService

import (
	"context"
	"errors"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"gorm.io/gorm"
	"net/http"
)

func (s *UsersService) Register(user *models.User) (string, *customError.CustomError) {
	user, err := s.UsersRepository.Register(user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", customError.New(http.StatusConflict, err)
		}
		return "", customError.New(http.StatusInternalServerError, err)
	}
	token, err := s.JWT.GenerateAccessToken(user.ID)
	if err != nil {
		return "", customError.New(http.StatusInternalServerError, err)
	}
	err = s.RDB.SaveToken(context.Background(), user.ID, token)
	if err != nil {
		return "", customError.New(http.StatusInternalServerError, err)
	}
	return token, nil
}
