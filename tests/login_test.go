package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/utils/password"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestLogin(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	type Test struct {
		name         string
		Login        string `json:"login"`
		Password     string `json:"password"`
		expectedCode int
	}
	tests := []Test{
		{
			name:         "login user",
			Login:        profile1.Login,
			Password:     "passworD!2",
			expectedCode: http.StatusOK,
		},
		{
			name:         "login user2",
			Login:        profile2.Login,
			Password:     "passworD!2",
			expectedCode: http.StatusOK,
		},
		{
			name:         "try login unknown user",
			Login:        "user_3",
			Password:     "passworD!2",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "try login user with bad password",
			Login:        "user_1",
			Password:     "password12",
			expectedCode: http.StatusBadRequest,
		},
	}
	hash, _ := password.Encrypt([]byte("passworD!2"), []byte(config.CryptoKey))
	UsersService.Register(&models.User{
		ID:       profile1.ID,
		Login:    profile1.Login,
		Password: hash,
	})
	UsersService.Register(&models.User{
		ID:       profile2.ID,
		Login:    profile2.Login,
		Password: hash,
	})
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test)
			url := "/api/v1/login"
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
		})
	}
}
