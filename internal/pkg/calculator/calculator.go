package calculator

import (
	"context"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/orchestrator/repository"
	"github.com/xLeSHka/calc/internal/pkg/cache"
	"github.com/xLeSHka/calc/internal/pkg/counter"
	"github.com/xLeSHka/calc/internal/pkg/token"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"strconv"
)

type Calculator struct {
	Tasks   chan *models.Task
	Results chan *models.Task
	Log     *zap.Logger
	Repo    repository.Repo
	Cache   *cache.Cache
	Counter *counter.Counter
}

func New(
	log *zap.Logger,
	repo repository.Repo,
	times *cache.Cache,
	counter *counter.Counter,
) *Calculator {
	return &Calculator{
		Log:     log,
		Tasks:   make(chan *models.Task, 50),
		Results: make(chan *models.Task, 50),
		Repo:    repo,
		Cache:   times,
		Counter: counter,
	}
}
func (c *Calculator) SendTask(taskID, expressionID, taskTime int64, arg1, arg2 float64, operation models.Operation) {
	task := &models.Task{
		ID:            taskID,
		ExpressionID:  expressionID,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operation,
		OperationTime: taskTime,
	}
	c.Tasks <- task
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
				c.SendTask(taskId, expressionID, taskTime, a, b, models.Addition)
				for {
					t := <-c.Results
					if t.ID != taskId && t.ExpressionID != expressionID {
						c.Results <- t
						continue
					}
					node.Token.Token = strconv.FormatFloat(*t.Result, 'f', 5, 64)
					break
				}
			} else if node.Token.Token == "-" {
				taskId := c.Counter.Int()
				taskTime := c.Cache.SubtractionTime().Milliseconds()
				c.SendTask(taskId, expressionID, taskTime, a, b, models.Subtraction)
				for {
					t := <-c.Results
					if t.ID != taskId && t.ExpressionID != expressionID {
						c.Results <- t
						continue
					}
					node.Token.Token = strconv.FormatFloat(*t.Result, 'f', 5, 64)
					break
				}
			} else if node.Token.Token == "*" {
				taskId := c.Counter.Int()
				taskTime := c.Cache.MultiplicationTime().Milliseconds()
				c.SendTask(taskId, expressionID, taskTime, a, b, models.Multiplication)
				for {
					t := <-c.Results
					if t.ID != taskId && t.ExpressionID != expressionID {
						c.Results <- t
						continue
					}
					node.Token.Token = strconv.FormatFloat(*t.Result, 'f', 5, 64)
					break
				}
			} else if node.Token.Token == "/" {
				if b == 0 {
					ctx.Done()
					return ErrDivisionByZero
				}
				taskId := c.Counter.Int()
				taskTime := c.Cache.DivisionTime().Milliseconds()
				c.SendTask(taskId, expressionID, taskTime, a, b, models.Division)
				for {
					t := <-c.Results
					if t.ID != taskId && t.ExpressionID != expressionID {
						c.Results <- t
						continue
					}
					node.Token.Token = strconv.FormatFloat(*t.Result, 'f', 5, 64)
					break
				}
			} else if node.Token.Token == "^" {
				taskId := c.Counter.Int()
				taskTime := c.Cache.ExponentiationTime().Milliseconds()
				c.SendTask(taskId, expressionID, taskTime, a, b, models.Exponentiation)
				for {
					t := <-c.Results
					if t.ID != taskId && t.ExpressionID != expressionID {
						c.Results <- t
						continue
					}
					node.Token.Token = strconv.FormatFloat(*t.Result, 'f', 5, 64)
					break
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
				c.SendTask(taskId, expressionID, taskTime, a, 0.0, models.UnaryMinus)
				for {
					t := <-c.Results
					if t.ID != taskId && t.ExpressionID != expressionID {
						c.Results <- t
						continue
					}
					node.Token.Token = strconv.FormatFloat(*t.Result, 'f', 5, 64)
					break
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
			if a <= 0 || a == 1 {
				ctx.Done()
				return ErrLogNotDefinedFor
			}
			if b <= 0.0 {
				ctx.Done()
				return ErrLogOutOfFuncDomain
			}
			taskId := c.Counter.Int()
			taskTime := c.Cache.LogarithmTime().Milliseconds()
			c.SendTask(taskId, expressionID, taskTime, a, b, models.Logarithm)
			for {
				t := <-c.Results
				if t.ID != taskId && t.ExpressionID != expressionID {
					c.Results <- t
					continue
				}
				node.Token.Token = strconv.FormatFloat(*t.Result, 'f', 5, 64)
				break
			}
			//node.Token.Token = strconv.FormatFloat(, 'f', 5, 64)
		}
		if node.Token.Token == "sqrt" {
			a, err := strconv.ParseFloat(node.Left.Token.Token, 64)
			if err != nil {
				ctx.Done()
				return err
			}
			if a < 0.0 {
				ctx.Done()
				return ErrSqrtOutOfDomain
			}
			taskId := c.Counter.Int()
			taskTime := c.Cache.SquareRootTime().Milliseconds()
			c.SendTask(taskId, expressionID, taskTime, a, 0.0, models.SquareRoot)
			for {
				t := <-c.Results
				if t.ID != taskId && t.ExpressionID != expressionID {
					c.Results <- t
					continue
				}
				node.Token.Token = strconv.FormatFloat(*t.Result, 'f', 5, 64)
				break
			}
			//node.Token.Token = strconv.FormatFloat(math.Sqrt(a), 'f', 5, 64)
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
