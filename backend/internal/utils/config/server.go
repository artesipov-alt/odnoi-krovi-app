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
// Middleware применяются так, что первое переданное в списке становится самым внешним слоем.
// Это позволяет соблюдать логический порядок: Recovery -> RequestID -> Logging -> Mux.
func (s *MyServer) Use(middlewares ...func(http.Handler) http.Handler) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		s.Handler = middlewares[i](s.Handler)
	}
}
