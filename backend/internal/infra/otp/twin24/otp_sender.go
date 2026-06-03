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

const (
	// VarTwinCode - имя переменной для кода подтверждения в сценарии бота
	VarTwinCode = "twin_code"
)

type OTPSender struct {
	defaultExecData string
	cidData         string
	apiClient       *Client
	authClient      *Client
	auth            *Auth
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
	}
}

// CreateTaskRequest описывает запрос на создание задания на обзвон (POST /cis/api/v1/telephony/autoCall)
type CreateTaskRequest struct {
	// Name — название задания (обязательное)
	Name string `json:"name"`

	// DefaultExec — всегда "robot" (обязательное)
	DefaultExec string `json:"defaultExec"`

	// DefaultExecData — ID сценария бота для использования в обзвоне (обязательное)
	DefaultExecData string `json:"defaultExecData"`

	// SecondExec — действие при переадресации: "end" (завершить), "ignore" (ничего не делать), "ch" (передать вызов на канал) (обязательное)
	SecondExec string `json:"secondExec"`

	// SecondExecData — ID канала для перевода (обязательное, если SecondExec = "ch")
	SecondExecData string `json:"secondExecData,omitempty"`

	// CidType — определяемый номер: "default" (по умолчанию для транка), "gornum" (один номер), "pool" (группа номеров) (обязательное)
	CidType string `json:"cidType"`

	// CidData — ID сущности (номер телефона с которого идут обзвоны) в cidType (обязательное, если cidType = "gornum" или "pool")
	CidData string `json:"cidData,omitempty"`

	// CallStrategy — стратегия обзвона: "STEP_2_STEP" (последовательная), "PARALLEL" (параллельная)
	CallStrategy string `json:"callStrategy,omitempty"`

	// StartType — режим запуска: "manual" (вручную), "time" (в указанное время) (обязательное)
	StartType string `json:"startType"`

	// StartMoment — дата и время начала обзвона в формате "ГГГГ-ММ-ДД ЧЧ:ММ" (обязательное, если startType = "time")
	StartMoment string `json:"startMoment,omitempty"`

	// CPS — интенсивность обзвона. 1 + N/100 для N звонков в секунду, или 1 - N/100 для 1 звонка в N секунд (обязательное)
	CPS float64 `json:"cps"`

	// CheckPhone — проверка корректности формата номера при добавлении кандидата
	CheckPhone bool `json:"checkPhone,omitempty"`

	// TaskComment — комментарий к заданию
	TaskComment string `json:"taskComment,omitempty"`

	// WebhookUrls — URL адреса для отправки webhook
	WebhookUrls []WebhookUrl `json:"webhookUrls,omitempty"`

	// Lifetime — срок действия задания в секундах. По истечении задание переходит в статус HALTED
	Lifetime int `json:"lifetime,omitempty"`

	// AdditionalOptions — дополнительные параметры вызовов (обязательное)
	AdditionalOptions AdditionalOptions `json:"additionalOptions"`

	// RedialStrategyOptions — настройки правил перезвона (обязательное)
	RedialStrategyOptions RedialStrategyOptions `json:"redialStrategyOptions"`

	// PhoneNormalization — нормализация номеров: nil (отключена), "RU" (нормализация РФ)
	PhoneNormalization string `json:"phoneNormalization,omitempty"`
}

// WebhookUrl описывает URL для отправки webhook уведомлений
type WebhookUrl struct {
	// URL — адрес, куда будет отправлен webhook
	URL string `json:"url,omitempty"`

	// PartialResults — отправлять ли промежуточные результаты
	PartialResults bool `json:"partialResults,omitempty"`

	// Delay — задержка перед отправкой webhook в секундах
	Delay int `json:"delay,omitempty"`

	// Events — список событий, по которым нужно отправить webhook
	// Доступные события: CALL_ENDED, CANDIDATE_CHANGED, CALL_REDIRECTED, RECALL_SCHEDULED, EFFICIENCY_REACHED, AUTOCALL_STATUS_CHANGED
	Events map[string]EventConfig `json:"events,omitempty"`
}

// EventConfig описывает настройку события для webhook
type EventConfig struct {
	// Name — имя события (например, "CALL_ENDED")
	Name string `json:"name"`

	// Value — отправлять ли webhook по данному событию
	Value bool `json:"value"`
}

// AdditionalOptions содержит дополнительные параметры вызовов (обязательное)
type AdditionalOptions struct {
	// FullListMethod — считать ли звонок результативным. Всегда "reject"
	FullListMethod string `json:"fullListMethod"`

	// FullListTime — через сколько секунд считать звонок результативным
	FullListTime int `json:"fullListTime"`

	// UseTr — учитывать ли время получателя вызова
	UseTr bool `json:"useTr,omitempty"`

	// AllowCallTimeFrom — начало интервала доступного для дозвона в секундах (если useTr = true)
	AllowCallTimeFrom int `json:"allowCallTimeFrom,omitempty"`

	// AllowCallTimeTo — конец интервала доступного для дозвона в секундах (если useTr = true)
	AllowCallTimeTo int `json:"allowCallTimeTo,omitempty"`

	// RecordCall — записывать ли звонки
	RecordCall bool `json:"recordCall"`

	// RecTrimLeft — на сколько секунд обрезать начало записи (если recordCall = true)
	RecTrimLeft int `json:"recTrimLeft,omitempty"`

	// DetectRobot — включать ли систему определения человек/робот
	DetectRobot bool `json:"detectRobot,omitempty"`

	// DetectRobotMode — режим определения: "back" (фоновая), "block" (с блокировкой)
	DetectRobotMode string `json:"detectRobotMode,omitempty"`
}

// RedialStrategyOptions содержит настройки правил перезвона (обязательное)
type RedialStrategyOptions struct {
	// RedialStrategyEn — использовать ли правила перезвона
	RedialStrategyEn bool `json:"redialStrategyEn"`

	// CandidateLimit — максимальное количество вызовов кандидату
	CandidateLimit *LimitConfig `json:"candidateLimit,omitempty"`

	// NumberLimit — максимальное количество вызовов по номеру
	NumberLimit *LimitConfig `json:"numberLimit,omitempty"`

	// Busy — правило перезвона при статусе "Занято" (обязательное)
	Busy RedialRule `json:"busy"`

	// NoAnswer — правило перезвона при статусе "Нет ответа" (обязательное)
	NoAnswer RedialRule `json:"noAnswer"`

	// AnswerMash — правило перезвона при статусе "Автоответчик" (обязательное)
	AnswerMash RedialRule `json:"answerMash"`

	// Congestion — правило перезвона при статусе "Ошибка вызова" (обязательное)
	Congestion RedialRule `json:"congestion"`

	// AnswerNoList — правило перезвона если вызов нерезультативен (обязательное)
	AnswerNoList RedialRule `json:"answerNoList"`
}

// LimitConfig описывает лимит на количество вызовов
type LimitConfig struct {
	// Redial — активировать ли лимит по максимальному количеству вызовов
	Redial bool `json:"redial"`

	// Count — максимальное количество вызовов (если redial = true)
	Count int `json:"count,omitempty"`
}

// RedialRule описывает правило перезвона для конкретного статуса
type RedialRule struct {
	// Redial — активировать ли правило перезвона
	Redial bool `json:"redial"`

	// Time — задержка перед перезвоном в секундах (если redial = true)
	Time int `json:"time,omitempty"`

	// Count — количество перезвонов (если redial = true)
	Count int `json:"count,omitempty"`
}

type CreateTaskResponse struct {
	ID struct {
		Identity string `json:"identity"`
	} `json:"id"`
}

// AddCandidatesRequest описывает запрос на добавление кандидатов (POST /cis/api/v1/telephony/autoCallCandidate/batch)
type AddCandidatesRequest struct {
	// ForceStart — если true, задание автоматически запустится после добавления кандидатов
	ForceStart bool `json:"forceStart,omitempty"`

	// Batch — массив кандидатов для обзвона (обязательное)
	Batch []CandidateBatch `json:"batch"`
}

// CandidateBatch описывает одного кандидата для добавления в задание
type CandidateBatch struct {
	// Phone — номера телефона кандидата (обязательное)
	Phone []string `json:"phone"`

	// Variables — объект с переменными по кандидату (доступны в сценарии бота)
	Variables map[string]string `json:"variables,omitempty"`

	// CallbackData — данные, которые вернутся в webhook о результате звонка (например, ID клиента)
	CallbackData map[string]string `json:"callbackData,omitempty"`

	// ClientExternalID — внешний идентификатор клиента
	ClientExternalID string `json:"clientExternalId,omitempty"`

	// Timezone — таймзона кандидата (отклонение от UTC+0 в минутах)
	Timezone int `json:"timezone,omitempty"`

	// AutoCallID — идентификатор задания на обзвон (обязательное)
	AutoCallID string `json:"autoCallId"`
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

	taskReq := CreateTaskRequest{
		Name:            fmt.Sprintf("OTP %s", phone),
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
			FullListTime:   12, // успешный звонок через 12 секунд
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

	// подготавливаем переменные для бота
	// VarTwinCode обязательна, другие переменные можно добавлять по мере необходимости
	variables := map[string]string{
		VarTwinCode: code,
	}

	// добавляем кандидата с OTP кодом
	if err := s.AddCandidates(ctx, AddCandidatesRequest{
		ForceStart: false,
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
	time.Sleep(2 * time.Second)

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
