package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"go.uber.org/zap"
	"net/http"
)

func (r *Router) GetExpression(c *gin.Context) {
	var req GetExpressionReq
	if err := c.ShouldBindUri(&req); err != nil {
		r.Log.Error("GetExpression: Failed bind uri", zap.Error(err))
		customError.New(http.StatusBadRequest, fmt.Errorf("GetExpression: error: %w", err)).SendError(c)
		c.Abort()
		return
	}
	expr, cErr := r.service.GetExpression(*req.ID)
	if cErr != nil {
		r.Log.Error("GetExpression: Failed get expression", zap.Error(cErr))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, mapExpression(expr))
}
func mapExpression(expression *models.Expression) *GetExpressionResp {
	return &GetExpressionResp{
		ID:     expression.ID,
		Status: expression.Status,
		Result: expression.Result,
	}
}
