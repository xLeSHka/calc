package calculator

import "errors"

var (
	ErrInvalidExpression   = errors.New("invalid expression")
	ErrInvalidExpression1   = errors.New("invalid expression1")
	ErrInvalidExpression2   = errors.New("invalid expression2")
	ErrInvalidExpression3   = errors.New("invalid expression3")
	ErrDivisionByZero      = errors.New("division by zero")
	ErrMissLeftParanthesis = errors.New("miss left paranthesis")
	ErrMissRightParanthesis = errors.New("miss right paranthesis")
)
