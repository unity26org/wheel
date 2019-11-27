package newapp

import (
	"github.com/unity26org/wheel/generator/gencommon"
	"github.com/unity26org/wheel/templates/templateapp"
	"github.com/unity26org/wheel/templates/templateapp/app/handlers"
	"github.com/unity26org/wheel/templates/templateapp/app/myself"
	"github.com/unity26org/wheel/templates/templateapp/app/session"
	"github.com/unity26org/wheel/templates/templateapp/app/session/sessionmailer"
	"github.com/unity26org/wheel/templates/templateapp/app/usertemplate"
	"github.com/unity26org/wheel/templates/templateapp/commons/app/handler"
	"github.com/unity26org/wheel/templates/templateapp/commons/app/model"
	"github.com/unity26org/wheel/templates/templateapp/commons/app/view"
	"github.com/unity26org/wheel/templates/templateapp/commons/conversor"
	"github.com/unity26org/wheel/templates/templateapp/commons/crypto"
	"github.com/unity26org/wheel/templates/templateapp/commons/locale"
	"github.com/unity26org/wheel/templates/templateapp/commons/logtemplate"
	"github.com/unity26org/wheel/templates/templateapp/commons/mailer"
	"github.com/unity26org/wheel/templates/templateapp/config"
	"github.com/unity26org/wheel/templates/templateapp/config/configlocales"
	"github.com/unity26org/wheel/templates/templateapp/db/entities"
	"github.com/unity26org/wheel/templates/templateapp/db/schema"
	"github.com/unity26org/wheel/templates/templateapp/routes"
)

var templateVar gencommon.TemplateVar
var rootAppPath string

func prependRootAppPathToPath(path []string) []string {
	return append([]string{rootAppPath}, path...)
}

func Generate(options map[string]interface{}) {
	// Main vars
	templateVar = gencommon.TemplateVar{
		AppName:       options["app_name"].(string),
		AppRepository: options["app_repository"].(string),
		SecretKey:     gencommon.SecureRandom(128),
	}
	rootAppPath = gencommon.BuildRootAppPath(options["app_repository"].(string))

	// APP Root path
	gencommon.CreateRootAppPath(rootAppPath)

	// APP Handlers
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.MyselfPath), handlers.MyselfContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.SessionPath), handlers.SessionContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.UserPath), handlers.UserContent, templateVar)

	// APP myself
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(myself.ViewPath), myself.ViewContent, templateVar)

	// APP session
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(session.ModelPath), session.ModelContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(session.ViewPath), session.ViewContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.PasswordRecoveryEnPath), sessionmailer.PasswordRecoveryEnContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.PasswordRecoveryPtBrPath), sessionmailer.PasswordRecoveryPtBrContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.SignUpEnPath), sessionmailer.SignUpEnContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.SignUpPtBrPath), sessionmailer.SignUpPtBrContent, templateVar)

	// APP user
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(usertemplate.ModelPath), usertemplate.ModelContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(usertemplate.ViewPath), usertemplate.ViewContent, templateVar)

	// COMMONS APPs
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handler.Path), handler.Content, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.Path), model.Content, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.OrderingPath), model.OrderingContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.PaginationPath), model.PaginationContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.SearchEnginePath), model.SearchEngineContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(view.Path), view.Content, templateVar)

	// COMMONS conversor
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(conversor.Path), conversor.Content, templateVar)

	// COMMONS crypto
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(crypto.Path), crypto.Content, templateVar)

	// COMMONS locale
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(locale.Path), locale.Content, templateVar)

	// COMMONS log
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(logtemplate.Path), logtemplate.Content, templateVar)

	// COMMONS mailer
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(mailer.Path), mailer.Content, templateVar)

	// config
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.Path), config.Content, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.AppPath), config.AppContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.DatabasePath), config.DatabaseContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(config.EmailPath), config.EmailContent, templateVar)

	// config certs
	gencommon.GenerateCertificates(rootAppPath)

	// config locales
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(configlocales.EnPath), configlocales.EnContent, templateVar)
	gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(configlocales.PtBrPath), configlocales.PtBrContent, templateVar)

	// db
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(entities.SessionPath), entities.SessionContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(entities.UserPath), entities.UserContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(schema.Path), schema.Content, templateVar)

	// routes
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.AuthorizePath), routes.AuthorizeContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.MiddlewarePath), routes.MiddlewareContent, templateVar)
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.Path), routes.Content, templateVar)

	// main
	gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(templates.MainPath), templates.MainContent, templateVar)
	if options["git_ignore"].(bool) {
		gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(templates.GitIgnorePath), templates.GitIgnoreContent, templateVar)
	}

	// Final
	gencommon.NotifyNewApp(rootAppPath)
}
