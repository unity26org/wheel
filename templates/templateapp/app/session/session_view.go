package session

var ViewPath = []string{"app", "session", "session_view.go"}

var ViewContent = `package session

import (
	"bytes"
	"html/template"
	"{{ .AppRepository }}/app/user"
	"{{ .AppRepository }}/commons/app/view"
	"{{ .AppRepository }}/commons/log"
	"{{ .AppRepository }}/config"
	"{{ .AppRepository }}/db/entities"
)

type SignInSuccess struct {
	Message view.SystemMessage ` + "`" + `json:"system_message"` + "`" + `
	Token   string             ` + "`" + `json:"token"` + "`" + `
	Expires int                ` + "`" + `json:"expires"` + "`" + `
}

type SignOutSuccess struct {
	Message view.SystemMessage ` + "`" + `json:"system_message"` + "`" + `
}

type SignUpSuccess struct {
	UserFirstName string
	AppName       string
}

type PasswordRecoveryInstructions struct {
	UserFirstName          string
	LinkToPasswordRecovery string
}

func SignInSuccessMessage(mType string, content string, token string) SignInSuccess {
	return SignInSuccess{Message: view.SystemMessage{mType, content}, Token: token, Expires: config.App.TokenExpirationSeconds}
}

func SignOutSuccessMessage(mType string, content string) SignOutSuccess {
	return SignOutSuccess{Message: view.SystemMessage{mType, content}}
}

func RefreshSuccessMessage(mType string, content string, token string) SignInSuccess {
	return SignInSuccessMessage(mType, content, token)
}

func SignUpSuccessMessage(mType string, content string, token string) SignInSuccess {
	return SignInSuccessMessage(mType, content, token)
}

func SignUpMailer(currentUser *entities.User) string {
	var content bytes.Buffer

	data := SignUpSuccess{UserFirstName: user.FirstName(currentUser), AppName: config.App.AppName}

	tmpl, err := template.ParseFiles("./app/session/mailer/sign_up." + currentUser.Locale + ".html")
	if err != nil {
		log.Error.Println(err)
	}

	err = tmpl.Execute(&content, &data)

	return content.String()
}

func PasswordRecoveryInstructionsMailer(currentUser *entities.User, token string) string {
	var content bytes.Buffer

	data := PasswordRecoveryInstructions{UserFirstName: user.FirstName(currentUser), LinkToPasswordRecovery: config.App.ResetPasswordUrl + "?token=" + token}

	tmpl, err := template.ParseFiles("./app/session/mailer/password_recovery." + currentUser.Locale + ".html")
	if err != nil {
		log.Error.Println(err)
	}

	err = tmpl.Execute(&content, &data)

	return content.String()
}`
