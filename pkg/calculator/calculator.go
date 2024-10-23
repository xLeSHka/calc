package pkg

import "github.com/xLeSHka/calculator/interenal/exprparser"

func Calc(expression string) (float64, error) {
	tokens := exprparser.Tokenize(expression)

	postfixed, err := exprparser.ToPostfixNotation(tokens)
	if err != nil {
		return 0, err
	}
	calculated, err := exprparser.Calculate(postfixed)
	if err != nil {
		return 0, err
	}
	return calculated, nil
}
