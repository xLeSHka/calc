package models

import "github.com/google/uuid"

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

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" `
	Login    string    `gorm:"unique;not null"`
	Password []byte    `gorm:"not null"`
}
type Expression struct {
	ID         int64     `gorm:"type:serial;primary_key"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	Expression string    `gorm:"not null"`
	Status     string    `gorm:"not null"`
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
