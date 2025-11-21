package v1

import (
	"context"

	greetv1 "github.com/artesipov-alt/odnoi-krovi-app/poolservice/gen/api/greet/v1"
)

type GreetServer struct{}

// Greet обрабатывает запрос приветствия
func (s *GreetServer) Greet(
	ctx context.Context,
	req *greetv1.GreetRequest,
) (*greetv1.GreetResponse, error) {
	return &greetv1.GreetResponse{
		Greeting: "Привет, " + req.Name + "!",
	}, nil
}
