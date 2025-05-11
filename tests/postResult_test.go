package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xLeSHka/calc/internal/models"
	proto "github.com/xLeSHka/calc/internal/pkg/api"
	"github.com/xLeSHka/calc/internal/utils/password"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestPostTask(t *testing.T) {
	fn, quit, err := setUp()
	assert.Nil(t, err)
	defer func() {
		quit <- syscall.SIGINT
		fn()
	}()
	addr := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort+4)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client = proto.NewCalcServiceClient(conn)
	type TestPostTask struct {
		ToPost       *PostTask
		name         string
		expectedCode int
	}
	tests := []TestPostTask{
		TestPostTask{
			name:         "Succes post task ",
			ToPost:       &PostTask{ID: 1, ExpressionID: 1, Result: 4.0},
			expectedCode: http.StatusOK,
		},
	}
	for i := 0; i < 10; i++ {
		solveExpr()
		time.Sleep(100 * time.Millisecond)
	}
	hash, _ := password.Encrypt([]byte("passworD!2"), []byte(config.CryptoKey))
	UsersService.Register(&models.User{
		ID:       profile1.ID,
		Login:    profile1.Login,
		Password: hash,
	})
	UsersService.CreateExpression("2+2", profile1.ID)
	time.Sleep(100 * time.Millisecond)
	task, _ := AgentService.GetTask()
	tests[0].ToPost.ID = task.ID
	tests[0].ToPost.ExpressionID = task.ExpressionID
	tests[0].ToPost.Result = task.Result
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var res float32
			res = float32(tests[0].ToPost.Result)
			_, err := client.PostResult(context.TODO(), &proto.PostResultRequest{
				Id:           tests[0].ToPost.ID,
				ExpressionId: tests[0].ToPost.ExpressionID,
				Result:       &res,
			})
			assert.Nil(t, err)
			time.Sleep(100 * time.Millisecond)
			var expression models.Expression
			result := db.Model(&models.Expression{}).Where("id = ?", test.ToPost.ExpressionID).First(&expression)
			assert.Nil(t, result.Error)
			assert.Equal(t, test.ToPost.Result, *expression.Result)
		})
	}
}
