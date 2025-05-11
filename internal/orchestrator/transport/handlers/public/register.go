package public

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xLeSHka/calc/internal/models"
	"github.com/xLeSHka/calc/internal/pkg/customError"
	"github.com/xLeSHka/calc/internal/utils/password"
	"net/http"
)

func (r *Router) Register(c *gin.Context) {
	var reqData RegisterReq
	if err := r.validator.ShouldBindJSON(c, &reqData); err != nil {
		customError.New(http.StatusBadRequest, err).SendError(c)
		c.Abort()
		return
	}
	encrypted, err := password.Encrypt([]byte(*reqData.Password), r.cryptoKey)
	if err != nil {
		customError.New(http.StatusInternalServerError, err).SendError(c)
		c.Abort()
		return
	}
	user := &models.User{
		ID:       uuid.New(),
		Login:    *reqData.Login,
		Password: encrypted,
	}
	token, cErr := r.usersService.Register(user)
	if cErr != nil {
		cErr.SendError(c)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
