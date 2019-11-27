package handler

var Path = []string{"commons", "app", "handler", "handler.go"}

var Content = `package handler

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"{{ .AppRepository }}/commons/app/view"
	"{{ .AppRepository }}/commons/log"
)

func ApiRoot(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("handler: ApiRoot")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "Yeah! Your API is working!"))
}

func Error401(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: Error401")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	json.NewEncoder(w).Encode(view.SetUnauthorizedErrorMessage())
}

func Error403(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: Error403")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	json.NewEncoder(w).Encode(view.SetForbiddenErrorMessage())
}

func Error404(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: Error404")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(view.SetNotFoundErrorMessage())
}

func QueryParamsToMapCriteria(param string, mapParams map[string][]string) map[string]string {
	var criteria, value string
	var checkParam, removePrefix, removeSufix *regexp.Regexp
	var err error

	query := make(map[string]string)

	checkParam, err = regexp.Compile(param + ` + "`" + `\[[a-zA-Z0-9\-\_]+\](\[\]){0,1}` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	removePrefix, err = regexp.Compile(param + ` + "`" + `\[` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	removeSufix, err = regexp.Compile(` + "`" + `\](\[\]){0,1}` + "`" + `)
	if err != nil {
		log.Warn.Println(err)
	}

	for key := range mapParams {
		if checkParam.MatchString(key) {
			criteria = key
			criteria = removeSufix.ReplaceAllString(criteria, "")
			criteria = removePrefix.ReplaceAllString(criteria, "")
			value = strings.Join(mapParams[key], ",")

			query[criteria] = value
		}
	}

	return query
}`
