package config

import "net/http"

type MyServer struct {
	*http.Server
}

// NewServer создает новый экземпляр сервера с заданным mux
func NewServer(mux http.Handler) *MyServer {
	return &MyServer{
		&http.Server{
			Addr:    GetEnv("SERVER_PORT", ":8080"),
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
