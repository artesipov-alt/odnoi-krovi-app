package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	errorWebhookTimeout = 5 * time.Second
	// Вместимость очереди алертов: при шквале 5xx лишние алерты
	// отбрасываются, чтобы не завалить вебхук и не задерживать ответы.
	errorWebhookQueueSize = 10
	// Сколько байт тела ответа попадает в алерт.
	errorWebhookBodyLimit = 1024
)

// statusRecorder запоминает статус-код и сниппет тела, записанного обработчиком.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	body        bytes.Buffer
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status = code
		r.wroteHeader = true
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.status = http.StatusOK
		r.wroteHeader = true
	}
	if r.body.Len() < errorWebhookBodyLimit {
		n := min(len(b), errorWebhookBodyLimit-r.body.Len())
		r.body.Write(b[:n])
	}
	return r.ResponseWriter.Write(b)
}

// errorAlert — payload, отправляемый на вебхук.
type errorAlert struct {
	Env     string `json:"env"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Status  int    `json:"status"`
	Error   string `json:"error,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
	Time    string `json:"time"`
}

// ErrorWebhookMiddleware отправляет алерт на внешний вебхук, если ответ получил статус 5xx.
// Должен стоять первым в цепочке (вне sloghttp.Recovery), чтобы ловить
// в том числе 500, записанные Recovery при панике.
// В поле error попадает тело ответа (например, detail из JSON-ошибки Huma),
// а при панике — её текст. Пустой url полностью выключает middleware.
func ErrorWebhookMiddleware(url, env string) func(http.Handler) http.Handler {
	if url == "" {
		return func(next http.Handler) http.Handler { return next }
	}

	client := &http.Client{Timeout: errorWebhookTimeout}
	queue := make(chan errorAlert, errorWebhookQueueSize)

	// Один воркер на всё время жизни процесса: HTTP-запрос на вебхук
	// никогда не блокирует обработку запроса клиента.
	go func() {
		for alert := range queue {
			sendErrorAlert(client, url, alert)
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			var panicMsg string
			alertSent := false

			sendAlert := func(status int, desc string) {
				if alertSent {
					return
				}
				alertSent = true
				select {
				case queue <- errorAlert{
					Env:     env,
					Method:  r.Method,
					Path:    r.URL.Path,
					Status:  status,
					Error:   desc,
					TraceID: rec.Header().Get("X-Trace-Id"),
					Time:    time.Now().Format(time.RFC3339),
				}:
				default:
					slog.Warn("ErrorWebhook: очередь переполнена, алерт пропущен", "path", r.URL.Path)
				}
			}

			// Отправляем алерт в том числе при панике (defer сработает до выхода).
			defer func() {
				if panicMsg != "" {
					// Recovery, который стоит ниже по цепочке, запишет 500.
					sendAlert(http.StatusInternalServerError, "panic: "+panicMsg)
					return
				}
				if rec.status >= 500 {
					sendAlert(rec.status, strings.TrimSpace(rec.body.String()))
				}
			}()

			// Перехватываем панику только чтобы запомнить её текст,
			// и пробрасываем дальше — sloghttp.Recovery отработает как раньше.
			func() {
				defer func() {
					if v := recover(); v != nil && v != http.ErrAbortHandler {
						panicMsg = fmt.Sprint(v)
						panic(v)
					}
				}()
				next.ServeHTTP(rec, r)
			}()
		})
	}
}

func sendErrorAlert(client *http.Client, url string, alert errorAlert) {
	body, err := json.Marshal(alert)
	if err != nil {
		slog.Warn("ErrorWebhook: не удалось сериализовать алерт", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), errorWebhookTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		slog.Warn("ErrorWebhook: не удалось создать запрос", "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		slog.Warn("ErrorWebhook: не удалось отправить алерт", "error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		slog.Warn("ErrorWebhook: вебхук вернул неуспешный статус", "status", resp.StatusCode)
	}
}
