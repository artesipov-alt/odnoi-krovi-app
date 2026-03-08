package dto

// ============================================
// SignIn User (for Messengers)
// ============================================

// MessengerSignInInput представляет запрос на вход пользователя через мессенджер
type MessengerSignInInput struct {
	InternalKey string `header:"X-Internal-Key" doc:"Секретный ключ для ботов и клиник"`
	Body        MessengerSignInBody
}

// MessengerSignInBody представляет тело запроса на вход пользователя через мессенджер
type MessengerSignInBody struct {
	FullName   *string         `json:"fullName,omitempty" doc:"Полное имя пользователя" minLength:"2" maxLength:"255" example:"Иван Иванов"`
	ProviderID int64           `json:"providerId" doc:"ID провайдера" example:"123456789"`
	MetaData   *map[string]any `json:"metaData,omitempty" doc:"Метаданные пользователя"`
}

// MessengerSignInOutput представляет ответ на вход пользователя
type MessengerSignInOutput struct {
	Body MessengerSignInResult
}

// MessengerSignInResult представляет результат входа пользователя
type MessengerSignInResult struct {
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
