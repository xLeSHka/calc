package token

import "errors"

var (
	ErrSqrtOutOfDomain           = errors.New("sqrt not defined for a < 0")
	ErrLogOutOfFuncDomain        = errors.New("log(a,x) out of function`s domain")
	ErrLogNotDefinedFor          = errors.New("log not defined for a <= 0 or a == 1")
	ErrAssociativityReq          = errors.New("associativity requeried")
	ErrNonOperationTokenAsc      = errors.New("non-operation token cant have an associativity")
	ErrUnknownOperator           = errors.New("unknown operator")
	ErrUnknownSymbol             = errors.New("unknown symbol")
	ErrUnexpectedSymbol          = errors.New("unexpected symbol")
	ErrNonBalancedParanthesis    = errors.New("non-balanced on paranthesis expression")
	ErrMisedSepOrParanth         = errors.New("paranthesis or separator missed")
	ErrConvertRPNToNT            = errors.New("failed convert RPN to node tree")
	ErrDivisionByZero            = errors.New("division by zero")
	ErrFailedCalculateExpression = errors.New("failed calculate expression")
)
