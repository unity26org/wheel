package handlers

var UserPath = []string{"app", "handlers", "user_handler.go"}

var UserContent = `package handlers

import (
	"encoding/json"
	"{{ .AppRepository }}/app/user"
	"{{ .AppRepository }}/commons/app/handler"
	"{{ .AppRepository }}/commons/app/view"
	"{{ .AppRepository }}/commons/log"
	"{{ .AppRepository }}/db/entities"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
  type UserPermittedParams struct {
  	Name     string ` + "`" + `json:"name"` + "`" + `
  	Email    string ` + "`" + `json:"email"` + "`" + `
  	Password string ` + "`" + `json:"password"` + "`" + `
  	Locale   string ` + "`" + `json:"locale"` + "`" + `
  	Admin    bool   ` + "`" + `json:"admin"` + "`" + `
  }

	var userNew = entities.User{}

	log.Info.Println("Handler: UserCreate")
	w.Header().Set("Content-Type", "application/json")

	var userParams UserPermittedParams
	_ = json.NewDecoder(r.Body).Decode(&userParams)

	handler.SetPermittedParamsToEntity(&userParams, &userNew)

	valid, errs := user.Create(&userNew)

	if valid {
		json.NewEncoder(w).Encode(user.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "user was successfully created"), User: user.SetJson(userNew)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not created", errs))
	}
}

func UserUpdate(w http.ResponseWriter, r *http.Request) {
  type UserPermittedParams struct {
  	Name     string ` + "`" + `json:"name"` + "`" + `
  	Email    string ` + "`" + `json:"email"` + "`" + `
  	Locale   string ` + "`" + `json:"locale"` + "`" + `
  	Admin    bool   ` + "`" + `json:"admin"` + "`" + `
  }

	log.Info.Println("Handler: UserUpdate")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	userCurrent, err := user.Find(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not updated", []error{err}))
		return
	}

	var userParams UserPermittedParams
	_ = json.NewDecoder(r.Body).Decode(&userParams)

	handler.SetPermittedParamsToEntity(&userParams, &userCurrent)

	if valid, errs := user.Update(&userCurrent); valid {
		json.NewEncoder(w).Encode(user.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "user was successfully updated"), User: user.SetJson(userCurrent)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not updated", errs))
	}
}

func UserUpdatePassword(w http.ResponseWriter, r *http.Request) {
  type UserPermittedParams struct {
  	Password             string ` + "`" + `json:"password"` + "`" + `
  }

	log.Info.Println("Handler: UserUpdate")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	userCurrent, err := user.Find(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user password was not updated", []error{err}))
		return
	}

	var userParams UserPermittedParams
	_ = json.NewDecoder(r.Body).Decode(&userParams)

	handler.SetPermittedParamsToEntity(&userParams, &userCurrent)

	if valid, errs := user.Update(&userCurrent); valid {
		json.NewEncoder(w).Encode(user.SuccessfullySavedJson{SystemMessage: view.SetSystemMessage("notice", "user password was successfully updated"), User: user.SetJson(userCurrent)})
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user password was not updated", errs))
	}
}


func UserDestroy(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: UserDestroy")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userCurrent, err := user.Find(params["id"])

	if err == nil && user.Destroy(&userCurrent) {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("alert", "user was not found"))
	}
}

func UserShow(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: UserShow")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userCurrent, err := user.Find(params["id"])

	if err == nil {
		json.NewEncoder(w).Encode(user.SetJson(userCurrent))
	} else {
		json.NewEncoder(w).Encode(view.SetSystemMessage("alert", "user was not found"))
	}
}

func UserList(w http.ResponseWriter, r *http.Request) {
	var i, page, entries, pages int
	var userList []entities.User
  
	userJsons := []user.Json{}

	log.Info.Println("Handler: UserList")
	w.Header().Set("Content-Type", "application/json")

	criteria := handler.QueryParamsToMapCriteria("search", r.URL.Query())
	order := userSanitizeOrder(r.FormValue("order"))

	userList, page, pages, entries = user.Paginate(criteria, order, r.FormValue("page"), r.FormValue("per_page"))

	for i = 0; i < len(userList); i++ {
		userJsons = append(userJsons, user.SetJson(userList[i]))
	}

	pagination := view.MainPagination{CurrentPage: page, TotalPages: pages, TotalEntries: entries}
	json.NewEncoder(w).Encode(user.PaginationJson{Pagination: pagination, Users: userJsons})
}

func userSanitizeOrder(value string) string {
	var allowedParams = []*regexp.Regexp{
		regexp.MustCompile(` + "`" + `(?i)\A\s*id(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*name(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*email(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*password(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*admin(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*locale(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*created_at(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*updated_at(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `),
		regexp.MustCompile(` + "`" + `(?i)\A\s*deleted_at(\s+(DESC|ASC)){0,1}\s*\z` + "`" + `)}

	for _, allowedParam := range allowedParams {
		if allowedParam.MatchString(value) {
			return value
		}
	}

	return ""
}`
