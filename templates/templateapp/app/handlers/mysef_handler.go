package handlers

var MyselfPath = []string{"app", "handlers", "myself_handler.go"}

var MyselfContent = `package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"{{ .AppRepository }}/app/models/myself"
	"{{ .AppRepository }}/app/models/user"
	"{{ .AppRepository }}/commons/app/handler"
	"{{ .AppRepository }}/commons/app/view"
	"{{ .AppRepository }}/commons/log"
)

func MyselfUpdate(w http.ResponseWriter, r *http.Request) {
  type MyselfPermittedParams struct {
  	Name     string ` + "`" + `json:"name"` + "`" + `
  	Locale   string ` + "`" + `json:"locale"` + "`" + `
  }

	log.Info.Println("Handler: MyselfUpdate")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current

	var myselfParams MyselfPermittedParams
	err := json.NewDecoder(r.Body).Decode(&myselfParams)
  if err != nil {
		log.Error.Println("could not parser input JSON")
		w.WriteHeader(500)
    json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "could not parser input JSON", []error{err}))
    return
  }

	handler.SetPermittedParamsToEntity(&myselfParams, &userMyself)

	if valid, errs := user.Save(&userMyself); valid {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully updated"))
	} else {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "user was not updated", errs))
	}
}

func MyselfUpdatePassword(w http.ResponseWriter, r *http.Request) {
  type MyselfPasswordParams struct {
  	Password             string ` + "`" + `json:"password"` + "`" + `
    PasswordConfirmation string ` + "`" + `json:"password_confirmation"` + "`" + `
  }

	var errs []error
	var valid bool

	log.Info.Println("Handler: MyselfChangePassword")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current

	var myselfParams MyselfPasswordParams
	err := json.NewDecoder(r.Body).Decode(&myselfParams)
  if err != nil {
		log.Error.Println("could not parser input JSON")
		w.WriteHeader(500)
    json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "could not parser input JSON", []error{err}))
    return
  }

  userMyself.Password = myselfParams.Password

	if !user.Exists(&userMyself) {
		errs = append(errs, errors.New("invalid user"))
	} else if myselfParams.Password != myselfParams.PasswordConfirmation {
		errs = append(errs, errors.New("password confirmation does not match new password"))
	} else if valid, errs = user.Save(&userMyself); valid {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "password was successfully changed"))
	}

	if len(errs) > 0 {
		json.NewEncoder(w).Encode(view.SetErrorMessage("alert", "password could not be changed", errs))
	}
}

func MyselfDestroy(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: MyselfDestroy")
	w.Header().Set("Content-Type", "application/json")

	userMyself := user.Current

	if user.Destroy(&userMyself) {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("notice", "user was successfully destroyed"))
	} else {
		json.NewEncoder(w).Encode(view.SetDefaultMessage("alert", "user could not be destroyed"))
	}
}

func MyselfShow(w http.ResponseWriter, r *http.Request) {
	log.Info.Println("Handler: MyselfShow")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(myself.SetJson(user.Current))
}`
