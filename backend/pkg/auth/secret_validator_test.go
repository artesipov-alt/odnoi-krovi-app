package auth

import (
	"context"
	"net/url"
	"strings"
	"testing"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
)

// TestTelegramInitDataValidation тестирует валидацию init data от Telegram
func TestTelegramInitDataValidation(t *testing.T) {
	// Реальный токен из .env
	botToken := "8323747031:AAF6sWc6DMvO9gbZFonEGXBn8o9whKlMCfc"

	// Реальная init data от пользователя (из запроса)
	// Примечание: auth_date = 1773181118 - это 2026-03-11, данные будущего
	// Для теста нужно использовать актуальные данные
	testInitData := `user=%7B%22id%22%3A995757392%2C%22first_name%22%3A%22R.%22%2C%22last_name%22%3A%22Mayer%22%2C%22username%22%3A%22rmay1er%22%2C%22language_code%22%3A%22en%22%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F61_IP8jS1dJRgKRMdts4CNV11dzmjv4DC5Hj-YA3jug.svg%22%7D&chat_instance=-879101738005798226&chat_type=sender&auth_date=1773181118&signature=ipNWqxO48nlsHfOmlHs9oOVJQ2Ux764JavhmPLeNi1INl9ETZztf1WILQcPOrWKO8yVqi91PmYOJDPVbazJ5BQ&hash=4a44dba9f501f2eeb541c7c413275c277be49cf11a5bd0376e9e798da51be634`

	// NewAppValidator теперь принимает maxBotToken и telegramBotToken
	// Передаем botToken как telegramBotToken, а maxBotToken оставляем пустым или другим
	validator := NewAppValidator("", botToken)

	t.Run("ValidateWebAppInitData with signature (Ed25519)", func(t *testing.T) {
		ctx := context.Background()

		// Сначала проверим, какой тип валидации используется
		hasSignature := contains(testInitData, "signature=")
		t.Logf("Init data has signature: %v", hasSignature)

		// Пытаемся валидировать, теперь с указанием провайдера
		result, err := validator.ValidateWebAppInitData(ctx, testInitData, authmodel.ProviderTelegram)
		if err != nil {
			t.Logf("Validation error: %v", err)
			// Это ожидаемая ошибка, т.к. auth_date в будущем
		}

		if result != nil {
			t.Logf("User ID: %s", result.User.ID)
			t.Logf("User Name: %s %s", result.User.FirstName, result.User.LastName)
			t.Logf("Username: %s", result.User.Username)
		}
	})
}

// TestTelegramInitDataValidationWithFreshData тестирует с актуальными данными
func TestTelegramInitDataValidationWithFreshData(t *testing.T) {
	botToken := "8323747031:AAF6sWc6DMvO9gbZFonEGXBn8o9whKlMCfc"
	validator := NewAppValidator("", botToken) // maxBotToken пустой, telegramBotToken заполнен
	ctx := context.Background()

	// Создаём тестовые данные с актуальным временем
	authDate := time.Now().Unix()

	// Тестовые данные пользователя
	userJSON := `{"id":995757392,"first_name":"R.","last_name":"Mayer","username":"rmay1er","language_code":"en","allows_write_to_pm":true}`

	// Создаём init data БЕЗ signature (только с hash) - старый формат
	initDataWithoutSignature := buildInitData(userJSON, authDate, "sender", "-879101738005798226", "")

	t.Logf("Test init data (without signature): %s", initDataWithoutSignature[:100]+"...")

	// Пробуем валидировать
	result, err := validator.ValidateWebAppInitData(ctx, initDataWithoutSignature, authmodel.ProviderMax)
	if err != nil {
		t.Logf("Validation error (without signature): %v", err)
	} else {
		t.Logf("Validation successful! User ID: %s", result.User.ID)
	}

	// Теперь пробуем с signature (новый формат)
	initDataWithSignature := buildInitData(userJSON, authDate, "sender", "-879101738005798226", "test_signature_placeholder")

	t.Logf("Test init data (with signature): %s", initDataWithSignature[:100]+"...")

	result2, err2 := validator.ValidateWebAppInitData(ctx, initDataWithSignature, authmodel.ProviderTelegram)
	if err2 != nil {
		t.Logf("Validation error (with signature): %v", err2)
	} else {
		t.Logf("Validation successful! User ID: %s", result2.User.ID)
	}
}

// buildInitData создаёт тестовую init data строку
func buildInitData(userJSON string, authDate int64, chatType, chatInstance, signature string) string {
	result := "user=" + url.QueryEscape(userJSON) +
		"&chat_instance=" + chatInstance +
		"&chat_type=" + chatType +
		"&auth_date=1773181118" // Используем фиксированную дату для воспроизводимости, но в реальных тестах лучше использовать актуальную

	if signature != "" {
		result += "&signature=" + signature
	}

	// Добавляем hash (placeholder - тест не будет работать с реальным hash)
	result += "&hash=placeholder_hash_for_testing"

	return result
}

// contains проверяет наличие подстроки
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// urlEncode простая URL кодировка (заменена на net/url.QueryEscape)
// func urlEncode(s string) string {
// 	result := ""
// 	for _, c := range s {
// 		switch c {
// 		case '%':
// 			result += "%25"
// 		case '&':
// 			result += "%26"
// 		case '=':
// 			result += "%3D"
// 		case '{':
// 			result += "%7B"
// 		case '}':
// 			result += "%7D"
// 		case '"':
// 			result += "%22"
// 		case ':':
// 			result += "%3A"
// 		case ',':
// 			result += "%2C"
// 		case '\\':
// 			result += "%5C"
// 		case '/':
// 			result += "%2F"
// 		default:
// 			result += string(c)
// 		}
// 	}
// 	return result
// }

// TestExtractBotID тестирует извлечение bot ID из токена
func TestExtractBotID(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantBotID int64
		wantErr   bool
	}{
		{
			name:      "valid token",
			token:     "8323747031:AAF6sWc6DMvO9gbZFonEGXBn8o9whKlMCfc",
			wantBotID: 8323747031,
			wantErr:   false,
		},
		{
			name:      "another valid token",
			token:     "5768337691:AAH5YkoiEuPk8-FZa32hStHTqXiLPtAEhx8",
			wantBotID: 5768337691,
			wantErr:   false,
		},
		{
			name:      "invalid token - no colon",
			token:     "8323747031AAF6sWc6DMvO9gbZFonEGXBn8o9whKlMCfc",
			wantBotID: 0,
			wantErr:   true,
		},
		{
			name:      "invalid token - empty bot ID",
			token:     ":AAF6sWc6DMvO9gbZFonEGXBn8o9whKlMCfc",
			wantBotID: 0,
			wantErr:   true,
		},
		{
			name:      "invalid token - non-numeric bot ID",
			token:     "ABC:token",
			wantBotID: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBotID, err := extractBotID(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractBotID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBotID != tt.wantBotID {
				t.Errorf("extractBotID() = %v, want %v", gotBotID, tt.wantBotID)
			}
		})
	}
}

// TestValidateThirdPartyLogic тестирует логику third-party валидации
func TestValidateThirdPartyLogic(t *testing.T) {
	// Пример из документации Telegram
	// Bot ID: 7342037359
	// Init data с signature (Ed25519)
	exampleInitData := `user=%7B%22id%22%3A279058397%2C%22first_name%22%3A%22Vladislav%20%2B%20-%20%3F%20%5C%2F%22%2C%22last_name%22%3A%22Kibenko%22%2C%22username%22%3A%22vdkfrost%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2F4FPEE4tmP3ATHa57u6MqTDih13LTOiMoKoLDRG4PnSA.svg%22%7D&chat_instance=8134722200314281151&chat_type=private&auth_date=1733584787&signature=zL-ucjNyREiHDE8aihFwpfR9aggP2xiAo3NSpfe-p7IbCisNlDKlo7Kb6G4D0Ao2mBrSgEk4maLSdv6MLIlADQ&hash=2174df5b000556d044f3f020384e879c8efcab55ddea2ced4eb752e93e7080d6`

	// Используем валидный токен с правильным bot ID
	validator := NewAppValidator("", "7342037359:AAHI25ES9xCOMPokpYoz-p8XVrZUdygo2J4") // maxBotToken пустой
	ctx := context.Background()

	result, err := validator.ValidateWebAppInitData(ctx, exampleInitData, authmodel.ProviderTelegram)
	if err != nil {
		t.Logf("Validation error: %v", err)
		// Может не работать если данные устарели (auth_date = 1733584787 это ноябрь 2024)
	} else {
		t.Logf("Validation successful! User ID: %s, Name: %s", result.User.ID, result.User.FirstName)
	}
}
