package calculator

import (
	"errors"
	"github.com/xLeSHka/calc/orchestrator/internal/app/token"
	"testing"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "/",
			expression:     "1/2",
			expectedResult: 0.5,
		},
		{
			name:           "sqrt",
			expression:     "sqrt(64)",
			expectedResult: 8,
		},
		{
			name:           "log",
			expression:     "log(2,8)",
			expectedResult: 3,
		},
		{
			name:           "log(sqrt(),sqrt())",
			expression:     "log(sqrt(4),sqrt(64))",
			expectedResult: 3,
		},
		{
			name:           "big_expression",
			expression:     "log(18,18)^(-9)/3.14*(-12-3)*3/10+2*sqrt(4)",
			expectedResult: 2.56689,
		},
		{
			name:           "verybig",
			expression:     "log(18,18)^(-9)/3.14*(-12-3)*3/(-10)+2*sqrt(4)",
			expectedResult: 5.43311,
		},
		{
			name:           "omg",
			expression:     "log(sqrt(4),sqrt(64))^sqrt(81)*(-1)",
			expectedResult: -19683,
		},
	}
	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil {
				t.Errorf("successful case %s returns error", testCase.expression)
			}
			if val != testCase.expectedResult {
				t.Errorf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}
	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "simple1",
			expression:  "1+1*",
			expectedErr: token.ErrConvertRPNToNT,
		},
		{
			name:        "priority",
			expression:  "2+2**2",
			expectedErr: token.ErrConvertRPNToNT,
		},
		{
			name:        "right paranthes",
			expression:  "((2+2-*(2",
			expectedErr: token.ErrMissRightParanthesis,
		},
		{
			name:        "left paranthes",
			expression:  "2+2)-2",
			expectedErr: token.ErrMissLeftParanthesis,
		},
		{
			name:        "empty",
			expression:  "",
			expectedErr: token.ErrConvertRPNToNT,
		},
		{
			name:        "division by zero",
			expression:  "10/0",
			expectedErr: ErrDivisionByZero,
		},
		{
			name:        "invaid operator",
			expression:  "10&0",
			expectedErr: token.ErrUnknownSymbol,
		},
		{
			name:        "log bad req",
			expression:  "log(-2,8)",
			expectedErr: ErrLogNotDefinedFor,
		},
		{
			name:        "log another bad req",
			expression:  "log(1,8)",
			expectedErr: ErrLogNotDefinedFor,
		},
		{
			name:        "log another bad req",
			expression:  "log(16,(-1))",
			expectedErr: ErrLogOutOfFuncDomain,
		},
		{
			name:        "sqrt bad req",
			expression:  "sqrt(50-50-50)",
			expectedErr: ErrSqrtOutOfDomain,
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err == nil {
				t.Errorf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
			if err != nil {
				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("expected err: %v, recieved %v", testCase.expectedErr, err)
				}
			}
		})
	}
}
