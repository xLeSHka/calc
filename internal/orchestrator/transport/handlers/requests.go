package handlers

import "github.com/xLeSHka/calc/internal/models"

type CreateExpressionReq struct {
	Expression *string `json:"expression" binding:"required"`
}
type CreateExpressionResp struct {
	ID int64 `json:"id"`
}
type GetExpressionReq struct {
	ID *int64 `uri:"id" binding:"required,gte=1"`
}

type GetExpressionResp struct {
	ID     int64    `json:"id"`
	Status string   `json:"status"`
	Result *float64 `json:"result,omitempty"`
}
type GetExpressionsReq struct {
	Size int `form:"size" binding:"omitempty,gte=0"`
	Page int `form:"page" binding:"omitempty,gte=0"`
}
type GetTaskResp struct {
	ID            int64            `json:"id"`
	ExpressionID  int64            `json:"expression_id"`
	Arg1          float64          `json:"arg1"`
	Arg2          float64          `json:"arg2"`
	Operation     models.Operation `json:"operation"`
	OperationTime int64            `json:"operation_time"`
}

type PostResult struct {
	ID           *int64  `json:"id" binding:"required,gte=1"`
	ExpressionID *int64  `json:"expression_id" binding:"required"`
	Result       float64 `json:"result" binding:"omitempty"`
	Error        *string `json:"error" binding:"omitempty"`
}
