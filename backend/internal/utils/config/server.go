package config

import (
	"fmt"
	"net/http"
)

type MyServer struct {
	*http.Server
}

// NewServer создает новый экземпляр сервера с заданным mux
func NewServer(port int, mux http.Handler) *MyServer {
	return &MyServer{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

// Use добавляет middleware в цепочку обработки.
// Принимает стандартные функции-обертки func(http.Handler) http.Handler.
// Middleware применяются в порядке их добавления (первый добавленный будет внешним).
func (s *MyServer) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, mw := range middlewares {
		s.Handler = mw(s.Handler)
	}
}
