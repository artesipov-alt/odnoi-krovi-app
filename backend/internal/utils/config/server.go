package config

import "net/http"

type MyServer struct {
	*http.Server
}

// NewServerConfig создает конфигурацию сервера из переменных окружения
func NewServer(mux http.Handler) *MyServer {
	return &MyServer{
		&http.Server{
			Addr:    GetEnv("SERVER_PORT", ":8080"),
			Handler: mux,
		},
	}
}

func (s *MyServer) Use(middleware ...http.Handler) error {
	for _, m := range middleware {
		s.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.ServeHTTP(w, r)
			s.Handler.ServeHTTP(w, r)
		})
	}
	return nil
}
