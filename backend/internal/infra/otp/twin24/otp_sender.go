package twin24

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type OTPSender struct {
	client *Client
	auth   *Auth
}

func NewOTPSender(client *Client, auth *Auth) *OTPSender {
	return &OTPSender{client: client, auth: auth}
}

func (s *OTPSender) SendOTP(ctx context.Context, phone, code string) error {
	_, err := s.auth.GetToken(ctx)
	if err != nil {
		return fmt.Errorf("get token: %w", err)
	}

	// taskID, err := s.createTask(ctx, token)
	// if err != nil {
	// 	return fmt.Errorf("create task: %w", err)
	// }

	// ждём пока задание инициализируется
	time.Sleep(5 * time.Second)

	// if err := s.addCandidate(ctx, token, taskID, phone, code); err != nil {
	// 	return fmt.Errorf("add candidate: %w", err)
	// }

	return nil
}

func (s *OTPSender) doWithAuth(ctx context.Context, method, path string, body any) (*http.Response, error) {
	token, err := s.auth.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	// graceful refresh — токен протух
	if resp.StatusCode == http.StatusUnauthorized {
		resp.Body.Close()
		token, err = s.auth.RefreshToken(ctx)
		if err != nil {
			return nil, err
		}
		_ = token
		return s.client.Do(ctx, method, path, body)
	}

	return resp, nil
}
