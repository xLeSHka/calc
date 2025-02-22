package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (r *Router) GetExpressions(c *gin.Context) {

	exprs, cErr := r.service.GetExpressions()
	if cErr != nil {
		r.Log.Error("GetExpressions: Failed get expressions", zap.Error(cErr))
		cErr.SendError(c)
		c.Abort()
		return
	}
	resp := make([]*GetExpressionResp, len(exprs))
	for i, expr := range exprs {
		resp[i] = mapExpression(expr)
	}
	c.JSON(http.StatusOK, resp)
}
