package models

type Operation byte

const (
	Addition Operation = iota
	Subtraction
	Division
	Multiplication
	Exponentiation
	UnaryMinus
	Logariphm
	SquareRoot
)

type Expression struct {
	ID         int64  `gorm:"type:serial;primary_key"`
	Expression string `gorm:"not null"`
	Status     string `gorm:"not null"`
	Result     *float64
}

func (_ Expression) TableName() string {
	return "expressions"
}

type Task struct {
	ID            int64
	ExpressionID  int64
	Arg1          float64
	Arg2          float64
	Operation     Operation
	OperationTime int64
	Result        *float64
}
