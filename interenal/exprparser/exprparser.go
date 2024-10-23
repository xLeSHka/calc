package exprparser

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	digits string = "0123456789."
)

func getPrecedance(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}
func Tokenize(expr string) []string {
	var tokens []string
	var nums string
	for _, r := range expr {
		strRune := string(r)
		if strings.ContainsAny(strRune, digits) {
			nums += strRune
		} else {
			if nums != "" {
				tokens = append(tokens, nums)
				nums = ""
			}
			if strRune != " " {
				tokens = append(tokens, strRune)
			}
		}
	}
	if nums != "" {
		tokens = append(tokens, nums)
	}
	return tokens
}

func ToPostfixNotation(tokens []string) ([]string, error) {
	var output []string
	var operators []string
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			for len(operators) > 0 && getPrecedance(operators[len(operators)-1]) >= getPrecedance(token) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)

		case "(":
			operators = append(operators, token)
		case ")":
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, fmt.Errorf("miss left paranthes")
			}
			operators = operators[:len(operators)-1]
		default:
			output = append(output, token)
		}
	}
	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, fmt.Errorf("miss left paranthes")
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}
	return output, nil
}

func Calculate(RPNTokens []string) (float64, error) {
	var stack []float64
	for _, token := range RPNTokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, fmt.Errorf("invalid expression")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, fmt.Errorf("division by zero")
				}
				stack = append(stack, a/b)
			default:
				return 0, fmt.Errorf("invalid operator %s", token)
			}
		}
	}
	if len(stack) != 1 {
		return 0, fmt.Errorf("invalid expression")
	}
	return stack[0], nil
}
