package templatecrud

var HandlerContent = `package handlers

import (
	"encoding/json"
	"{{ .AppRepository }}/app/{{ .EntityName.LowerCase }}"
	"{{ .AppRepository }}/commons/app/handler"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/commons/app/view"
	"{{ .AppRepository }}/commons/log"
	"{{ .AppRepository }}/db/entities"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
)

func {{ .EntityName.CamelCase }}Create(w http.ResponseWriter, r *http.Request) {
	var new{{ .EntityName.CamelCase }} = entities.{{ .EntityName.CamelCase }}{}

	log.Info.Println("Handler: {{ .EntityName.CamelCase }}Create")
	w.Header().Set("Content-Type", "application/json")

	{{ .EntityName.LowerCamelCase }}SetParams(&new{{ .EntityName.CamelCase }}, r)

	valid, errs := {{ .EntityName.LowerCase }}.Create(&new{{ .EntityName.CamelCase }})

	if valid {
		json.NewEncoder(w).Encode({{ .EntityName.LowerCase }}.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "{{ .EntityName.SnakeCase }} was successfully created"), {{ .EntityName.CamelCase }}: {{ .EntityName.LowerCase }}.SetJson(new{{ .EntityName.CamelCase }})})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "{{ .EntityName.SnakeCase }} was not created", errs))
	}
}

func {{ .EntityName.CamelCase }}Update(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: {{ .EntityName.CamelCase }}Update")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	{{ .EntityName.LowerCamelCase }}Current, err := {{ .EntityName.LowerCase }}.Find(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "{{ .EntityName.SnakeCase }} was not updated", []error{err}))
		return
	}

	{{ .EntityName.LowerCamelCase }}SetParams(&{{ .EntityName.LowerCamelCase }}Current, r)

	if valid, errs := {{ .EntityName.LowerCase }}.Update(&{{ .EntityName.LowerCamelCase }}Current); valid {
		json.NewEncoder(w).Encode({{ .EntityName.LowerCase }}.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "{{ .EntityName.SnakeCase }} was successfully updated"), {{ .EntityName.CamelCase }}: {{ .EntityName.LowerCase }}.SetJson({{ .EntityName.LowerCamelCase }}Current)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "{{ .EntityName.SnakeCase }} was not updated", errs))
	}
}

func {{ .EntityName.CamelCase }}Destroy(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: {{ .EntityName.CamelCase }}Destroy")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	{{ .EntityName.LowerCamelCase }}Current, err := {{ .EntityName.LowerCase }}.Find(params["id"])

	if err == nil && {{ .EntityName.LowerCase }}.Destroy(&{{ .EntityName.LowerCamelCase }}Current) {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "{{ .EntityName.SnakeCase }} was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("alert", "{{ .EntityName.SnakeCase }} was not found"))
	}
}

func {{ .EntityName.CamelCase }}Show(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: {{ .EntityName.CamelCase }}Show")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	{{ .EntityName.LowerCamelCase }}Current, err := {{ .EntityName.LowerCase }}.Find(params["id"])

	if err == nil {
		json.NewEncoder(w).Encode({{ .EntityName.LowerCase }}.SetJson({{ .EntityName.LowerCamelCase }}Current))
	} else {
		json.NewEncoder(w).Encode(view.SetSystemMessage("alert", "{{ .EntityName.SnakeCase }} was not found"))
	}
}

func {{ .EntityName.CamelCase }}List(w http.ResponseWriter, r *http.Request) {
	var i, page, entries, pages int
	var {{ .EntityName.LowerCamelCase }}List []entities.{{ .EntityName.CamelCase }}

	{{ .EntityName.LowerCamelCase }}Jsons := []{{ .EntityName.LowerCase }}.Json{}

	log.Info.Println("Handler: {{ .EntityName.CamelCase }}List")
	w.Header().Set("Content-Type", "application/json")

	criteria := handler.QueryParamsToMapCriteria("search", r.URL.Query())
	order := {{ .EntityName.LowerCamelCase }}SanitizeOrder(r.FormValue("order"))

	{{ .EntityName.LowerCamelCase }}List, page, pages, entries = {{ .EntityName.LowerCase }}.Paginate(criteria, order, r.FormValue("page"), r.FormValue("per_page"))

	for i = 0; i < len({{ .EntityName.LowerCamelCase }}List); i++ {
		{{ .EntityName.LowerCamelCase }}Jsons = append({{ .EntityName.LowerCamelCase }}Jsons, {{ .EntityName.LowerCase }}.SetJson({{ .EntityName.LowerCamelCase }}List[i]))
	}

	pagination := view.MainPagination{CurrentPage: page, TotalPages: pages, TotalEntries: entries}
	json.NewEncoder(w).Encode({{ .EntityName.LowerCase }}.PaginationJson{Pagination: pagination, {{ .EntityName.CamelCasePlural }}: {{ .EntityName.LowerCamelCase }}Jsons})
}

func {{ .EntityName.LowerCamelCase }}SetParams({{ .EntityName.LowerCamelCase }}Set *entities.{{ .EntityName.CamelCase }}, r *http.Request) {
  {{- $filteredEntityColumns := filterEntityColumnsNotForeignKeys .EntityColumns }}
	var allowedParams = []string{ {{- range $index, $element := $filteredEntityColumns }} "{{ $element.NameSnakeCase }}" {{- if isNotLastIndex $index $filteredEntityColumns }}, {{- end }} {{- end }} }

	r.ParseMultipartForm(100 * 1024)

	for key := range r.Form {
		for _, allowedParam := range allowedParams {
			if key == allowedParam {
				model.SetColumnValue({{ .EntityName.LowerCamelCase }}Set, allowedParam, r.FormValue(allowedParam))
				break
			}
		}
	}
}

func {{ .EntityName.LowerCamelCase }}SanitizeOrder(value string) string {
	var allowedParams = []*regexp.Regexp{
    regexp.MustCompile(` + "`" + `(?i)\A\s*id(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
    {{- $filteredEntityColumns := filterEntityColumnsNotForeignKeys .EntityColumns }}
  	{{- range $index, $element := $filteredEntityColumns }} 
    regexp.MustCompile(` + "`" + `(?i)\A\s*{{ $element.NameSnakeCase }}(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `), {{- end }}
    regexp.MustCompile(` + "`" + `(?i)\A\s*created_at(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
    regexp.MustCompile(` + "`" + `(?i)\A\s*updated_at(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `)}

	for _, allowedParam := range allowedParams {
		if allowedParam.MatchString(value) {
			return value
		}
	}

	return ""
}`
