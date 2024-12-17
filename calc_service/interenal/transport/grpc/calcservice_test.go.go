package grpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xLeSHka/calculator/pkg/calculator"
)

func TestCalcHandlerSuccessCase(t *testing.T) {
	testCasesSucces := []struct {
		name           string
		Expression     string `json:"expression"`
		expectedResult float64
		expectedStatus int
		expectedErr    error
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
			name:           "simple",
			Expression:     "1+1*",
			expectedErr:    calculator.ErrInvalidExpression,
			expectedStatus: 400,
		},
		{
			name:           "priority",
			Expression:     "2+2**2",
			expectedErr:    calculator.ErrInvalidExpression,
			expectedStatus: 400,
		},
		{
			name:           "right paranthes",
			Expression:     "((2+2-*(2",
			expectedErr:    calculator.ErrMissRightParanthesis,
			expectedStatus: 400,
		},
		{
			name:           "left paranthes",
			Expression:     "2+2)-2",
			expectedErr:    calculator.ErrMissLeftParanthesis,
			expectedStatus: 400,
		},
		{
			name:           "empty",
			Expression:     "",
			expectedErr:    calculator.ErrInvalidExpression,
			expectedStatus: 400,
		},
		{
			name:           "division by zero",
			Expression:     "10/0",
			expectedErr:    calculator.ErrDivisionByZero,
			expectedStatus: 400,
		},
		{
			name:           "invaid operator",
			Expression:     "10&0",
			expectedErr:    calculator.ErrInvalidExpression,
			expectedStatus: 400,
		},
	}
	// testLogger := logger.New()
	// ctx := context.Background()

	for _, testCase := range testCasesSucces {
		t.Run(testCase.name, func(t *testing.T) {
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(testCase)
			if err != nil {
				t.Errorf("failed encode test cases", err)
				return
			}
			client := http.Client{}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", &b)
			resp, err := client.Do(req)
			if resp.StatusCode != testCase.expectedStatus {
				t.Errorf("Resp status: expected %v, recieved %v", testCase.expectedStatus, resp.StatusCode)
				return
			}
			if err != nil {
				if err.Error() != testCase.expectedErr.Error() {
					t.Errorf("Error message: expected %v, recieved %v", testCase.expectedErr.Error(), err.Error())
					return
				}
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Error: %v", err)
				return
			}
			var Response Response
			err = json.Unmarshal(data, &Response)
			if err != nil {
				t.Errorf("Error: %v", err)
			}
			if Response.Result != testCase.expectedResult {
				t.Errorf("expected status %.2f, recieved %2f", testCase.expectedResult, Response.Result)
			}

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