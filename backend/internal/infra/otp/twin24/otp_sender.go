package twin24

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type OTPSender struct {
	apiClient  *Client
	authClient *Client
	auth       *Auth
}

func NewOTPSender(apiBaseURL string, authClient *Client, auth *Auth) *OTPSender {
	return &OTPSender{
		apiClient:  NewClient(apiBaseURL),
		authClient: authClient,
		auth:       auth,
	}
}

type CreateTaskRequest struct {
	Name                  string                `json:"name"`
	DefaultExec           string                `json:"defaultExec"`
	DefaultExecData       string                `json:"defaultExecData"`
	SecondExec            string                `json:"secondExec"`
	SecondExecData        string                `json:"secondExecData,omitempty"`
	CidType               string                `json:"cidType"`
	CidData               string                `json:"cidData,omitempty"`
	CallStrategy          string                `json:"callStrategy,omitempty"`
	StartType             string                `json:"startType"`
	StartMoment           string                `json:"startMoment,omitempty"`
	CPS                   float64               `json:"cps"`
	CheckPhone            bool                  `json:"checkPhone,omitempty"`
	TaskComment           string                `json:"taskComment,omitempty"`
	WebhookUrls           []WebhookUrl          `json:"webhookUrls,omitempty"`
	Lifetime              int                   `json:"lifetime,omitempty"`
	AdditionalOptions     AdditionalOptions     `json:"additionalOptions"`
	RedialStrategyOptions RedialStrategyOptions `json:"redialStrategyOptions"`
	PhoneNormalization    string                `json:"phoneNormalization,omitempty"`
}

type WebhookUrl struct {
	URL            string                 `json:"url,omitempty"`
	PartialResults bool                   `json:"partialResults,omitempty"`
	Delay          int                    `json:"delay,omitempty"`
	Events         map[string]EventConfig `json:"events,omitempty"`
}

type EventConfig struct {
	Name  string `json:"name"`
	Value bool   `json:"value"`
}

type AdditionalOptions struct {
	FullListMethod    string `json:"fullListMethod"`
	FullListTime      int    `json:"fullListTime"`
	UseTr             bool   `json:"useTr,omitempty"`
	AllowCallTimeFrom int    `json:"allowCallTimeFrom,omitempty"`
	AllowCallTimeTo   int    `json:"allowCallTimeTo,omitempty"`
	RecordCall        bool   `json:"recordCall"`
	RecTrimLeft       int    `json:"recTrimLeft,omitempty"`
	DetectRobot       bool   `json:"detectRobot,omitempty"`
	DetectRobotMode   string `json:"detectRobotMode,omitempty"`
}

type RedialStrategyOptions struct {
	RedialStrategyEn bool         `json:"redialStrategyEn"`
	CandidateLimit   *LimitConfig `json:"candidateLimit,omitempty"`
	NumberLimit      *LimitConfig `json:"numberLimit,omitempty"`
	Busy             RedialRule   `json:"busy"`
	NoAnswer         RedialRule   `json:"noAnswer"`
	AnswerMash       RedialRule   `json:"answerMash"`
	Congestion       RedialRule   `json:"congestion"`
	AnswerNoList     RedialRule   `json:"answerNoList"`
}

type LimitConfig struct {
	Redial bool `json:"redial"`
	Count  int  `json:"count,omitempty"`
}

type RedialRule struct {
	Redial bool `json:"redial"`
	Time   int  `json:"time,omitempty"`
	Count  int  `json:"count,omitempty"`
}

type CreateTaskResponse struct {
	ID struct {
		Identity string `json:"identity"`
	} `json:"id"`
}

type AddCandidatesRequest struct {
	ForceStart bool             `json:"forceStart,omitempty"`
	Batch      []CandidateBatch `json:"batch"`
}

type CandidateBatch struct {
	Phone            []string          `json:"phone"`
	Variables        map[string]string `json:"variables,omitempty"`
	CallbackData     map[string]string `json:"callbackData,omitempty"`
	ClientExternalID string            `json:"clientExternalId,omitempty"`
	Timezone         int               `json:"timezone,omitempty"`
	AutoCallID       string            `json:"autoCallId"`
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
	// TODO: получить botScenarioID из конфигурации
	botScenarioID := os.Getenv("TWIN24_BOT_SCENARIO_ID")

	taskReq := CreateTaskRequest{
		Name:            fmt.Sprintf("OTP %s", phone),
		DefaultExec:     "robot",
		DefaultExecData: botScenarioID,
		SecondExec:      "end",
		CidType:         "gornum",
		CidData:         "12b5eaa8-925d-4e60-a02a-69049a89a916",
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
			FullListTime:   30, // успешный звонок через 10 секунд
			RecordCall:     false,
		},
		RedialStrategyOptions: RedialStrategyOptions{
			RedialStrategyEn: false, // перезвоны не нужны для OTP
			Busy:             RedialRule{Redial: false},
			NoAnswer:         RedialRule{Redial: false},
			AnswerMash:       RedialRule{Redial: false},
			Congestion:       RedialRule{Redial: false},
			AnswerNoList:     RedialRule{Redial: false},
		},
		PhoneNormalization: "RU",
		Lifetime:           300, // 5 минут
	}

	taskID, err := s.CreateTask(ctx, taskReq)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}

	// ждём инициализации задания (рекомендация Twin24)
	time.Sleep(5 * time.Second)

	// добавляем кандидата с OTP кодом
	if err := s.AddCandidates(ctx, AddCandidatesRequest{
		ForceStart: false,
		Batch: []CandidateBatch{
			{
				Phone:      []string{phone},
				Variables:  map[string]string{"twin_code": code},
				AutoCallID: taskID,
			},
		},
	}); err != nil {
		return fmt.Errorf("add candidate: %w", err)
	}

	// ждём перед стартом (рекомендация Twin24)
	time.Sleep(5 * time.Second)

	// запускаем задание
	if err := s.StartTask(ctx, taskID); err != nil {
		return fmt.Errorf("start task: %w", err)
	}

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
