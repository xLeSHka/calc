package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/pkg/config"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type JWT struct {
	secret []byte
}

func New(config config.Config) *JWT {
	return &JWT{secret: []byte(config.RandomSecret)}
}

func (j JWT) CreateToken(data jwt.MapClaims, workUntil time.Time) (string, error) {
	data["exp"] = workUntil.Unix()
	data["iat"] = time.Now().Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(j.secret)
}

func (j JWT) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token.Claims.(jwt.MapClaims), nil
}
func (j JWT) GenerateAccessToken(id uuid.UUID) (string, error) {
	return j.CreateToken(jwt.MapClaims{
		"id": id,
	}, time.Now().Add(time.Hour*24))
}
