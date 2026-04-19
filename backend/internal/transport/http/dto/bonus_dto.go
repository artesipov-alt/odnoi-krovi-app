package dto

import (
	"github.com/danielgtaylor/huma/v2"
)

// ImportBonusesInput represents the request input for importing bonuses from an Excel file.
type ImportBonusesInput struct {
	RawBody huma.MultipartFormFiles[struct {
		File huma.FormFile `form:"file" required:"true"`
	}]
}

// ImportBonusesOutput represents the response output for the bonus import operation.
type ImportBonusesOutput struct {
	Body ImportBonusesResult
}

// ImportBonusesResult contains statistics about the import operation.
type ImportBonusesResult struct {
	TotalRows int      `json:"totalRows" doc:"Total number of data rows processed"`
	Imported  int      `json:"imported" doc:"Number of bonuses successfully imported"`
	Skipped   int      `json:"skipped" doc:"Number of rows skipped (duplicates or errors)"`
	Errors    []string `json:"errors,omitempty" doc:"List of error messages for skipped rows"`
}
