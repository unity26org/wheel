package handler

var Path = []string{"commons", "app", "handler", "handler.go"}

var Content = `package handler

import (
	"encoding/json"
	"{{ .AppRepository }}/commons/app/view"
	"{{ .AppRepository }}/commons/log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
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
	w.WriteHeader(404)

	json.NewEncoder(w).Encode(view.SetNotFoundErrorMessage())
}

func Error400(w http.ResponseWriter, r *http.Request, jsonParseError bool) {
	log.Info.Println("Handler: Error400")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	if jsonParseError {
		json.NewEncoder(w).Encode(view.SetBadRequestInvalidJsonErrorMessage())
	} else {
		json.NewEncoder(w).Encode(view.SetBadRequestErrorMessage())
	}
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
}

func SetPermittedParamsToEntity(params interface{}, entity interface{}) {
	setPermittedParams(params, entity, []string{}, []string{})
}

func SetPermittedParamsToEntityWithExceptions(params interface{}, entity interface{}, excepts []string) {
	setPermittedParams(params, entity, excepts, []string{})
}

func SetPermittedParamsToEntityButOnly(params interface{}, entity interface{}, only []string) {
	setPermittedParams(params, entity, []string{}, only)
}

func setPermittedParams(params interface{}, entity interface{}, excepts []string, only []string) {
	valParams := reflect.ValueOf(params).Elem()
	valEntity := reflect.ValueOf(entity).Elem()
	returnEntity := reflect.ValueOf(entity)

	for i := 0; i < valParams.NumField(); i++ {
		paramValueField := valParams.Field(i)
		paramTypeField := valParams.Type().Field(i)

		if !inExcept(excepts, paramTypeField.Name) && inOnly(only, paramTypeField.Name) {
			for j := 0; j < valEntity.NumField(); j++ {
				entityTypeField := valEntity.Type().Field(j)
				if paramTypeField.Name == entityTypeField.Name {
					returnEntity.Elem().Field(j).Set(paramValueField)
				}
			}
		}
	}
}

func inExcept(excepts []string, value string) bool {
	if len(excepts) == 0 {
		return false
	} else {
		for _, element := range excepts {
			if value == element {
				return true
			}
		}

		return false
	}
}

func inOnly(only []string, value string) bool {
	if len(only) == 0 {
		return true
	} else {
		for _, element := range only {
			if value == element {
				return true
			}
		}

		return false
	}
}`
