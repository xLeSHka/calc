package users

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
	ID         int64    `json:"id"`
	Expression string   `json:"expression"`
	Status     string   `json:"status"`
	Result     *float64 `json:"result,omitempty"`
}
type GetExpressionsReq struct {
	Size int `form:"size" binding:"omitempty,gte=0"`
	Page int `form:"page" binding:"omitempty,gte=0"`
}

func mapExpression(expression *models.Expression) *GetExpressionResp {
	return &GetExpressionResp{
		ID:         expression.ID,
		Expression: expression.Expression,
		Status:     expression.Status,
		Result:     expression.Result,
	}
}

//type GetTaskResp struct {
//
//}
//
//type PostResult struct {
//
//}
