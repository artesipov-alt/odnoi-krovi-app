package http

import (
	"context"
	"net/http"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/auth/cmd"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/danielgtaylor/huma/v2"
)

// AuthHandler обрабатывает HTTP запросы для операций с пользователями
type AuthHandler struct {
	appAuthHandler      *cmd.MiniAppAuthHandler
	externalAuthHandler *cmd.ExternalAuthHandler
}

func NewAuthHandler(appAuthHandler *cmd.MiniAppAuthHandler, externalAuthHandler *cmd.ExternalAuthHandler) *AuthHandler {
	return &AuthHandler{
		appAuthHandler:      appAuthHandler,
		externalAuthHandler: externalAuthHandler,
	}
}

// Register регистрирует маршруты пользователя в Huma API
func (h *AuthHandler) Register(api huma.API) {
	// Аутентификация пользователя через Telegram
	huma.Register(api, huma.Operation{
		OperationID:   "auth-user-via-telegram",
		Method:        http.MethodPost,
		Path:          "/v1/auth/signin/telegram",
		Summary:       "Аутентификация пользователя через Telegram",
		Description:   "Аутентифицирует пользователя в системе через Telegram",
		Tags:          []string{"auth-v1"},
		DefaultStatus: http.StatusOK,
	}, h.TelegramSignIn)

	// Аутентификация пользователя через Max
	huma.Register(api, huma.Operation{
		OperationID:   "auth-user-via-max",
		Method:        http.MethodPost,
		Path:          "/v1/auth/signin/max",
		Summary:       "Аутентификация пользователя через Max",
		Description:   "Аутентифицирует пользователя в системе через Max",
		Tags:          []string{"auth-v1"},
		DefaultStatus: http.StatusOK,
	}, h.MaxSignIn)

	// Аутентификация пользователя через внешний сервис
	huma.Register(api, huma.Operation{
		OperationID:   "auth-user-via-service",
		Method:        http.MethodPost,
		Path:          "/v1/auth/signin/service",
		Summary:       "Аутентификация пользователя через внешний сервис",
		Description:   "Аутентифицирует пользователя в системе через внешний сервис",
		Tags:          []string{"auth-v1"},
		DefaultStatus: http.StatusOK,
	}, h.ServiceSignIn)
}

func (h *AuthHandler) TelegramSignIn(ctx context.Context, input *dto.MiniAppSignInInput) (*dto.MiniAppSignInOutput, error) {
	idn := &authmodel.Identity{
		ProviderName: authmodel.ProviderTelegram,
		AppInitData:  input.Body.AppInitData,
		Metadata:     input.Body.MetaData,
	}

	var metadata *authmodel.Metadata
	if input.Body.MetaData != nil {
		metadata = authmodel.NewUserMetadata(*input.Body.MetaData)
	}

	authdata, err := h.appAuthHandler.Handle(ctx, idn, metadata)
	if err != nil {
		return nil, err
	}

	return &dto.MiniAppSignInOutput{Body: dto.MiniAppSignInResult{
		UserID:      authdata.UserID,
		AccessToken: authdata.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   authdata.ExpiresAt.Format(time.RFC3339),
	}}, nil
}

func (h *AuthHandler) MaxSignIn(ctx context.Context, input *dto.MiniAppSignInInput) (*dto.MiniAppSignInOutput, error) {
	idn := &authmodel.Identity{
		ProviderName: authmodel.ProviderMax,
		AppInitData:  input.Body.AppInitData,
		Metadata:     input.Body.MetaData,
	}

	var metadata *authmodel.Metadata
	if input.Body.MetaData != nil {
		metadata = authmodel.NewUserMetadata(*input.Body.MetaData)
	}

	authdata, err := h.appAuthHandler.Handle(ctx, idn, metadata)
	if err != nil {
		return nil, err
	}

	return &dto.MiniAppSignInOutput{Body: dto.MiniAppSignInResult{
		UserID:      authdata.UserID,
		AccessToken: authdata.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   authdata.ExpiresAt.Format(time.RFC3339),
	}}, nil
}

func (h *AuthHandler) ServiceSignIn(ctx context.Context, input *dto.MessengerSignInInput) (*dto.MessengerSignInOutput, error) {

	idn := &authmodel.Identity{
		ProviderName:   authmodel.ProviderService,
		ProviderUserID: input.Body.ProviderID,
		ServiceKey:     input.InternalKey,
		Metadata:       input.Body.MetaData,
	}

	usr := &usermodel.User{}
	if input.Body.FullName != nil {
		usr.FullName = *input.Body.FullName
	}

	var metadata *authmodel.Metadata
	if input.Body.MetaData != nil {
		metadata = authmodel.NewUserMetadata(*input.Body.MetaData)
		usr.OriginSource = metadata.UTMData.Campaign
	}

	authdata, err := h.externalAuthHandler.Handle(ctx, idn, usr, metadata)
	if err != nil {
		return nil, err
	}

	return &dto.MessengerSignInOutput{Body: dto.MessengerSignInResult{
		UserID:      authdata.UserID,
		AccessToken: authdata.AccessToken,
		TokenType:   "Bearer",
		ExpiresAt:   authdata.ExpiresAt.Format(time.RFC3339),
	}}, nil
}
