package calculator

import (
	"strconv"
	"strings"
)

const (
	digits string = "0123456789."
)

// Функция определения приоритета токена
func getPrecedence(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

// Функция превращения выражения в массив токенов
func tokenize(expr string) []string {
	var tokens []string
	var nums string
	for _, r := range expr {

		// Если руна - цифра, то идем по выражению до следующего символа,
		// не являющегося цифрой
		strRune := string(r)
		if strings.ContainsAny(strRune, digits) {
			nums += strRune
		} else {
			// Записываем все число как один токен и обнуляем запись числа
			if nums != "" {
				tokens = append(tokens, nums)
				nums = ""
			}
			//Пропускаем пробелы и записываем каждый символ как отдельный токен
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

// Функция конвертации выражения в postfix notation, которая прывычнее компьютеру
// Пример (A+B)-C ==> (AB+)C-
func toPostfixNotation(tokens []string) ([]string, error) {
	var output []string
	var operators []string
	for _, token := range tokens {
		switch token {
		// проходим по токенам, если токен - оператор, то проходим по очереди операторов
		// пока приоритет текущего оператора не будет больше приоритета предыдущего
		case "+", "-", "*", "/":
			for len(operators) > 0 && getPrecedence(operators[len(operators)-1]) >= getPrecedence(token) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)

		case "(":
			operators = append(operators, token)
		// Если встретили закрывающую скобку, то проходим по очереди операторов, пока не встретим открывающую скобку,
		// если такой нет - возвращаем ошибку
		case ")":
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, ErrMissLeftParanthesis
			}
			operators = operators[:len(operators)-1]
		default:
			output = append(output, token)
		}
	}
	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, ErrMissRightParanthesis
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}
	return output, nil
}

// Функция вычисления postfix notation выражения
func calculate(RPNTokens []string) (float64, error) {
	var stack []float64
	for _, token := range RPNTokens {
		// пробуем парсить токен как число
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			// иначе парсим как оператор и проверяем его валидность
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
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
					return 0, ErrDivisionByZero
				}
				stack = append(stack, a/b)
			default:
				return 0, ErrInvalidExpression2
			}
		}
	}
	if len(stack) != 1 {
		return 0, ErrInvalidExpression1
	}
	return stack[0], nil
}
func Calc(expression string) (float64, error) {
	if expression == "" {
		return 0, ErrInvalidExpression
	}
	tokens := tokenize(expression)
	postfixed, err := toPostfixNotation(tokens)
	if err != nil {
		return 0, err
	}
	calculated, err := calculate(postfixed)
	if err != nil {
		return 0, err
	}
	return calculated, nil
}
