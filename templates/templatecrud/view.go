package templatecrud

var ViewContent = `package {{ .EntityName.LowerCase }}

import (
	"time"
	"{{ .AppRepository }}/app/entities"
	"{{ .AppRepository }}/commons/app/view"
)

type PaginationJson struct {
	Pagination view.MainPagination ` + "`" + `json:"pagination"` + "`" + `
	{{ .EntityName.CamelCasePlural }} []Json ` + "`" + `json:"{{ .EntityName.SnakeCasePlural }}"` + "`" + `
}

type SuccessfullySavedJson struct {
	SystemMessage view.SystemMessage ` + "`" + `json:"system_message"` + "`" + `
	{{ .EntityName.CamelCase }} Json ` + "`" + `json:"{{ .EntityName.SnakeCase }}"` + "`" + `
}

type Json struct {
	ID uint64 ` + "`" + `json:"id"` + "`" + `
  {{- range .EntityColumns }}
  {{- if not .IsForeignKey }}
  {{ .Name }} {{ .Type }} ` + "`" + `json:"{{ .NameSnakeCase }}"` + "`" + `
  {{- end }}  
  {{- end }}
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func SetJson({{ .EntityName.LowerCamelCase }} entities.{{ .EntityName.CamelCase }}) Json {
	return Json{
		ID: {{ .EntityName.LowerCamelCase }}.ID,
    {{- range .EntityColumns }}
    {{- if not .IsForeignKey }}
    {{ .Name }}: {{ $.EntityName.LowerCamelCase }}.{{ .Name }},
    {{- end }}
    {{- end }}
		CreatedAt: {{ .EntityName.LowerCamelCase }}.CreatedAt,
		UpdatedAt: {{ .EntityName.LowerCamelCase }}.UpdatedAt,
	}
}`
