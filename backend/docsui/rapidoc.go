package docsui

import (
	"net/http"
)

// RapiDocHandler serves the RapiDoc documentation UI.
func RapiDocHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<!doctype html>
<html>
<head>
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width, minimum-scale=1, initial-scale=1, user-scalable=yes">
				<title>RapiDoc</title>
				<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/rapidoc@latest/dist/rapidoc-min.css">
</head>
<body>
				<rapi-doc
								id="rapidoc-container"
								theme="dark"
								layout="row"
								render-style="focused"
								show-header="false"
								allow-authentication="true"
								allow-spec-url-load="true"
								allow-spec-file-load="true"
				>
				</rapi-doc>
				<script type="module" src="https://cdn.jsdelivr.net/npm/rapidoc@latest/dist/rapidoc-min.js"></script>
				<script>
								document.addEventListener('DOMContentLoaded', function() {
												const rapidoc = document.getElementById('rapidoc-container');
												if (rapidoc) {
																rapidoc.setAttribute('spec-url', window.location.origin + '/openapi.json');
												}
								});
				</script>
</body>
</html>`))
}
