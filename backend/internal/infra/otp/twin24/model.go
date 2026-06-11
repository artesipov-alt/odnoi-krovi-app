package twin24

const (
	// VarTwinCode - имя переменной для кода подтверждения в сценарии бота
	VarTwinCode = "twin_code"
)

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
