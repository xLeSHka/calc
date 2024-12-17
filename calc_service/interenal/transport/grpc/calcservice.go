package grpc

import (
	"context"

	"github.com/xLeSHka/calculator/pkg/api/calc"
)

type Service interface {
	Calculate(ctx context.Context, expression string) (float64, error)
}

type CalcService struct {
	calc.UnimplementedCalcServiceServer
	service Service
}

func NewCalcService(srv Service) *CalcService {
	return &CalcService{service: srv}
}
//Импелентируем grpc сервер
func (cs *CalcService) Calculate(ctx context.Context, req *calc.CalculateRequest) (*calc.CalculateResponse, error) {
	res, err := cs.service.Calculate(ctx, req.GetExpression())
	if err != nil {
		return nil, err
	}
	return &calc.CalculateResponse{Result: float32(res)}, nil
}
