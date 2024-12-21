package calculator

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	"github.com/xLeSHka/calc/pkg/app/token"
)

func checkedDivisionByZero(a, b float64) (float64, error) {
	if b == 0 {
		return 0, token.ErrDivisionByZero
	}
	return a / b, nil
}

func calculate(node *token.Node, wg *sync.WaitGroup, cancelCtx context.CancelFunc) error {
	defer wg.Done()
	if node == nil {
		return nil
	}
	wg1 := &sync.WaitGroup{}
	if node.Left != nil {
		wg1.Add(1)
		go calculate(node.Left, wg1, cancelCtx)
	}
	if node.Right != nil {
		wg1.Add(1)
		go calculate(node.Right, wg1, cancelCtx)
	}
	if node.Left == nil && node.Right == nil {
		return nil
	}
	wg1.Wait()
	switch node.Token.Type {
	case token.Operator:
		switch node.Token.Associativity {
		case token.Left:
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				return err
			}
			b, err := strconv.ParseFloat(node.Right.Token.Token, 64)
			if err != nil {
				return err
			}
			if node.Token.Token == "+" {
				node.Token.Token = strconv.FormatFloat(a+b, 'f', 5, 64)
			} else if node.Token.Token == "-" {
				node.Token.Token = strconv.FormatFloat(a-b, 'f', 5, 64)
			} else if node.Token.Token == "*" {
				node.Token.Token = strconv.FormatFloat(a*b, 'f', 5, 64)
			} else if node.Token.Token == "/" {
				res, err := checkedDivisionByZero(a, b)
				if err != nil {
					return err
				}
				node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
			} else if node.Token.Token == "^" {
				node.Token.Token = strconv.FormatFloat(math.Pow(a, b), 'f', 5, 64)
			} else {
				return token.ErrUnknownOperator
			}
		case token.Right:
			if node.Token.Token == "-" {
				node.Token.Token = "-" + node.Left.Token.Token
			} else {
				return token.ErrUnknownOperator
			}
		}
	case token.Function:
		if node.Token.Token == "log" {
			b, err := strconv.ParseFloat(node.Right.Token.Token, 64)
			if err != nil {
				return err
			}
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				return err
			}
			if a <= 0 || a == 1 {
				return token.ErrLogNotDefinedFor
			}
			if b <= 0.0 {
				return token.ErrLogOutOfFuncDomain
			}
			node.Token.Token = strconv.FormatFloat(math.Log(b)/math.Log(a), 'f', 5, 64)
		}
		if node.Token.Token == "sqrt" {
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				return err
			}
			if a < 0.0 {
				return token.ErrSqrtOutOfDomain
			}
			node.Token.Token = strconv.FormatFloat(math.Sqrt(a), 'f', 5, 64)
		}
	}
	return nil
}

func Calc(expression string) (float64, error) {
	expressionBT, err := token.TokenizeExpression(expression)
	if err != nil {
		return 0, err
	}

	ctx, cancelCtx := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	// errCh := make(chan error)
	// tokens, err := token.Tokenize(expression)
	// if err != nil {
	// 	fmt.Println("error 1")
	// 	return 0, err
	// }
	go  calculate(expressionBT, wg, cancelCtx)
	select {
	case <-ctx.Done():
		return 0, fmt.Errorf("context done")
	default:
		wg.Wait()
		res, err := strconv.ParseFloat(expressionBT.Token.Token, 64)
		if err != nil {
			return 0, err
		}
		return res, nil
	}

}
