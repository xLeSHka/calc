package agentService

import (
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/pkg/customError"
)

func (s *AgentService) GetTask() (*models.Task, *customError.CustomError) {
	task, cErr := s.Calculator.GetTask()
	if cErr != nil {
		return nil, cErr
	}
	return task, nil
}
