package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/utils/jwt"
	"go.uber.org/zap"
	"net/http"
)

func (r *Router) GetExpressions(c *gin.Context) {
	userID, err := jwt.Parse(c)
	if err != nil {
		err.(*customError.CustomError).SendError(c)
		c.Abort()
		return
	}
	var req GetExpressionsReq
	if err := r.validator.ShouldBindQuery(c, &req); err != nil {
		r.Log.Error("GetExpressions: Failed bind query", zap.Error(err))
		customError.New(http.StatusBadRequest, fmt.Errorf("GetExpressions: error: %w", err)).SendError(c)
		c.Abort()
		return
	}
	if req.Size == 0 {
		req.Size = 10
	}
	exprs, total, cErr := r.usersService.GetExpressions(req.Size, req.Page, userID)
	if cErr != nil {
		r.Log.Error("GetExpressions: Failed get expressions", zap.Error(cErr))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	resp := make([]*GetExpressionResp, len(exprs))
	for i, expr := range exprs {
		resp[i] = mapExpression(expr)
	}
	c.JSON(http.StatusOK, resp)
}
