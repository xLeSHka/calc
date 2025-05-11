package usersService

import (
	"context"
	"errors"
	"fmt"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/utils/password"
	"gorm.io/gorm"
	"net/http"
)

func (s *UsersService) Login(login, pw string) (string, *customError.CustomError) {
	user, err := s.UsersRepository.Login(login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", customError.New(http.StatusNotFound, fmt.Errorf("User Not Found"))
		}
		return "", customError.New(http.StatusInternalServerError, err)
	}
	err = password.Compare([]byte(pw), user.Password, s.cryptoKey)
	if err != nil {
		return "", customError.New(http.StatusUnauthorized, fmt.Errorf("Invalid password"))
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
