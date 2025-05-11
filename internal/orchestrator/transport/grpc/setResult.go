package grpc

import (
	"context"
	proto "github.com/xLeSHka/calc/internal/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *Router) PostResult(ctx context.Context, req *proto.PostResultRequest) (*emptypb.Empty, error) {
	cErr := r.agentService.SetResult(req.GetId(), req.GetExpressionId(), float64(req.GetResult()), req.Error)
	if cErr != nil {
		r.log.Error("PostResult: SetResult", zap.Error(cErr))

		return nil, cErr
	}
	return &emptypb.Empty{}, nil
}
