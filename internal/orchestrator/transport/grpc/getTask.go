package grpc

import (
	"context"
	proto "github.com/xLeSHka/calc/internal/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *Router) GetTask(ctx context.Context, req *emptypb.Empty) (*proto.GetTaskResponse, error) {
	task, cErr := r.agentService.GetTask()
	if cErr != nil {
		r.log.Error("GetTask: failed get task", zap.Error(cErr.Err))

		return nil, cErr.Err
	}
	return &proto.GetTaskResponse{
		Id:            task.ID,
		ExpressionId:  task.ExpressionID,
		Arg1:          float32(task.Arg1),
		Arg2:          float32(task.Arg2),
		Operation:     proto.Operation(task.Operation),
		OperationTime: task.OperationTime,
	}, nil
}
