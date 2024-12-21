package calculator

import (
	"context"
	"math"
	"strconv"

	"github.com/xLeSHka/calc/pkg/app/token"
	"golang.org/x/sync/errgroup"
)

func checkedDivisionByZero(a, b float64) (float64, error) {
	if b == 0 {
		return 0, token.ErrDivisionByZero
	}
	return a / b, nil
}

func calculate(node *token.Node, ctx context.Context) error {
	if node == nil {
		return nil
	}
	eg, ctx2 := errgroup.WithContext(context.Background())

	if node.Left != nil {
		eg.Go(func() error {
			return calculate(node.Left, ctx2)
		})
	}
	if node.Right != nil {
		eg.Go(func() error {
			return calculate(node.Right, ctx2)
		})
	}

	if node.Left == nil && node.Right == nil {
		return nil
	}
	if err := eg.Wait(); err != nil {
		ctx.Done()
		return err
	}
	switch node.Token.Type {
	case token.Operator:
		switch node.Token.Associativity {
		case token.Left:
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			b, err := strconv.ParseFloat(node.Right.Token.Token, 64)
			if err != nil {
				ctx.Done()
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
					ctx.Done()
					return token.ErrDivisionByZero
				}
				node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
			} else if node.Token.Token == "^" {
				node.Token.Token = strconv.FormatFloat(math.Pow(a, b), 'f', 5, 64)
			} else {
				ctx.Done()
				return token.ErrUnknownOperator

			}
		case token.Right:
			if node.Token.Token == "-" {
				node.Token.Token = "-" + node.Left.Token.Token
			} else {
				ctx.Done()
				return token.ErrUnknownOperator
			}
		}
	case token.Function:
		if node.Token.Token == "log" {
			b, err := strconv.ParseFloat(node.Right.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			if a <= 0 || a == 1 {
				ctx.Done()
				return token.ErrLogNotDefinedFor
			}
			if b <= 0.0 {
				ctx.Done()
				return token.ErrLogOutOfFuncDomain
			}
			node.Token.Token = strconv.FormatFloat(math.Log(b)/math.Log(a), 'f', 5, 64)
		}
		if node.Token.Token == "sqrt" {
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			if a < 0.0 {
				ctx.Done()
				return token.ErrSqrtOutOfDomain
			}
			node.Token.Token = strconv.FormatFloat(math.Sqrt(a), 'f', 5, 64)
		}
	}
	return nil
}

func Calc(expression string) (float64, error) {
	expressionNT, err := token.TokenizeExpression(expression)
	if err != nil {
		return 0, err
	}
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return calculate(expressionNT, ctx)
	})
	if err := eg.Wait(); err != nil {
		// logger.New().Error(ctx, "error calculate expression", zap.String("error:", err.Error()))
		return 0, err
	}
	res, err := strconv.ParseFloat(expressionNT.Token.Token, 64)
	if err != nil {
		return 0, err
	}
	return res, nil

}
