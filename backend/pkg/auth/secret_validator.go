package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

// Provider тип для определения платформы Web App
type Provider string

const (
	ProviderTelegram Provider = "telegram_bot"
	ProviderMax      Provider = "max_bot"
)

// ParseProvider преобразует строку в Provider
func ParseProvider(s string) (Provider, error) {
	switch s {
	case "telegram", "tg":
		return ProviderTelegram, nil
	case "max":
		return ProviderMax, nil
	default:
		return "", fmt.Errorf("unknown provider: %s", s)
	}
}

// WebAppInitData содержит все распарсенные данные из initData (кроме hash).
type WebAppInitData struct {
	User       *WebAppUserInfo
	AuthDate   int64
	StartParam string
	ChatType   string
	RawParams  map[string]string
}

// WebAppUserInfo содержит информацию о пользователе из Web App.
type WebAppUserInfo struct {
	ID           string
	FirstName    string
	LastName     string
	Username     string
	LanguageCode string
	PhotoURL     *string
}

// AppValidator валидирует данные для конкретного провайдера.
type AppValidator struct {
	maxBotToken string
	tgBotToken  string
}

// NewAppValidator создает валидатор для конкретного провайдера.
func NewAppValidator(maxBotToken, telegramBotToken string) *AppValidator {
	return &AppValidator{
		maxBotToken: maxBotToken,
		tgBotToken:  telegramBotToken,
	}
}

// ValidateWebAppInitData валидирует initData.
func (v *AppValidator) ValidateWebAppInitData(ctx context.Context, initData string, provider authmodel.ProviderName) (*WebAppInitData, error) {
	if initData == "" {
		return nil, errors.New("initData is empty")
	}

	p := Provider(provider)
	var botToken string

	switch p {
	case ProviderTelegram:
		botToken = v.tgBotToken
	case ProviderMax:
		// Для Max используем старую логику валидации
		return v.validateMaxInitData(initData)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Проверяем, какой тип подписи используется в initData
	// Если есть поле "signature" - используем third-party валидацию (Ed25519)
	// Если есть только "hash" - используем старую валидацию (HMAC-SHA256)
	hasSignature := strings.Contains(initData, "signature=")

	expIn := 24 * time.Hour
	var parsed initdata.InitData
	var err error

	if hasSignature {
		// Third-party валидация (Ed25519) - использует Bot ID
		botID, err := extractBotID(botToken)
		if err != nil {
			return nil, fmt.Errorf("failed to extract bot ID from token: %w", err)
		}
		if err := initdata.ValidateThirdParty(initData, botID, expIn); err != nil {
			return nil, fmt.Errorf("invalid init data (third-party): %w", err)
		}
		parsed, err = initdata.Parse(initData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse init data: %w", err)
		}
	} else {
		// Старая валидация (HMAC-SHA256) - использует Bot Token
		if err := initdata.Validate(initData, botToken, expIn); err != nil {
			return nil, fmt.Errorf("invalid init data: %w", err)
		}
		parsed, err = initdata.Parse(initData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse init data: %w", err)
		}
	}

	result := &WebAppInitData{
		AuthDate:  int64(parsed.AuthDateRaw),
		RawParams: make(map[string]string),
	}

	// Заполняем информацию о пользователе (User - не указатель, проверяем по ID)
	if parsed.User.ID > 0 {
		user := &WebAppUserInfo{
			ID:           strconv.FormatInt(parsed.User.ID, 10),
			FirstName:    parsed.User.FirstName,
			LastName:     parsed.User.LastName,
			Username:     parsed.User.Username,
			LanguageCode: parsed.User.LanguageCode,
		}
		if parsed.User.PhotoURL != "" {
			user.PhotoURL = &parsed.User.PhotoURL
		}
		result.User = user
	}

	// Заполняем дополнительные параметры
	if parsed.StartParam != "" {
		result.StartParam = parsed.StartParam
	}
	// Chat - не указатель, проверяем через Known()
	if parsed.Chat.Type.Known() {
		result.ChatType = string(parsed.Chat.Type)
	}

	// Сохраняем все параметры в RawParams
	if parsed.QueryID != "" {
		result.RawParams["query_id"] = parsed.QueryID
	}
	if parsed.User.ID > 0 {
		userJSON, _ := json.Marshal(parsed.User)
		result.RawParams["user"] = url.QueryEscape(string(userJSON))
	}
	if parsed.Chat.Type.Known() {
		result.RawParams["chat_instance"] = strconv.FormatInt(parsed.ChatInstance, 10)
		result.RawParams["chat_type"] = string(parsed.Chat.Type)
	}
	result.RawParams["auth_date"] = strconv.Itoa(parsed.AuthDateRaw)
	if parsed.StartParam != "" {
		result.RawParams["start_param"] = parsed.StartParam
	}

	return result, nil
}

// extractBotID извлекает ID бота из токена (формат: "123456789:ABCdefGHIjklMNOpqrsTUVwxyz")
func extractBotID(token string) (int64, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 || parts[0] == "" {
		return 0, fmt.Errorf("invalid token format: expected 'botID:botToken', got '%s'", token)
	}
	botID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse bot ID: %w", err)
	}
	return botID, nil
}

// validateMaxInitData валидирует initData для Max (оставляем старую логику)
func (v *AppValidator) validateMaxInitData(initData string) (*WebAppInitData, error) {
	rawParams, err := parseRawQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse initData: %w", err)
	}

	hash, ok := rawParams["hash"]
	if !ok {
		return nil, errors.New("missing hash parameter")
	}
	delete(rawParams, "hash")
	delete(rawParams, "signature")

	authDateStr, ok := rawParams["auth_date"]
	if !ok {
		return nil, errors.New("missing auth_date parameter")
	}
	authDate, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid auth_date: %w", err)
	}
	if time.Now().Unix()-authDate > 24*60*60 {
		return nil, errors.New("auth_date is too old")
	}

	keys := make([]string, 0, len(rawParams))
	for k := range rawParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var checkStrings []string
	for _, k := range keys {
		checkStrings = append(checkStrings, k+"="+rawParams[k])
	}
	dataCheckString := joinStrings(checkStrings, "\n")

	secretKey := v.computeSecretKey(ProviderMax) // Pass provider to computeSecretKey

	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(dataCheckString))
	computedHash := hex.EncodeToString(h.Sum(nil))

	computedBytes, err := hex.DecodeString(computedHash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode computed hash: %w", err)
	}
	providedBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, fmt.Errorf("invalid hash format: %w", err)
	}
	if !hmac.Equal(computedBytes, providedBytes) {
		return nil, errors.New("invalid hash: data has been tampered with")
	}

	result := &WebAppInitData{
		AuthDate:  authDate,
		RawParams: make(map[string]string),
	}
	maps.Copy(result.RawParams, rawParams)

	userJSON, ok := rawParams["user"]
	if ok {
		userJSONDecoded, err := url.QueryUnescape(userJSON)
		if err != nil {
			return nil, fmt.Errorf("failed to decode user parameter: %w", err)
		}

		var user WebAppUserInfo
		if err := json.Unmarshal([]byte(userJSONDecoded), &user); err != nil {
			return nil, fmt.Errorf("failed to parse user data: %w", err)
		}
		result.User = &user
	}

	if startParam, ok := rawParams["start_param"]; ok {
		if decoded, err := url.QueryUnescape(startParam); err == nil {
			result.StartParam = decoded
		}
	}

	if chatType, ok := rawParams["chat_type"]; ok {
		result.ChatType = chatType
	}

	return result, nil
}

// computeSecretKey вычисляет secret_key в зависимости от провайдера.
func (v *AppValidator) computeSecretKey(p Provider) []byte {
	var botToken string
	switch p {
	case ProviderTelegram:
		botToken = v.tgBotToken
	case ProviderMax:
		botToken = v.maxBotToken
	default:
		// Should not happen with current logic, but as a fallback
		botToken = v.tgBotToken
	}

	h := hmac.New(sha256.New, []byte(botToken))
	h.Write([]byte("WebAppData"))
	return h.Sum(nil)
}

func parseRawQuery(qs string) (map[string]string, error) {
	result := make(map[string]string)
	pairs := strings.Split(qs, "&")
	for _, pair := range pairs {
		if pair == "" {
			continue
		}
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid query parameter: %s", pair)
		}
		result[kv[0]] = kv[1]
	}
	return result, nil
}

func joinStrings(ss []string, sep string) string {
	if len(ss) == 0 {
		return ""
	}
	result := ss[0]
	for i := 1; i < len(ss); i++ {
		result += sep + ss[i]
	}
	return result
}

// ValidateMock оставлена для совместимости.
func (v *AppValidator) ValidateMock(initData string) (providerID string, role string) {
	splited := strings.Split(initData, "-")
	if len(splited) < 2 {
		return "", ""
	}
	// This mock validation should probably be updated to consider the provider
	// For now, it uses tgBotToken as a default for comparison
	if splited[1] != v.tgBotToken { // Assuming mock is for Telegram or a generic token
		return "", ""
	}
	providerID = splited[0]
	role = "USER"
	return providerID, role
}
