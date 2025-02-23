package models

type Operation byte

const (
	Addition Operation = iota
	Subtraction
	Division
	Multiplication
	Exponentiation
	UnaryMinus
	Logarithm
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
	ID            int64     `json:"id"`
	ExpressionID  int64     `json:"expression_id"`
	Arg1          float64   `json:"arg1"`
	Arg2          float64   `json:"arg2"`
	Operation     Operation `json:"operation"`
	OperationTime int64     `json:"operation_time"`
	Result        float64   `json:"result,omitempty"`
	Error         *string   `json:"error,omitempty"`
}
