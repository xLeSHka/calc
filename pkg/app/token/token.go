package token

import (
	"strings"
	"unicode"
)

// делаем enum для удобного определения типа токена
type Type int

const (
	Operator Type = iota
	LParanthesis
	RParanthesis
	IntLiteral
	FloatLiteral
	Function
	Separator
)

// делаем enum для удобного определения ассщциативности токена
type Associativity int

const (
	None Associativity = iota
	Right
	Left
)

// создаем структуру для удобного представления токенов
type Token struct {
	Type
	Associativity
	Token string
}

// создаем структуру дерева чтобы удобно ходить по ней
type Node struct {
	Left  *Node
	Right *Node
	Token *Token
}

// делаем enum для удобного отслеживания состояния парсера выражений
type State int

const (
	S0 State = iota
	S1
	S2
	S3
	S4
	S5
)

// сконструктор токенов
func NewToken(token string, tokenType Type, asc Associativity) (*Token, error) {
	if tokenType == Operator && asc == None {
		return nil, ErrAssociativityReq
	}
	if tokenType != Operator && asc != None {
		return nil, ErrNonOperationTokenAsc
	}
	return &Token{
		tokenType,
		asc,
		token,
	}, nil
}

// получаем приоритет каждого оператора
func (t *Token) GetPrecedence() (int, error) {
	leftAssociativity := map[string]int{
		"+": 2,
		"-": 2,
		"/": 3,
		"*": 3,
		"^": 5,
	}
	rightAssociativity := map[string]int{
		"-": 4,
	}

	switch t.Associativity {
	case Left:
		if _, ok := leftAssociativity[t.Token]; ok {
			return leftAssociativity[t.Token], nil
		} else {
			return 0, ErrUnknownOperator
		}
	case Right:
		if _, ok := rightAssociativity[t.Token]; ok {
			return rightAssociativity[t.Token], nil
		} else {
			return 0, ErrUnknownOperator
		}
	case None:
		return 0, ErrTokenIsNotAnOperator
	default:
		return 0, ErrTokenIsNotAnOperator
	}
}

func TokenizeExpression(expression string) (*Node, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return nil, err
	}
	return toNT(tokens)
}

// токенизируем выражение
func tokenize(expr string) ([]*Token, error) {
	state := S0
	validOperators := "+-*^/"
	tokens := make([]*Token, 0)

	var buf string
	var bufTokenType Type
	for _, t := range expr {
		isDigit := unicode.IsDigit(t)
		isLetter := unicode.Is(unicode.Latin, t)
		isLParanth := t == '('
		isRParanth := t == ')'
		isParanth := isLParanth || isRParanth
		isPoint := t == '.'
		isSep := t == ','
		isOp := strings.ContainsAny(string(t), validOperators)

		if !(isDigit || isLetter || isParanth || isPoint || isSep || isOp) {
			return nil, ErrUnknownSymbol
		}

		switch state {
		case S0:
			if isOp || isParanth {
				state = S1
			} else if isDigit {
				state = S2
			} else if isLetter {
				state = S4
			} else if isPoint || isSep {
				return nil, ErrUnexpectedSymbol
			}
		case S1:
			if isDigit {
				state = S2
			} else if isLetter {
				state = S4
			} else if isPoint {
				return nil, ErrUnexpectedSymbol
			}
		case S2:
			bufTokenType = IntLiteral
			if isPoint {
				state = S3
			} else if isParanth || isOp || isSep || isPoint {
				state = S5
			} else if isLetter {
				return nil, ErrUnexpectedSymbol
			}
		case S3:
			bufTokenType = FloatLiteral
			if isParanth || isOp || isSep || isPoint {
				state = S5
			} else if isLetter {
				return nil, ErrUnexpectedSymbol
			}
		case S4:
			bufTokenType = Function
			if isLParanth {
				state = S5
			} else if isRParanth || isSep || isOp || isPoint {
				return nil, ErrUnexpectedSymbol
			}
		case S5:
			if isOp || isParanth || isSep {
				state = S1
			} else if isDigit {
				state = S2
			} else if isLetter {
				state = S4
			} else if isPoint {
				return nil, ErrUnexpectedSymbol
			}
		default:
			break
		}
		fn := func() error {
			if isOp {
				//обработка unary negation
				if len(tokens) == 0 || tokens[len(tokens)-1].Type == LParanthesis {
					tkn, err := NewToken(string(t), Operator, Right)
					if err != nil {
						return err
					}
					tokens = append(tokens, tkn)
				} else {
					tkn, err := NewToken(string(t), Operator, Left)
					if err != nil {
						return err
					}
					tokens = append(tokens, tkn)
				}
			} else if isParanth {
				if isRParanth {
					tkn, err := NewToken(string(t), RParanthesis, None)
					if err != nil {
						return err
					}
					tokens = append(tokens, tkn)
				} else {
					tkn, err := NewToken(string(t), LParanthesis, None)
					if err != nil {
						return err
					}
					tokens = append(tokens, tkn)
				}
			} else if isSep {
				tkn, err := NewToken(string(t), Separator, None)
				if err != nil {
					return err
				}
				tokens = append(tokens, tkn)
			}
			return nil
		}

		switch state {
		case S1:
			err := fn()
			if err != nil {
				return nil, err
			}
		case S2, S3, S4:
			buf += string(t)
		case S5:
			tkn, err := NewToken(buf, bufTokenType, None)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tkn)
			buf = ""
			err = fn()
			if err != nil {
				return nil, err
			}
		}
	}
	if buf != "" {
		tkn, err := NewToken(buf, bufTokenType, None)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tkn)
	}
	return toRPN(tokens)
}

// перевод выражения в обратную польскую нотацию для удобного чтения
func toRPN(tokens []*Token) ([]*Token, error) {

	var queue []*Token
	var stack []*Token

	for _, tkn := range tokens {
		switch tkn.Type {
		case IntLiteral, FloatLiteral:
			queue = append(queue, tkn)
		case LParanthesis, Function:
			stack = append(stack, tkn)
		case Operator:
			for len(stack) > 0 && (stack[len(stack)-1].Type == Operator) {
				stackPrec, err := stack[len(stack)-1].GetPrecedence()
				if err != nil {
					return nil, err
				}
				tknPrec, err := tkn.GetPrecedence()
				if err != nil {
					return nil, err
				}
				if (stackPrec > tknPrec) || (stackPrec == tknPrec && tkn.Associativity == Left) {
					stackToQueue(&queue, &stack)
				} else {
					break
				}
			}
			stack = append(stack, tkn)
		case RParanthesis:
			if len(stack) == 0 {
				return nil, ErrMissLeftParanthesis
			}
			for stack[len(stack)-1].Type != LParanthesis {
				stackToQueue(&queue, &stack)
				if len(stack) == 0 {
					return nil, ErrMissLeftParanthesis
				}
			}
			stack = stack[:len(stack)-1]
			if len(stack) > 0 && stack[len(stack)-1].Type == Function {
				stackToQueue(&queue, &stack)
			}
		case Separator:
			if len(stack) == 0 {
				return nil, ErrMisedSepOrParanth
			}
			for stack[len(stack)-1].Type != LParanthesis {
				stackToQueue(&queue, &stack)
				if len(stack) == 0 {
					return nil, ErrMisedSepOrParanth
				}
			}
		}
	}
	for len(stack) > 0 {
		if stack[len(stack)-1].Type == LParanthesis {
			return nil, ErrMissRightParanthesis
		}
		stackToQueue(&queue, &stack)
	}
	return queue, nil
}

func stackToQueue(queue *[]*Token, stack *[]*Token) {
	token := (*stack)[len(*stack)-1]
	*queue = append(*queue, token)
	*stack = (*stack)[:len(*stack)-1]
}

// Переводим выражение в дерево нод для более удобного взаимодействия, тк так лечше будет сделать конкурентную обработку выражения
func toNT(expr []*Token) (*Node, error) {
	var queue []*Node
	for _, t := range expr {
		switch t.Type {
		case IntLiteral, FloatLiteral:
			queue = append(queue, &Node{Token: t})
		case Operator:
			switch t.Associativity {
			case Left:
				if len(queue) < 2 {
					return nil, ErrConvertRPNToNT
				}
				stackTwoToQueue(&queue, t)
			case Right:
				if len(queue) < 1 {
					return nil, ErrConvertRPNToNT
				}
				stackOneToQueue(&queue, t)
			}
		case Function:
			if t.Token == "log" {
				if len(queue) < 2 {
					return nil, ErrConvertRPNToNT
				}
				stackTwoToQueue(&queue, t)
			} else if t.Token == "sqrt" {
				if len(queue) < 1 {
					return nil, ErrConvertRPNToNT
				}
				stackOneToQueue(&queue, t)
			}
		}
	}
	if len(queue) != 1 {
		return nil, ErrConvertRPNToNT
	}
	return queue[0], nil
}

// перемещаем две ноды из стека в очередь, причем в левую ноду помещаем предпоследний элемент,
// когда как в правую последний
func stackTwoToQueue(queue *[]*Node, token *Token) {
	right := (*queue)[len(*queue)-1]
	*queue = (*queue)[:len(*queue)-1]
	left := (*queue)[len(*queue)-1]
	*queue = (*queue)[:len(*queue)-1]
	node := &Node{Right: right, Left: left, Token: token}
	*queue = append(*queue, node)
}

// перемещаем ноду из стека в очередь, причем на левую позицию
func stackOneToQueue(queue *[]*Node, token *Token) {
	left := (*queue)[len(*queue)-1]
	*queue = (*queue)[:len(*queue)-1]
	node := &Node{Right: nil, Left: left, Token: token}
	*queue = append(*queue, node)
}
