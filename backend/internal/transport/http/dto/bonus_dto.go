package dto

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// ImportBonusesInput представляет входные данные запроса для импорта бонусов из файла Excel.
type ImportBonusesInput struct {
	RawBody huma.MultipartFormFiles[struct {
		File huma.FormFile `form:"file" required:"true"`
	}]
}

// ImportBonusesOutput представляет выходные данные ответа для операции импорта бонусов.
type ImportBonusesOutput struct {
	Body ImportBonusesResult
}

// ImportBonusesResult содержит статистику об операции импорта.
type ImportBonusesResult struct {
	TotalRows int      `json:"totalRows" doc:"Общее количество обработанных строк данных"`
	Imported  int      `json:"imported" doc:"Количество успешно импортированных бонусов"`
	Skipped   int      `json:"skipped" doc:"Количество пропущенных строк (дубликаты или ошибки)"`
	Errors    []string `json:"errors,omitempty" doc:"Список сообщений об ошибках для пропущенных строк"`
}

type AssignedBonusItem struct {
	ID           string     `json:"id" doc:"Уникальный идентификатор присвоенного бонуса"`
	UserID       *string    `json:"userId,omitempty" doc:"ID пользователя, связанного с присвоенным бонусом"`
	PartnerName  string     `json:"partnerName" doc:"Название партнера, предоставляющего присвоенный бонус"`
	Description  string     `json:"description" doc:"Описание присвоенного бонуса"`
	Target       string     `json:"target" doc:"Целевая аудитория для присвоенного бонуса" enum:"dog,cat,all"`
	Recipient    string     `json:"recipient" doc:"Получатель присвоенного бонуса"`
	Category     string     `json:"category" doc:"Категория присвоенного бонуса" enum:"food,preparation,other,lock"`
	Subcategory  *string    `json:"subcategory,omitempty" doc:"Подкатегория присвоенного бонуса"`
	PromoCode    string     `json:"promoCode" doc:"Промокод для присвоенного бонуса"`
	ExpiresAt    time.Time  `json:"expiresAt" doc:"Дата и время истечения срока присвоенного бонуса"`
	PlatformName string     `json:"platformName" doc:"Название платформы, где можно использовать присвоенный бонус"`
	PlatformURL  *string    `json:"platformUrl,omitempty" doc:"URL платформы для присвоенного бонуса"`
	Stage        string     `json:"stage" doc:"Текущая стадия присвоенного бонуса" enum:"unused,used"`
	CreatedAt    *time.Time `json:"createdAt,omitempty" doc:"Временная метка создания присвоенного бонуса"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty" doc:"Временная метка последнего обновления присвоенного бонуса"`
	AssignedAt   *time.Time `json:"assignedAt,omitempty" doc:"Временная метка, когда бонус был присвоен"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty" doc:"Временная метка удаления присвоенного бонуса, если мягко удалено"`
}

// AssignedBonus представляет присвоенные бонусы, распределенные по категориям.
type AssignedBonus struct {
	TotalPriority int                 `json:"priority" doc:"Количество приоритетов"`
	Food          []AssignedBonusItem `json:"food" doc:"Бонусы категории food"`
	Preparation   []AssignedBonusItem `json:"preparation" doc:"Бонусы категории preparation"`
	Other         []AssignedBonusItem `json:"other" doc:"Бонусы категории other"`
}

// AssignedBonusOutput представляет выходные данные ответа для присвоенных бонусов.
type AssignedBonusOutput struct {
	Body AssignedBonus
}
