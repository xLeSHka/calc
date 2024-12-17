package calculator

import "errors"

var (
	ErrInvalidExpression   = errors.New("invalid expression")
	ErrDivisionByZero      = errors.New("division by zero")
	ErrMissLeftParanthesis = errors.New("miss left paranthesis")
	ErrMissRightParanthesis = errors.New("miss right paranthesis")
)
