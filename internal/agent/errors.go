package agent

import "errors"

var (
	ErrDivisionByZero     = errors.New("division by zero")
	ErrSqrtOutOfDomain    = errors.New("sqrt not defined for a < 0")
	ErrLogOutOfFuncDomain = errors.New("log(a,x) out of function`s domain")
	ErrLogNotDefinedFor   = errors.New("log not defined for a <= 0 or a == 1")
)
