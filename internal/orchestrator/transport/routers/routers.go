package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/pkg/customError"
)

type Service interface {
	GetExpressions(size, page int) ([]*models.Expression, int64, *customError.CustomError)
	GetExpression(id int64) (*models.Expression, *customError.CustomError)
	CreateExpression(expression string) (int64, *customError.CustomError)
	SetResult(id, expressionID int64, result *float64, error *string) *customError.CustomError
	GetTask() (*models.Task, *customError.CustomError)
}
type Routers struct {
	Public *gin.RouterGroup
}

func CreateRouter(g *gin.Engine) *Routers {
	public := g.Group("/api/v1")
	return &Routers{
		Public: public,
	}
}
