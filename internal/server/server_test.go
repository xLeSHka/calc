package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/pkg/app/calculator"
	"github.com/xLeSHka/calc/pkg/app/token"
)

func TestCalcHandlerSuccessCase(t *testing.T) {
	testCasesSucces := []struct {
		name           string
		Expression     string `json:"expression"`
		expectedResult float64
		expectedStatus int
		expectedErrMsg string
		method         string
	}{
		{
			name:           "simple",
			Expression:     "1+1",
			expectedResult: 2,
			expectedStatus: 200,
		},
		{
			name:           "priority",
			Expression:     "(2+2)*2",
			expectedResult: 8,
			expectedStatus: 200,
		},
		{
			name:           "priority",
			Expression:     "2+2*2",
			expectedResult: 6,
			expectedStatus: 200,
		},
		{
			name:           "/",
			Expression:     "1/2",
			expectedResult: 0.5,
			expectedStatus: 200,
		},
		{
			name:           "simple1",
			Expression:     "1+1*",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", token.ErrConvertRPNToNT),
			expectedStatus: 422,
		},
		{
			name:           "priority",
			Expression:     "2+2**2",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", token.ErrConvertRPNToNT),
			expectedStatus: 422,
		},
		{
			name:           "right paranthes",
			Expression:     "((2+2-*(2",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", token.ErrMissRightParanthesis),
			expectedStatus: 422,
		},
		{
			name:           "left paranthes",
			Expression:     "2+2)-2",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", token.ErrMissLeftParanthesis),
			expectedStatus: 422,
		},
		{
			name:           "empty",
			Expression:     "",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", token.ErrConvertRPNToNT),
			expectedStatus: 422,
		},
		{
			name:           "division by zero",
			Expression:     "10/0",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", calculator.ErrDivisionByZero),
			expectedStatus: 422,
		},
		{
			name:           "invaid operator",
			Expression:     "10&0",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", token.ErrUnknownSymbol),
			expectedStatus: 422,
		},
		{
			name:           "internal",
			Expression:     "internal",
			expectedErrMsg: "code=500, message=Internal server error",
			expectedStatus: 500,
		},
		{
			name:           "log bad req",
			Expression:     "log(-2,8)",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", calculator.ErrLogNotDefinedFor),
			expectedStatus: 422,
		},
		{
			name:           "log another bad req",
			Expression:     "log(1,8)",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", calculator.ErrLogNotDefinedFor),
			expectedStatus: 422,
		},
		{
			name:           "log another bad req",
			Expression:     "log(16,(-1))",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s", calculator.ErrLogOutOfFuncDomain),
			expectedStatus: 422,
		},
		{
			name:           "sqrt bad req",
			Expression:     "sqrt(50-50-50)",
			expectedErrMsg: fmt.Sprintf("code=422, message=%s",calculator.ErrSqrtOutOfDomain),
			expectedStatus: 422,
		},
	}
	// testLogger := logger.New()
	// ctx := context.Background()

	for _, testCase := range testCasesSucces {
		t.Run(testCase.name, func(t *testing.T) {
			//создание body запроса
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(testCase)
			if err != nil {
				t.Error("failed encode test cases", err)
				return
			}
			//настраиваем запрос
			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", &b)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := echo.New().NewContext(req, rec)

			//создаем сервер для отправки запросов
			server := Server{
				server: e,
			}
			//делаем запрос на сервер
			err = server.calculate(c)
			resp := Response{}
			//проверяем если вернулась ошибка, ожидаемая ли она или нет
			if testCase.expectedStatus != 200 {
				recCode := strings.Split(err.Error(), "=")
				recCode = strings.Split(recCode[1], ",")
				recievedCode, err2 := strconv.Atoi(recCode[0])
				if err2 != nil {
					t.Errorf("unexpected error: %v", err)
				}
				//сравниваем полученную ошибку и ожидаемую
				assert.Equal(t, testCase.expectedErrMsg, err.Error())
				assert.Equal(t, testCase.expectedStatus, recievedCode)
			}
			//десериализируем тело ответа
			err = json.NewDecoder(rec.Body).Decode(&resp)
			if err != nil && err != io.EOF {
				t.Errorf("unexpected error: %v", err)
				return
			}
			//проверяем значение с ожидаемым
			assert.Equal(t, testCase.expectedResult, resp.Result)

		})
	}
}

// func TestCalcHandlerFailCase(t *testing.T) {
// 	testCasesFail := []struct {
// 		name           string
// 		Expression     string `json:"expression"`
// 		expectedErr    error
// 		expectedStatus int
// 	}{}
// 	for _, testCase := range testCasesFail {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			var b bytes.Buffer
// 			err := json.NewEncoder(&b).Encode(testCase)
// 			if err != nil {
// 				t.Error(err)
// 			}
// 			w := httptest.NewRecorder()
// 			req := httptest.NewRequest(http.MethodPost, "/", &b)
// 			CalcHandler(
// 			res := w.Result()
// 			defer res.Body.Close()
// 			data, err := io.ReadAll(res.Body)
// 			if err != nil {
// 				t.Errorf("Error: %v", err)
// 			}
// 			if string(data) != (testCase.expectedErr.Error() + "\n") {
// 				t.Errorf("expected error message %s, recieved %s", testCase.expectedErr.Error(), string(data))
// 			}
// 			if res.StatusCode != testCase.expectedStatus {
// 				t.Errorf("expected status %d, recieved %d", testCase.expectedStatus, res.StatusCode)
// 			}

// 		})
// 	}
// }
