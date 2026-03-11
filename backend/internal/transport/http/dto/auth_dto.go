package dto

// ============================================
// SignIn User (for Service)
// ============================================

// ServiceSignInInput представляет запрос на вход пользователя через сервис
type ServiceSignInInput struct {
	InternalKey string `header:"X-Internal-Key" doc:"Секретный ключ для ботов и клиник"`
	Body        ServiceSignInBody
}

// ServiceSignInBody представляет тело запроса на вход пользователя через сервис
type ServiceSignInBody struct {
	FullName     string          `json:"fullName,omitempty" doc:"Полное имя пользователя" minLength:"2" maxLength:"255" example:"Иван Иванов"`
	ProviderID   string          `json:"providerId" doc:"ID провайдера" example:"123456789"`
	ProviderName string          `json:"providerName" doc:"Имя провайдера" enum:"telegram_bot,max_bot,service"`
	MetaData     *map[string]any `json:"metaData,omitempty" doc:"Метаданные пользователя"`
}

// ServiceSignInOutput представляет ответ на вход пользователя
type ServiceSignInOutput struct {
	Body ServiceSignInResult
}

// ServiceSignInResult представляет результат входа пользователя
type ServiceSignInResult struct {
	UserID      string `json:"userId" doc:"ID пользователя на портале" example:"USR-ABCDEABCDE"`
	AccessToken string `json:"accessToken,omitempty" doc:"Access токен для аутентификации" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"tokenType" doc:"Тип токена" example:"Bearer"`
	ExpiresAt   string `json:"expiresAt" doc:"Время истечения токена" example:"2023-12-31T23:59:59Z"`
}

// ============================================
// SignIn User (for MiniApp)
// ============================================

// MiniAppSignInInput представляет запрос на вход пользователя через мини-приложение
type MiniAppSignInInput struct {
	Body MiniAppSignInBody
}

// MiniAppSignInBody представляет тело запроса на вход пользователя через мини-приложение
type MiniAppSignInBody struct {
	AppInitData string          `json:"appInitData" doc:"Зашифрованный токен бота для сверки"`
	MetaData    *map[string]any `json:"metaData,omitempty" doc:"Метаданные пользователя"`
}

// MiniAppSignInOutput представляет ответ на вход пользователя
type MiniAppSignInOutput struct {
	Body MiniAppSignInResult
}

// MiniAppSignInResult представляет результат входа пользователя
type MiniAppSignInResult struct {
	UserID      string `json:"userId" doc:"ID пользователя на портале" example:"USR-ABCDEABCDE"`
	AccessToken string `json:"accessToken,omitempty" doc:"Access токен для аутентификации" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"tokenType" doc:"Тип токена" example:"Bearer"`
	ExpiresAt   string `json:"expiresAt" doc:"Время истечения токена" example:"2023-12-31T23:59:59Z"`
}
