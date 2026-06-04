package twin24

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type OTPSender struct {
	defaultExecData string
	cidData         string
	apiClient       *Client
	authClient      *Client
	auth            *Auth
	redisClient     *redis.Client
}

// NewOTPSenderFromEnv создает OTPSender, читая конфигурацию из переменных окружения.
// Возвращает nil, если обязательные переменные не заданы.
func NewOTPSenderFromEnv(redisClient *redis.Client) *OTPSender {
	email := os.Getenv("TWIN24_EMAIL")
	password := os.Getenv("TWIN24_PASSWORD")
	baseURL := os.Getenv("TWIN24_BASE_URL")
	botID := os.Getenv("TWIN24_BOT_SCENARIO_ID")
	cid := os.Getenv("TWIN24_BOT_CID")
	iamURL := os.Getenv("TWIN24_IAM_URL")

	// Значения по умолчанию
	if baseURL == "" {
		baseURL = "https://twin24.ai"
	}
	if iamURL == "" {
		iamURL = "https://iam.twin24.ai"
	}

	if email == "" || password == "" || botID == "" || cid == "" {
		slog.Warn("⚠️ Twin24 credentials not configured, OTP calls disabled")
		return nil
	}

	authClient := NewClient(iamURL)
	auth := NewAuth(authClient, *redisClient, email, password)

	return &OTPSender{
		defaultExecData: botID,
		cidData:         cid,
		apiClient:       NewClient(baseURL),
		authClient:      authClient,
		auth:            auth,
		redisClient:     redisClient,
	}
}

// getOrCreateDailyTask возвращает ID задания на сегодня.
// Если задания нет или оно создано в другой день — создаёт новое.
// Использует Redis для хранения состояния между перезапусками.
func (s *OTPSender) getOrCreateDailyTask(ctx context.Context) (string, error) {
	now := time.Now()
	taskKey := fmt.Sprintf("twin24:daily_task:%s", now.Format("2006-01-02"))

	// проверяем, есть ли задание на сегодня в Redis
	taskID, err := s.redisClient.Get(ctx, taskKey).Result()
	if err == nil && taskID != "" {
		return taskID, nil
	}

	// создаём новое задание на сегодня
	taskReq := CreateTaskRequest{
		Name:            fmt.Sprintf("OTP Одной Крови %s ", now.Format("2006-01-02")),
		DefaultExec:     "robot",
		DefaultExecData: s.defaultExecData,
		SecondExec:      "end",
		CidType:         "gornum",
		CidData:         s.cidData,
		StartType:       "manual",
		CPS:             1.0,
		WebhookUrls: []WebhookUrl{
			{
				URL: "https://n8n.rmay1er.ru/webhook/twin-webhook",
				Events: map[string]EventConfig{
					"CALL_ENDED":         {Name: "CALL_ENDED", Value: true},
					"CANDIDATE_CHANGED":  {Name: "CANDIDATE_CHANGED", Value: true},
					"CALL_REDIRECTED":    {Name: "CALL_REDIRECTED", Value: true},
					"RECALL_SCHEDULED":   {Name: "RECALL_SCHEDULED", Value: true},
					"EFFICIENCY_REACHED": {Name: "EFFICIENCY_REACHED", Value: true},
				},
			},
		},
		AdditionalOptions: AdditionalOptions{
			FullListMethod: "reject",
			FullListTime:   12,
			RecordCall:     false,
		},
		RedialStrategyOptions: RedialStrategyOptions{
			RedialStrategyEn: true,
			CandidateLimit: &LimitConfig{
				Redial: true,
				Count:  6,
			},
			Busy:         RedialRule{Redial: false},
			NoAnswer:     RedialRule{Redial: true, Time: 26, Count: 2},
			AnswerMash:   RedialRule{Redial: false},
			Congestion:   RedialRule{Redial: false},
			AnswerNoList: RedialRule{Redial: false},
		},
		PhoneNormalization: "RU",
		// Lifetime:           300,
	}

	newTaskID, err := s.CreateTask(ctx, taskReq)
	if err != nil {
		return "", fmt.Errorf("create daily task: %w", err)
	}

	// сохраняем в Redis с TTL 48 часов (чтобы пережило день)
	err = s.redisClient.Set(ctx, taskKey, newTaskID, 48*time.Hour).Err()
	if err != nil {
		slog.Warn("⚠️ failed to save task ID to Redis", "error", err)
		// не критично, продолжаем работу
	}

	// ждём инициализации задания (рекомендация Twin24)
	time.Sleep(5 * time.Second)

	return newTaskID, nil
}

// AddCandidates добавляет кандидатов в задание на обзвон
func (s *OTPSender) AddCandidates(ctx context.Context, req AddCandidatesRequest) error {
	resp, err := s.doWithAuth(ctx, http.MethodPost, "/cis/api/v1/telephony/autoCallCandidate/batch", req)
	if err != nil {
		return fmt.Errorf("add candidates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}

// CreateTask создаёт задание на автоматический обзвон
func (s *OTPSender) CreateTask(ctx context.Context, req CreateTaskRequest) (string, error) {
	resp, err := s.doWithAuth(ctx, http.MethodPost, "/cis/api/v1/telephony/autoCall", req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}

	var result CreateTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.ID.Identity, nil
}

// StartTask запускает задание на обзвон
func (s *OTPSender) StartTask(ctx context.Context, taskID string) error {
	resp, err := s.doWithAuth(ctx, http.MethodPost, "/cis/api/v1/telephony/autoCall/"+taskID+"/play", nil)
	if err != nil {
		return fmt.Errorf("start task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

// PauseTask ставит задание на паузу
func (s *OTPSender) PauseTask(ctx context.Context, taskID string) error {
	resp, err := s.doWithAuth(ctx, http.MethodPost, "/cis/api/v1/telephony/autoCall/"+taskID+"/pause", nil)
	if err != nil {
		return fmt.Errorf("pause task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

// HaltTask останавливает задание навсегда
func (s *OTPSender) HaltTask(ctx context.Context, taskID string) error {
	resp, err := s.doWithAuth(ctx, http.MethodPost, "/cis/api/v1/telephony/autoCall/"+taskID+"/halted", nil)
	if err != nil {
		return fmt.Errorf("halt task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (s *OTPSender) SendOTP(ctx context.Context, phone, code string) error {
	// получаем или создаём задание на сегодня
	taskID, err := s.getOrCreateDailyTask(ctx)
	if err != nil {
		return fmt.Errorf("get or create daily task: %w", err)
	}

	// подготавливаем переменные для бота
	// VarTwinCode обязательна, другие переменные можно добавлять по мере необходимости
	variables := map[string]string{
		VarTwinCode: code,
	}

	// добавляем кандидата с OTP кодом
	if err := s.AddCandidates(ctx, AddCandidatesRequest{
		ForceStart: true, // запускаем задание сразу после добавления
		Batch: []CandidateBatch{
			{
				Phone:      []string{phone},
				Variables:  variables,
				AutoCallID: taskID,
			},
		},
	}); err != nil {
		return fmt.Errorf("add candidate: %w", err)
	}

	// ждём перед стартом (рекомендация Twin24)
	// time.Sleep(3 * time.Second)

	// запускаем задание
	// if err := s.StartTask(ctx, taskID); err != nil {
	// 	return fmt.Errorf("start task: %w", err)
	// }

	return nil
}

func (s *OTPSender) doWithAuth(ctx context.Context, method, path string, body any) (*http.Response, error) {
	token, err := s.auth.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := s.apiClient.DoWithBearer(ctx, method, path, body, token)
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
		return s.apiClient.DoWithBearer(ctx, method, path, body, token)
	}

	return resp, nil
}
