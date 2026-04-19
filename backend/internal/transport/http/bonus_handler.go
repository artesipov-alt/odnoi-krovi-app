package http

import (
	"bytes"
	"context"
	"io"
	"net/http"

	bonuscmd "github.com/artesipov-alt/odnoi-krovi-app/internal/application/bonus/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
	"github.com/danielgtaylor/huma/v2"
)

// BonusHandler handles HTTP requests for bonuses.
type BonusHandler struct {
	importHandler *bonuscmd.ImportBonusesHandler
}

// NewBonusHandler creates a new BonusHandler.
func NewBonusHandler(importHandler *bonuscmd.ImportBonusesHandler) *BonusHandler {
	return &BonusHandler{
		importHandler: importHandler,
	}
}

// Register registers the bonus routes.
func (h *BonusHandler) Register(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "import-bonuses",
		Method:      http.MethodPost,
		Path:        "/v1/admin/bonuses/import",
		Summary:     "Импорт бонусов из Excel",
		Description: "Загружает бонусы из Excel файла. Требуются права администратора.",
		Tags:        []string{"admin-v1"},
	}, h.ImportBonuses)
}

// ImportBonuses handles the import of bonuses from an uploaded Excel file.
func (h *BonusHandler) ImportBonuses(ctx context.Context, input *dto.ImportBonusesInput) (*dto.ImportBonusesOutput, error) {
	formData := input.RawBody.Data()

	fileContent, err := io.ReadAll(formData.File)
	if err != nil {
		return nil, huma.Error400BadRequest("Не удалось прочитать файл")
	}

	result, err := h.importHandler.Handle(ctx, bytes.NewReader(fileContent))
	if err != nil {
		return nil, err
	}

	return &dto.ImportBonusesOutput{
		Body: dto.ImportBonusesResult{
			TotalRows: result.TotalRows,
			Imported:  result.Imported,
			Skipped:   result.Skipped,
			Errors:    result.Errors,
		},
	}, nil
}
