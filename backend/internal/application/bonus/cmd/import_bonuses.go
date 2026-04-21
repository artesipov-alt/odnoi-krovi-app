package cmd

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	bonusrepo "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/xuri/excelize/v2"
)

// ImportBonusesHandler handles the business logic for importing bonuses from an Excel file.
type ImportBonusesHandler struct {
	repo bonusrepo.Repository
}

// NewImportBonusesHandler creates a new instance of ImportBonusesHandler.
func NewImportBonusesHandler(repo bonusrepo.Repository) *ImportBonusesHandler {
	return &ImportBonusesHandler{
		repo: repo,
	}
}

// ImportResult contains statistics about the import operation.
type ImportResult struct {
	TotalRows int
	Imported  int
	Skipped   int
	Errors    []string
}

// Handle processes the Excel file and imports bonuses into the database.
func (h *ImportBonusesHandler) Handle(ctx context.Context, file io.Reader) (*ImportResult, error) {
	result := &ImportResult{}

	// Read the Excel file
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, apperrors.BadRequest(fmt.Sprintf("Не удалось прочитать Excel файл: %v", err))
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, apperrors.Internal(err, "Не удалось получить строки из файла")
	}

	if len(rows) < 2 {
		return result, nil // No data rows
	}

	result.TotalRows = len(rows) - 1

	// Collect all promo codes to check for duplicates in DB
	var promoCodes []string
	for _, row := range rows[1:] {
		if len(row) > 6 {
			promoCodes = append(promoCodes, strings.TrimSpace(row[6]))
		}
	}

	existingCodes, err := h.repo.ExistsByPromoCodes(ctx, promoCodes)
	if err != nil {
		return nil, apperrors.Internal(err, "Ошибка проверки существующих промокодов")
	}

	existingMap := make(map[string]struct{})
	for _, code := range existingCodes {
		existingMap[code] = struct{}{}
	}

	var bonusesToCreate []*bonusmodel.Bonus

	// Skip header row
	for i, row := range rows[1:] {
		rowNum := i + 2
		b, err := h.parseRow(row)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Строка %d: %v", rowNum, err))
			result.Skipped++
			continue
		}

		if _, exists := existingMap[b.PromoCode]; exists {
			result.Skipped++
			continue
		}

		bonusesToCreate = append(bonusesToCreate, b)
	}

	if len(bonusesToCreate) > 0 {
		if err := h.repo.CreateBatch(ctx, bonusesToCreate); err != nil {
			return nil, apperrors.Internal(err, "Ошибка сохранения бонусов в базу данных")
		}
		result.Imported = len(bonusesToCreate)
	}

	return result, nil
}

// parseRow parses a single row from the Excel file and returns a bonus domain model.
func (h *ImportBonusesHandler) parseRow(row []string) (*bonusmodel.Bonus, error) {
	// Expected columns:
	// 0: ID (ignored)
	// 1: Partner Name
	// 2: Description
	// 3: Target (cat/dog/all)
	// 4: Recipient (donor/recipient/all)
	// 5: Category (food/preparation/other)
	// 6: Promo Code
	// 7: Expires At
	// 8: Platform Name
	// 9: Platform URL

	if len(row) < 9 {
		return nil, fmt.Errorf("недостаточно колонок")
	}

	partnerName := strings.TrimSpace(row[1])
	if partnerName == "" {
		return nil, fmt.Errorf("наименование партнера обязательно")
	}

	description := strings.TrimSpace(row[2])
	if description == "" {
		return nil, fmt.Errorf("описание бонуса обязательно")
	}

	target, err := mapTarget(strings.TrimSpace(row[3]))
	if err != nil {
		return nil, fmt.Errorf("ошибка в поле 'Для кого': %v", err)
	}

	recipient, err := mapRecipient(strings.TrimSpace(row[4]))
	if err != nil {
		return nil, fmt.Errorf("ошибка в поле 'Для кого (донор/реципиент)': %v", err)
	}

	category, err := mapCategory(strings.TrimSpace(row[5]))
	if err != nil {
		return nil, fmt.Errorf("ошибка в поле 'Категория': %v", err)
	}

	promoCode := strings.TrimSpace(row[6])
	if promoCode == "" {
		return nil, fmt.Errorf("промокод обязателен")
	}

	expiresAt, err := parseDate(strings.TrimSpace(row[7]))
	if err != nil {
		return nil, fmt.Errorf("ошибка в дате окончания: %v", err)
	}

	platformName := strings.TrimSpace(row[8])
	if platformName == "" {
		return nil, fmt.Errorf("площадка применения обязательна")
	}

	var platformURL *string
	if len(row) > 9 {
		url := strings.TrimSpace(row[9])
		if url != "" {
			platformURL = &url
		}
	}

	return &bonusmodel.Bonus{
		PartnerName:  partnerName,
		Description:  description,
		Target:       target,
		Recipient:    recipient,
		Category:     category,
		PromoCode:    promoCode,
		ExpiresAt:    expiresAt,
		PlatformName: platformName,
		PlatformURL:  platformURL,
		IsActive:     true,
	}, nil
}

func mapTarget(val string) (string, error) {
	switch strings.ToLower(val) {
	case "все", "all":
		return "all", nil
	case "кошка", "cat":
		return "cat", nil
	case "собака", "dog":
		return "dog", nil
	default:
		return "", fmt.Errorf("неизвестное значение '%s', допустимые: Все, Кошка, Собака", val)
	}
}

func mapRecipient(val string) (string, error) {
	switch strings.ToLower(val) {
	case "все", "all":
		return "all", nil
	case "донор", "donor":
		return "donor", nil
	case "реципиент", "recipient":
		return "recipient", nil
	default:
		return "", fmt.Errorf("неизвестное значение '%s', допустимые: Все, Донор, Реципиент", val)
	}
}

func mapCategory(val string) (string, error) {
	switch strings.ToLower(val) {
	case "корма", "food":
		return "food", nil
	case "препараты", "preparation":
		return "preparation", nil
	case "другое", "other":
		return "other", nil
	default:
		return "", fmt.Errorf("неизвестное значение '%s', допустимые: Корма, Препараты, Другое", val)
	}
}

func parseDate(val string) (time.Time, error) {
	if val == "" {
		return time.Time{}, fmt.Errorf("дата обязательна")
	}

	// Handle "Бессрочно" (no expiry) by returning a far future date
	if strings.ToLower(val) == "бессрочно" {
		return time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC), nil
	}

	// Try DD.MM.YYYY format first
	t, err := time.Parse("02.01.2006", val)
	if err == nil {
		return t, nil
	}

	// Try YYYY-MM-DD format
	t, err = time.Parse("2006-01-02", val)
	if err == nil {
		return t, nil
	}

	// Try MM-DD-YY format (e.g., 09-26-26)
	t, err = time.Parse("01-02-06", val)
	if err == nil {
		return t, nil
	}

	// Try parsing as float (Excel date serial number)
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		// Excel dates start from 1899-12-30
		t := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC).Add(time.Duration(f * float64(24*time.Hour)))
		return t, nil
	}

	return time.Time{}, fmt.Errorf("неподдерживаемый формат даты '%s', используйте DD.MM.YYYY", val)
}
