package calculator

import (
	"context"
	"fmt"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"github.com/xLeSHka/calc/internal/pkg/cache"
	"github.com/xLeSHka/calc/internal/pkg/counter"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/pkg/token"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	Task   *models.Task
	Result chan float64
	Error  chan error
}
type Calculator struct {
	Tasks   map[int64]*Task
	TasksCh chan int64
	Log     *zap.Logger
	Repo    repository.Repo
	Cache   *cache.Cache
	Counter *counter.Counter
	mu      *sync.Mutex
}

func New(
	log *zap.Logger,
	repo repository.Repo,
	times *cache.Cache,
	counter *counter.Counter,
) *Calculator {
	return &Calculator{
		Log:     log,
		Tasks:   make(map[int64]*Task),
		TasksCh: make(chan int64, 200),
		mu:      &sync.Mutex{},
		Repo:    repo,
		Cache:   times,
		Counter: counter,
	}
}
func (c *Calculator) Exists(taskID int64) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.Tasks[taskID]
	return ok
}
func (c *Calculator) GetTask() (*models.Task, *customError.CustomError) {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case id, ok := <-c.TasksCh:
		if !ok {
			return nil, customError.New(http.StatusInternalServerError, fmt.Errorf("Calculator.GetTask: task channel closed"))
		}
		task, _ := c.Tasks[id]
		return task.Task, nil
	default:
		return nil, customError.New(http.StatusNotFound, fmt.Errorf("Calculator.GetTask: task not found"))
	}
}
func (c *Calculator) SendTask(taskID, expressionID, taskTime int64, arg1, arg2 float64, operation models.Operation) (chan float64, chan error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	task := &models.Task{
		ID:            taskID,
		ExpressionID:  expressionID,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operation,
		OperationTime: taskTime,
	}
	res := make(chan float64)
	err := make(chan error)
	t := &Task{
		Task:   task,
		Result: res,
		Error:  err,
	}
	c.Tasks[taskID] = t
	c.TasksCh <- taskID
	return res, err
}
func (c *Calculator) RecieveResult(task *models.Task) {
	c.mu.Lock()
	defer c.mu.Unlock()
	t, _ := c.Tasks[task.ID]
	defer close(t.Result)
	defer close(t.Error)
	if task.Result != nil {
		t.Result <- *task.Result

	} else {
		t.Error <- fmt.Errorf("Failed calculate, error: %s", *task.Error)
	}
}
func (c *Calculator) calculate(node *token.Node, expressionID int64, ctx context.Context) error {
	if node == nil {
		return nil
	}
	eg, ctx2 := errgroup.WithContext(context.Background())

	if node.Left != nil {
		eg.Go(func() error {
			return c.calculate(node.Left, expressionID, ctx2)
		})
	}
	if node.Right != nil {
		eg.Go(func() error {
			return c.calculate(node.Right, expressionID, ctx2)
		})
	}

	if node.Left == nil && node.Right == nil {
		return nil
	}
	if err := eg.Wait(); err != nil {
		ctx.Done()
		return err
	}
	switch node.Token.Type {
	case token.Operator:
		switch node.Token.Associativity {
		case token.Left:
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			b, err := strconv.ParseFloat(node.Right.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			if node.Token.Token == "+" {
				taskId := c.Counter.Int()
				taskTime := c.Cache.AddictionTime().Milliseconds()
				resCh, errCh := c.SendTask(taskId, expressionID, taskTime, a, b, models.Addition)

				select {
				case res := <-resCh:
					node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
				case err = <-errCh:
					ctx.Done()
					return err
				}
			} else if node.Token.Token == "-" {
				taskId := c.Counter.Int()
				taskTime := c.Cache.SubtractionTime().Milliseconds()
				resCh, errCh := c.SendTask(taskId, expressionID, taskTime, a, b, models.Subtraction)

				select {
				case res := <-resCh:
					node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
				case err = <-errCh:
					ctx.Done()
					return err
				}
			} else if node.Token.Token == "*" {
				taskId := c.Counter.Int()
				taskTime := c.Cache.MultiplicationTime().Milliseconds()
				resCh, errCh := c.SendTask(taskId, expressionID, taskTime, a, b, models.Multiplication)

				select {
				case res := <-resCh:
					node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
				case err = <-errCh:
					ctx.Done()
					return err
				}
			} else if node.Token.Token == "/" {
				taskId := c.Counter.Int()
				taskTime := c.Cache.DivisionTime().Milliseconds()
				resCh, errCh := c.SendTask(taskId, expressionID, taskTime, a, b, models.Division)

				select {
				case res := <-resCh:
					node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
				case err = <-errCh:
					ctx.Done()
					return err
				}
			} else if node.Token.Token == "^" {
				taskId := c.Counter.Int()
				taskTime := c.Cache.ExponentiationTime().Milliseconds()
				resCh, errCh := c.SendTask(taskId, expressionID, taskTime, a, b, models.Exponentiation)

				select {
				case res := <-resCh:
					node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
				case err = <-errCh:
					ctx.Done()
					return err
				}
			} else {
				ctx.Done()
				return token.ErrUnknownOperator

			}
		case token.Right:
			if node.Token.Token == "-" {
				a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
				if err != nil {
					ctx.Done()
					return err
				}
				taskId := c.Counter.Int()
				taskTime := c.Cache.UnaryMinusTime().Milliseconds()
				resCh, errCh := c.SendTask(taskId, expressionID, taskTime, a, 0.0, models.UnaryMinus)

				select {
				case res := <-resCh:
					node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
				case err = <-errCh:
					ctx.Done()
					return err
				}
			} else {
				ctx.Done()
				return token.ErrUnknownOperator
			}
		}
	case token.Function:
		if node.Token.Token == "log" {
			b, err := strconv.ParseFloat(node.Right.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			taskId := c.Counter.Int()
			taskTime := c.Cache.LogarithmTime().Milliseconds()
			resCh, errCh := c.SendTask(taskId, expressionID, taskTime, a, b, models.Logarithm)

			select {
			case res := <-resCh:
				node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
			case err = <-errCh:
				ctx.Done()
				return err
			}
		}
		if node.Token.Token == "sqrt" {
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			taskId := c.Counter.Int()
			taskTime := c.Cache.SquareRootTime().Milliseconds()
			resCh, errCh := c.SendTask(taskId, expressionID, taskTime, a, 0.0, models.SquareRoot)

			select {
			case res := <-resCh:
				node.Token.Token = strconv.FormatFloat(res, 'f', 5, 64)
			case err = <-errCh:
				ctx.Done()
				return err
			}
		}
	}
	return nil
}

func (c *Calculator) Calc(expressionNT token.Node, expressionID int64) {
	expr := &models.Expression{
		ID:     expressionID,
		Status: "In process",
	}
	err := c.Repo.UpdateExpression(expr)
	if err != nil {
		c.Log.Error("Failed set expression status to In process", zap.Error(err))
		return
	}
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return c.calculate(&expressionNT, expressionID, ctx)
	})
	if err = eg.Wait(); err != nil {
		c.Log.Error("Unprocessable expression", zap.Error(err))
		expr := &models.Expression{
			ID:     expressionID,
			Status: "Unprocessable expression",
		}
		err = c.Repo.UpdateExpression(expr)
		if err != nil {
			c.Log.Error("Failed set expression status to Unprocessable expression", zap.Error(err))
		}
		return
	}
	res, err := strconv.ParseFloat(expressionNT.Token.Token, 64)
	if err != nil {
		c.Log.Error("Unprocessable expression", zap.Error(err))
		expr := &models.Expression{
			ID:     expressionID,
			Status: "Unprocessable expression",
		}
		err = c.Repo.UpdateExpression(expr)
		if err != nil {
			c.Log.Error("Failed set expression status to Unprocessable expression", zap.Error(err))
		}
		return
	}
	expr = &models.Expression{
		ID:     expressionID,
		Status: "Solved",
		Result: &res,
	}
	err = c.Repo.UpdateExpression(expr)
	if err != nil {
		c.Log.Error("Failed set expression status to Solved", zap.Error(err))
	}
	return
}
