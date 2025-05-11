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
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestGetTask(t *testing.T) {
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
	type TestGetTask struct {
		Expected     *models.Task
		name         string
		expectedCode int
	}
	tests := []TestGetTask{
		TestGetTask{
			name: "Succes get task ",
			Expected: &models.Task{
				Arg1:          2.0,
				Arg2:          2.0,
				Operation:     models.Addition,
				OperationTime: config.TimeAddiction.Milliseconds(),
			},
			expectedCode: http.StatusOK,
		},
		TestGetTask{
			name:         "Task not found ",
			expectedCode: http.StatusNotFound,
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			task, err := client.GetTask(context.TODO(), &emptypb.Empty{})
			if test.expectedCode != http.StatusOK {
				assert.NotNil(t, err)
			}
			if test.expectedCode == http.StatusOK {
				fmt.Print(t, err)
				assert.Nil(t, err)
				assert.Equal(t, test.Expected.Arg1, task.Arg1)
				assert.Equal(t, test.Expected.Arg2, task.Arg2)
				assert.Equal(t, test.Expected.Operation, task.Operation)
				assert.Equal(t, test.Expected.OperationTime, task.OperationTime)
			}
		})
	}
}
