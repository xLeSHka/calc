package public

import (
	"github.com/gin-gonic/gin"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"net/http"
)

func (r *Router) Login(c *gin.Context) {
	var reqData RegisterReq
	if err := r.validator.ShouldBindJSON(c, &reqData); err != nil {
		customError.New(http.StatusBadRequest, err).SendError(c)
		c.Abort()
		return
	}
	token, cErr := r.usersService.Login(*reqData.Login, *reqData.Password)
	if cErr != nil {
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
