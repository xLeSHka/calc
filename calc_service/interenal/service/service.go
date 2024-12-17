package service

import (
	"context"

	"github.com/xLeSHka/calculator/pkg/calculator"
)

type CalcService struct{}

func New() *CalcService {
	return &CalcService{}
}

func (s *CalcService) Calculate(ctx context.Context, expression string) (float64, error) {
	return calculator.Calc(expression)
}
