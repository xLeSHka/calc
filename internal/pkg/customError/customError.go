package customError

import "github.com/gin-gonic/gin"

type CustomError struct {
	Err        error
	StatusCode int
}

func (e *CustomError) Error() string {
	return e.Err.Error()
}
func (e *CustomError) SendError(c *gin.Context) {
	c.JSON(e.StatusCode, gin.H{"status": "error", "message": e.Err.Error()})
}
func New(statusCode int, err error) *CustomError {
	return &CustomError{Err: err, StatusCode: statusCode}
}
