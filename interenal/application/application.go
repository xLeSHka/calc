package application

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/xLeSHka/calculator/pkg/calculator"
	"github.com/xLeSHka/calculator/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Application struct {
	server *http.Server
}

func New(ctx context.Context, port int) *Application {
	mux := http.NewServeMux()
	mux.HandleFunc("/", CalcHandler)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: LoggingMiddleware(mux, ctx),
	}
	return &Application{server: server}
}

// Функция запуска приложения
// тут будем читать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
func (a *Application) Run() error {
	for {
		// читаем выражение для вычисления из командной строки
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("failed to read expression from console")
		}
		// убираем пробелы, чтобы оставить только вычислемое выражение
		text = strings.TrimSpace(text)
		// выходим, если ввели команду "exit"
		if text == "exit" {
			log.Println("aplication was successfully closed")
			return nil
		}
		//вычисляем выражение
		result, err := calculator.Calc(text)
		if err != nil {
			log.Println(text, " calculation failed wit error: ", err)
		} else {
			log.Println(text, "=", result)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}
type Response struct {
	Expression string  `json:"expression"`
	Result     float64 `json:"result"`
}

// Хендлер вычислителя выражений
func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	//декодируем выражение из тела запроса
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//пробуем вычислить выражение
	result, err := calculator.Calc(request.Expression)
	if err != nil {
		if errors.Is(err, calculator.ErrInvalidExpression) ||
			errors.Is(err, calculator.ErrDivisionByZero) ||
			errors.Is(err, calculator.ErrMissLeftParanthesis) ||
			errors.Is(err, calculator.ErrMissRightParanthesis) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		var b bytes.Buffer
		resp := Response{
			Expression: request.Expression,
			Result:     result,
		}
		err := json.NewEncoder(&b).Encode(resp)
		if err != nil {
			http.Error(w, ErrCanNotEncodeResp.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(b.Bytes()))
	}
}

// Функция старта сервера
func (s *Application) Start(ctx context.Context) error {
	eg := errgroup.Group{}
	eg.Go(func() error {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "starting server", zap.String("port", s.server.Addr))
		return s.server.ListenAndServe()
	})
	return eg.Wait()
}

// Функция плавного завершения
func (s *Application) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
