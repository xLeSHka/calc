package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/models"
	"go.uber.org/zap"
)

func (r *Router) GetTask(c *gin.Context) {
	task, cErr := r.service.GetTask()
	if cErr != nil {
		r.Log.Error("GetTask: failed get task", zap.Error(cErr.Err))
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.JSON(200, mapTask(task))
}
func mapTask(task *models.Task) *GetTaskResp {
	return &GetTaskResp{
		ID:            task.ID,
		ExpressionID:  task.ExpressionID,
		Arg1:          task.Arg1,
		Arg2:          task.Arg2,
		Operation:     task.Operation,
		OperationTime: task.OperationTime,
	}
}
