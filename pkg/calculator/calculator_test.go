package calculator

import (
	"errors"
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
			name:        "simple",
			expression:  "1+1*",
			expectedErr: ErrInvalidExpression,
		},
		{
			name:        "priority",
			expression:  "2+2**2",
			expectedErr: ErrInvalidExpression,
		},
		{
			name:        "right paranthes",
			expression:  "((2+2-*(2",
			expectedErr: ErrMissRightParanthesis,
		},
		{
			name:        "left paranthes",
			expression:  "2+2)-2",
			expectedErr: ErrMissLeftParanthesis,
		},
		{
			name:        "empty",
			expression:  "",
			expectedErr: ErrInvalidExpression,
		},
		{
			name:        "division by zero",
			expression:  "10/0",
			expectedErr: ErrDivisionByZero,
		},
		{
			name:        "invaid operator",
			expression:  "10&0",
			expectedErr: ErrInvalidExpression,
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
