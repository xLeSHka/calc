package tests

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/utils/password"
	"io"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestGetExpressions(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	type TestGetExpression struct {
		Expected     []models.Expression
		jwt          string
		name         string
		expectedCode int
	}
	tests := []TestGetExpression{
		TestGetExpression{
			jwt:  profile1JWT,
			name: "Succes get expressions ",
			Expected: []models.Expression{
				Expr2,
				Expr,
			},
			expectedCode: http.StatusOK,
		},
		{
			jwt:          profile3JWT,
			name:         "unauthorized",
			Expected:     []models.Expression{},
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
			req, err := http.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
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
				var expressions []models.Expression
				data, err := io.ReadAll(w.Result().Body)
				assert.Nil(t, err)
				err = json.Unmarshal(data, &expressions)
				assert.Nil(t, err)
				for i := 0; i < len(test.Expected); i++ {

					assert.Equal(t, test.Expected[i].Status, expressions[i].Status)
					assert.Equal(t, *test.Expected[i].Result, *expressions[i].Result)
				}
			}
		})
	}
}
