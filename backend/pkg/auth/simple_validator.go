package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

// SimpleValidator - упрощенный валидатор для извлечения данных из initData.
// Не выполняет валидацию подписи, только парсит и извлекает данные.
// Используется для тестирования или в случаях, когда валидация подписи не требуется.
type SimpleValidator struct{}

// NewSimpleValidator создает новый экземпляр SimpleValidator.
func NewSimpleValidator() *SimpleValidator {
	return &SimpleValidator{}
}

// ValidateWebAppInitData извлекает данные из initData без валидации подписи.
// Возвращает *WebAppInitData с извлеченными данными или ошибку.
func (v *SimpleValidator) ValidateWebAppInitData(ctx context.Context, initData string, provider authmodel.ProviderName) (*WebAppInitData, error) {
	if initData == "" {
		return nil, errors.New("initData is empty")
	}

	parsed, err := initdata.Parse(initData)
	if err != nil {
		return nil, err
	}

	result := &WebAppInitData{
		AuthDate:  int64(parsed.AuthDateRaw),
		RawParams: make(map[string]string),
	}

	// Извлекаем информацию о пользователе
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

	// Извлекаем StartParam
	if parsed.StartParam != "" {
		result.StartParam = parsed.StartParam
	}

	// Извлекаем ChatType
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
