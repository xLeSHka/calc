package agentService

import (
	"fmt"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"net/http"
)

func (s *AgentService) SetResult(id, expressionID int64, result float64, error *string) *customError.CustomError {
	task := &models.Task{
		ID:           id,
		ExpressionID: expressionID,
		Result:       result,
		Error:        error,
	}
	if !s.Calculator.Exists(id) {
		return customError.New(http.StatusNotFound, fmt.Errorf("Service.SerResult: task not found"))
	}
	s.Calculator.ReceiveResult(task)
	return nil
}
