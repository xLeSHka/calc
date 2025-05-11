package agent

import (
	"context"
	"fmt"
	"github.com/xLeSHka/calc/internal/models"
	proto "github.com/xLeSHka/calc/internal/pkg/api"
	config2 "github.com/xLeSHka/calc/internal/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"math"
	"sync"
	"time"
)

type Agent struct {
	ComputingPower int
	Jobs           chan models.Task
	Results        chan models.Task
	Wg             *sync.WaitGroup
	Log            *zap.Logger
	Shutdown       chan struct{}
	Client         proto.CalcServiceClient
	Conn           *grpc.ClientConn
}
type PostResult struct {
	ID           int64   `json:"id"`
	ExpressionID int64   `json:"expression_id"`
	Result       float64 `json:"result,omitempty"`
	Error        *string `json:"error,omitempty"`
}

func (a *Agent) Recieve() {

	for {
		select {
		case <-a.Shutdown:
			return
		default:
			t, err := a.Client.GetTask(context.TODO(), &emptypb.Empty{})
			if err != nil {
				a.Log.Error("Agent Request failed", zap.Error(err))
				time.Sleep(1 * time.Second)
				continue
			}
			task := models.Task{
				ID:            t.Id,
				ExpressionID:  t.ExpressionId,
				Arg1:          float64(t.Arg1),
				Arg2:          float64(t.Arg2),
				Operation:     models.Operation(t.Operation),
				OperationTime: t.OperationTime,
			}
			a.Log.Info("Agent task received", zap.Any("task", task))
			a.Jobs <- task
		}
	}
}
func (a *Agent) Send() {
	for {
		select {
		case <-a.Shutdown:
			for {
				select {
				case task, ok := <-a.Results:
					if !ok {
						return
					}
					sended := false
					for i := 0; i < 3; i++ {
						res := float32(task.Result)
						_, err := a.Client.PostResult(context.TODO(), &proto.PostResultRequest{
							Id:           task.ID,
							ExpressionId: task.ExpressionID,
							Result:       &res,
							Error:        task.Error,
						})
						if err != nil {
							a.Log.Error("Agent Request failed", zap.Error(err))
							time.Sleep(1 * time.Second)
							continue
						}
						sended = true
						break
					}
					if sended {
						a.Log.Info("Agent result send", zap.Any("task", task))
					} else {
						a.Log.Info("Agent failed send result", zap.Any("task", task))
					}
				default:
					return
				}
			}
		case task, ok := <-a.Results:
			if !ok {
				return
			}
			sended := false
			for i := 0; i < 3; i++ {
				res := float32(task.Result)
				_, err := a.Client.PostResult(context.TODO(), &proto.PostResultRequest{
					Id:           task.ID,
					ExpressionId: task.ExpressionID,
					Result:       &res,
					Error:        task.Error,
				})
				if err != nil {
					a.Log.Error("Agent Request failed", zap.Error(err))
					time.Sleep(1 * time.Second)
					continue
				}
				sended = true
				break
			}
			if sended {
				a.Log.Info("Agent result send", zap.Any("task", task))
			} else {
				a.Log.Info("Agent failed send result", zap.Any("task", task))
			}
		}
	}

}
func (a *Agent) Start(ctx context.Context) error {
	a.Log.Info("Starting agent")
	a.Wg.Add(a.ComputingPower)
	for i := 0; i < a.ComputingPower; i++ {
		go a.Worker()
	}
	go a.Recieve()
	go a.Send()
	return nil
}
func (a *Agent) Stop(ctx context.Context) error {
	close(a.Shutdown)
	close(a.Jobs)
	defer close(a.Results)
	a.Wg.Wait()
	a.Conn.Close()
	return nil
}
func (a *Agent) Worker() {
	defer a.Wg.Done()
	for j := range a.Jobs {
		dur, err := time.ParseDuration(fmt.Sprintf("%dms", j.OperationTime))
		if err != nil {
			continue
		}
		time.Sleep(dur)
		switch j.Operation {
		case models.Addition:
			res := j.Arg1 + j.Arg2
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.Subtraction:
			res := j.Arg1 - j.Arg2
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.Multiplication:
			res := j.Arg1 * j.Arg2
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.Division:
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
			}
			if j.Arg2 == 0 {
				errMsg := ErrDivisionByZero.Error()
				t.Error = &errMsg
				t.Result = 0.0
			} else {
				res := j.Arg1 / j.Arg2
				t.Result = res
				t.Error = nil
			}
			a.Results <- t
		case models.Exponentiation:
			res := math.Pow(j.Arg1, j.Arg2)
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.UnaryMinus:
			res := -j.Arg1
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
				Result:       res,
				Error:        nil,
			}
			a.Results <- t
		case models.Logarithm:
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
			}
			if j.Arg1 <= 0 || j.Arg1 == 1 {
				errMsg := ErrLogNotDefinedFor.Error()
				t.Error = &errMsg
				t.Result = .0
			} else if j.Arg2 <= 0.0 {
				errMsg := ErrLogOutOfFuncDomain.Error()
				t.Error = &errMsg
				t.Result = 0.0
			} else {
				res := math.Log(j.Arg2) / math.Log(j.Arg1)
				t.Result = res
				t.Error = nil
			}
			a.Results <- t
		case models.SquareRoot:
			t := models.Task{
				ID:           j.ID,
				ExpressionID: j.ExpressionID,
			}
			if j.Arg1 < 0 {
				errMsg := ErrSqrtOutOfDomain.Error()
				t.Error = &errMsg
				t.Result = 0.0
			} else {
				res := math.Sqrt(j.Arg1)
				t.Result = res
				t.Error = nil
			}
			a.Results <- t
		}
	}
}
func New(config config2.Config, lc fx.Lifecycle, log *zap.Logger) *Agent {
	agent := &Agent{
		Jobs:           make(chan models.Task, 100),
		Results:        make(chan models.Task, 100),
		ComputingPower: config.ComputingPower,
		Wg:             &sync.WaitGroup{},
		Log:            log,
		Shutdown:       make(chan struct{}),
	}
	addr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort+1)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	agent.Conn = conn
	client := proto.NewCalcServiceClient(conn)
	agent.Client = client
	log.Info("Agent created")
	lc.Append(fx.Hook{
		OnStart: agent.Start,
		OnStop:  agent.Stop,
	})
	log.Info("Agent started")
	return agent
}
