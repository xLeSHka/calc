package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/xLeSHka/calculator/pkg/api/calc"
	"github.com/xLeSHka/calculator/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	restServer *http.Server
	listener   net.Listener
}

// создаем новый grpc сервер с gateway`ем
func New(ctx context.Context, port, restPort int, service Service) (*Server, error) {
	//настраиваем листенер
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	//добавляем интерсептор логирования
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
		),
	}
	//создаем новый grpc сервер с интерсептором логирования
	grpcServer := grpc.NewServer(opts...)
	calc.RegisterCalcServiceServer(grpcServer, NewCalcService(service))

	//создаем gateway для grpc сервера
	restSrv := runtime.NewServeMux()

	if err := calc.RegisterCalcServiceHandlerServer(ctx, restSrv, NewCalcService(service)); err != nil {
		return nil, err
	}
	//добавляем мидлвар логирования для gateway
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: LoggingMiddleware(restSrv, ctx),
	}
	return &Server{grpcServer, httpServer, lis}, nil
}

// Функция запуска приложения
// тут будем читать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
// func (a *Application) Run() error {
// 	for {
// 		// читаем выражение для вычисления из командной строки
// 		log.Println("input expression")
// 		reader := bufio.NewReader(os.Stdin)
// 		text, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Println("failed to read expression from console")
// 		}
// 		// убираем пробелы, чтобы оставить только вычислемое выражение
// 		text = strings.TrimSpace(text)
// 		// выходим, если ввели команду "exit"
// 		if text == "exit" {
// 			log.Println("aplication was successfully closed")
// 			return nil
// 		}
// 		//вычисляем выражение
// 		result, err := calculator.Calc(text)
// 		if err != nil {
// 			log.Println(text, " calculation failed wit error: ", err)
// 		} else {
// 			log.Println(text, "=", result)
// 		}
// 	}
// }

type Request struct {
	Expression string `json:"expression"`
}
type Response struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

// Функция запуска сервера
func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{}
	eg.Go(func() error {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "starting grpc server", zap.Int("port", s.listener.Addr().(*net.TCPAddr).Port))
		return s.grpcServer.Serve(s.listener)
	})
	eg.Go(func() error {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "starting rest server", zap.String("port", s.restServer.Addr))
		return s.restServer.ListenAndServe()
	})
	return eg.Wait()
}

// Функция плавного завершения
func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	l := logger.GetLoggerFromCtx(ctx)
	if l != nil {
		l.Info(ctx, "gRPC server stopped")
	}
	return s.restServer.Shutdown(ctx)
}
