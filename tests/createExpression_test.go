package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/utils/password"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
	"time"
)

func TestCreateExpression(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	defer solveExpr()
	type TestCreateExpression struct {
		Expression   SendExpression
		jwt          string
		Expected     models.Expression
		name         string
		expectedCode int
	}
	tests := []TestCreateExpression{
		TestCreateExpression{
			name:       "Succes create expression",
			jwt:        profile1JWT,
			Expression: SendExpression{"2+2"},
			Expected: models.Expression{
				Expression: "2+2",
				Status:     "Solved",
				Result:     &r,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			name:       "Succes create expression 2",
			jwt:        profile1JWT,
			Expression: SendExpression{"log(18,18)^(-9)/3.14*(-12-3)*3/(-10)+2*sqrt(4)"},
			Expected: models.Expression{
				Expression: "log(18,18)^(-9)/3.14*(-12-3)*3/(-10)+2*sqrt(4)",
				Status:     "Solved",
				Result:     &r2,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			jwt:          profile1JWT,
			name:         "simple fail empty expression",
			Expression:   SendExpression{""},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			jwt:          profile1JWT,
			name:         "simple fail",
			Expression:   SendExpression{"1+1*"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			jwt:          profile1JWT,
			name:         "priority",
			Expression:   SendExpression{"2+2**2"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			jwt:          profile1JWT,
			name:         "right paranthes",
			Expression:   SendExpression{"((2+2-*(2"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			jwt:          profile1JWT,
			name:         "left paranthes",
			Expression:   SendExpression{"2+2)-2"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			jwt:        profile1JWT,
			name:       "division by zero",
			Expression: SendExpression{"10/0"},
			Expected: models.Expression{
				Expression: "10/0",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			jwt:          profile1JWT,
			name:         "invaid operator",
			Expression:   SendExpression{"10&0"},
			expectedCode: http.StatusUnprocessableEntity,
		},
		TestCreateExpression{
			jwt:        profile1JWT,
			name:       "log bad req",
			Expression: SendExpression{"log(-2,8)"},
			Expected: models.Expression{
				Expression: "log(-2,8)",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			jwt:        profile1JWT,
			name:       "log another bad req",
			Expression: SendExpression{"log(1,8)"},
			Expected: models.Expression{
				Expression: "log(1,8)",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			jwt:        profile1JWT,
			name:       "log another bad req 2",
			Expression: SendExpression{"log(16,(-1))"},
			Expected: models.Expression{
				Expression: "log(16,(-1))",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		TestCreateExpression{
			jwt:        profile1JWT,
			name:       "sqrt bad req",
			Expression: SendExpression{"sqrt(50-50-50)"},
			Expected: models.Expression{
				Expression: "sqrt(50-50-50)",
				Status:     "Unprocessable expression",
				Result:     nil,
			},
			expectedCode: http.StatusCreated,
		},
		{
			jwt:          profile3JWT,
			name:         "unknown user",
			Expression:   SendExpression{"2+2"},
			expectedCode: http.StatusUnauthorized,
		},
	}
	hash, err := password.Encrypt([]byte("passworD!2"), []byte(config.CryptoKey))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(hash, config.CryptoKey)
	UsersService.Register(&models.User{
		ID:       profile1.ID,
		Login:    profile1.Login,
		Password: hash,
	})
	RedisRepository.SaveToken(context.Background(), profile1.ID, profile1JWT)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			jsonData, _ := json.Marshal(test.Expression)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(jsonData))
			req.Header.Set("Authorization", "Bearer "+test.jwt)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusCreated {
				for i := 0; i < 10; i++ {
					solveExpr()
					time.Sleep(100 * time.Millisecond)
				}
				var expr models.Expression
				result := db.Model(&models.Expression{}).Where("expression = ?", test.Expected.Expression).Order("id DESC").First(&expr)
				assert.Nil(t, result.Error)
				assert.Equal(t, test.Expected.Status, expr.Status)
				if test.Expected.Status == "Solved" {
					assert.Equal(t, *test.Expected.Result, *expr.Result)
				} else {
					assert.Nil(t, expr.Result)
				}
			}
		})
	}
}
