package templatecrud

var CustomHandlerContent = `package handlers

import (
	"encoding/json"
	"{{ .AppRepository }}/commons/app/view"
	"{{ .AppRepository }}/commons/log"
	"net/http"
)

{{- range $element := .EntityColumns }}

func {{ $.EntityName.CamelCase }}{{ $element.Name }}(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: {{ $.EntityName.CamelCase }}{{ $element.Name }}")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "Handler: {{ $.EntityName.CamelCase }}{{ $element.Name }}"))
}

{{- end }}
`
