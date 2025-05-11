package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/internal/models"
	"net/http"
	"net/http/httptest"
	"syscall"
	"testing"
)

func TestRegister(t *testing.T) {
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
			name:         "register user",
			Login:        "user_1",
			Password:     "passworD!2",
			expectedCode: http.StatusOK,
		},
		{
			name:         "register user2",
			Login:        "user_2",
			Password:     "passworD!2",
			expectedCode: http.StatusOK,
		},
		{
			name:         "try register user",
			Login:        "user_1",
			Password:     "passworD!2",
			expectedCode: http.StatusConflict,
		},
		{
			name:         "try register user with bad password",
			Login:        "user_3",
			Password:     "password12",
			expectedCode: http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test)
			url := "/api/v1/register"
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			http2.ServeHTTP(w, req)
			defer func() {
				err := w.Result().Body.Close()
				assert.Nil(t, err)
			}()
			assert.Equal(t, test.expectedCode, w.Code)
			if test.expectedCode == http.StatusOK {
				var user models.User
				err := db.Model(&models.User{}).First(&user, "login = ?", test.Login).Error
				assert.Nil(t, err)
			}
		})
	}
}
