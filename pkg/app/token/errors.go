package token

import "errors"

var (
	ErrAssociativityReq          = errors.New("associativity requeried")
	ErrNonOperationTokenAsc      = errors.New("non-operation token cant have an associativity")
	ErrUnknownOperator           = errors.New("unknown operator")
	ErrUnknownSymbol             = errors.New("unknown symbol")
	ErrUnexpectedSymbol          = errors.New("unexpected symbol")
	ErrMissLeftParanthesis       = errors.New("left paranthesis missed")
	ErrMisedSepOrParanth         = errors.New("paranthesis or separator missed")
	ErrConvertRPNToNT            = errors.New("failed convert RPN to node tree")
	ErrTokenIsNotAnOperator      = errors.New("Token is not an operator")
	ErrFailedCalculateExpression = errors.New("failed calculate expression")
	ErrMissRightParanthesis      = errors.New("right paranthesis missed")
)
