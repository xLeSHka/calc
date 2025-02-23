package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"go.uber.org/zap"
	"net/http"
)

func (r *Router) PostResult(c *gin.Context) {
	var req PostResult
	if err := c.ShouldBindJSON(&req); err != nil {
		r.Log.Error("PostResult: Failed bind body", zap.Error(err))
		customError.New(http.StatusBadRequest, fmt.Errorf("PostResult: error: %w", err)).SendError(c)
		c.Abort()
		return
	}
	cErr := r.service.SetResult(*req.ID, *req.ExpressionID, req.Result, req.Error)
	if cErr != nil {
		r.Log.Error("PostResult: SetResult", zap.Error(cErr))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
