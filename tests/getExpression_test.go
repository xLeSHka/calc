package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/utils/password"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestGetExpression(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	type TestGetExpression struct {
		ID           int64
		jwt          string
		Expected     models.Expression
		name         string
		expectedCode int
	}
	tests := []TestGetExpression{
		TestGetExpression{
			jwt:          profile1JWT,
			name:         "Succes get expression ",
			ID:           1111,
			Expected:     Expr,
			expectedCode: http.StatusOK,
		},
		TestGetExpression{
			jwt:          profile1JWT,
			name:         "Succes get expression 2",
			ID:           1112,
			Expected:     Expr2,
			expectedCode: http.StatusOK,
		},
		TestGetExpression{
			jwt:          profile1JWT,
			name:         "bad reqquest",
			ID:           -1,
			expectedCode: http.StatusBadRequest,
		},
		TestGetExpression{
			jwt:          profile1JWT,
			name:         "not found",
			ID:           math.MaxInt64,
			expectedCode: http.StatusNotFound,
		},
		{
			jwt:          profile3JWT,
			name:         "unauthorized",
			ID:           423423,
			expectedCode: http.StatusUnauthorized,
		},
	}
	hash, _ := password.Encrypt([]byte("passworD!2"), []byte(config.CryptoKey))
	UsersService.Register(&models.User{
		ID:       profile1.ID,
		Login:    profile1.Login,
		Password: hash,
	})
	RedisRepository.SaveToken(context.Background(), profile1.ID, profile1JWT)
	db.Create(&Expr)
	db.Create(&Expr2)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/api/v1/expressions/"+fmt.Sprintf("%d", test.ID), nil)
			assert.Nil(t, err)
			req.Header.Set("Authorization", "Bearer "+test.jwt)
			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var expression models.Expression
				data, err := io.ReadAll(w.Result().Body)
				assert.Nil(t, err)
				err = json.Unmarshal(data, &expression)
				assert.Nil(t, err)
				assert.Equal(t, test.Expected.ID, expression.ID)
				assert.Equal(t, test.Expected.Status, expression.Status)
				assert.Equal(t, *test.Expected.Result, *expression.Result)
			}
		})
	}
}
