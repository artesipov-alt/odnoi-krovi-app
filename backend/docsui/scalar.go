package docsui

import (
	_ "embed"
	"net/http"
)

//go:embed scalar.html
var scalarHTML []byte

// ScalarDocsHandler serves the Scalar Docs documentation UI.
func ScalarDocsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(scalarHTML)
}
