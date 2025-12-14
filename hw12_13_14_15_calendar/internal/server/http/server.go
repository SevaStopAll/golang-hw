package internalhttp

import (
	"context"
	"fmt"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/app"
	"github.com/sevastopall/hw12_13_14_15_calendar/internal/logger"
	"net/http"
	"strconv"
)

type Server struct {
	host        string
	port        int
	logger      *logger.Logger
	application *app.App
}

func NewServer(logger *logger.Logger, app *app.App, host string, port int) *Server {
	return &Server{host: host,
		port:        port,
		logger:      logger,
		application: app}
}

func (s *Server) Start(ctx context.Context) error {
	http.HandleFunc("/", handler)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	err := http.ListenAndServe(s.host+":"+strconv.Itoa(s.port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func hello(w http.ResponseWriter, req *http.Request) {
	// Функции, служащие обработчиками, принимают
	// `http.ResponseWriter` и `http.Request` в качестве
	// аргументов. Объект `http.ResponseWriter` используется для заполнения
	// HTTP-ответа. Здесь наш простой ответ - это просто
	// "hello\n".
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	// Этот обработчик делает что-то более сложное,
	// прочитав все HTTP-заголовки запроса и выведя их в тело ответа.
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
