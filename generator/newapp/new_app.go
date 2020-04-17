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
	"github.com/unity26org/wheel/templates/templateapp/db/migrate"
	"github.com/unity26org/wheel/templates/templateapp/db/migrate/adapter"
	"github.com/unity26org/wheel/templates/templateapp/db/migrate/adapter/postgresql"
	"github.com/unity26org/wheel/templates/templateapp/db/schema"
	"github.com/unity26org/wheel/templates/templateapp/db/schema/data/col"
	"github.com/unity26org/wheel/templates/templateapp/routes"
	"time"
)

var templateVar gencommon.TemplateVar
var rootAppPath string

func prependRootAppPathToPath(path []string) []string {
	return append([]string{rootAppPath}, path...)
}

func generateHandlers() error {
	err := gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.MyselfPath), handlers.MyselfContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.SessionPath), handlers.SessionContent, templateVar)
	if err != nil {
		return err
	}
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handlers.UserPath), handlers.UserContent, templateVar)
	if err != nil {
		return err
	}

	return nil
}

func generateSession() error {
	err := gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(session.ModelPath), session.ModelContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(session.ViewPath), session.ViewContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.PasswordRecoveryEnPath), sessionmailer.PasswordRecoveryEnContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.PasswordRecoveryPtBrPath), sessionmailer.PasswordRecoveryPtBrContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.SignUpEnPath), sessionmailer.SignUpEnContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(sessionmailer.SignUpPtBrPath), sessionmailer.SignUpPtBrContent, templateVar)
	if err != nil {
		return err
	}

	return nil
}

func generateUser() error {
	err := gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(usertemplate.ModelPath), usertemplate.ModelContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(usertemplate.ViewPath), usertemplate.ViewContent, templateVar)
	if err != nil {
		return err
	}

	return nil
}

func generateCommmonsApp() error {
	err := gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(handler.Path), handler.Content, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.Path), model.Content, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.OrderingPath), model.OrderingContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.PaginationPath), model.PaginationContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(model.SearchEnginePath), model.SearchEngineContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(view.Path), view.Content, templateVar)
	if err != nil {
		return err
	}

	return nil
}

func generateCommonsOthers() error {
	// COMMONS conversor
	err := gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(conversor.Path), conversor.Content, templateVar)
	if err != nil {
		return err
	}

	// COMMONS crypto
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(crypto.Path), crypto.Content, templateVar)
	if err != nil {
		return err
	}

	// COMMONS locale
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(locale.Path), locale.Content, templateVar)
	if err != nil {
		return err
	}

	// COMMONS log
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(logtemplate.Path), logtemplate.Content, templateVar)
	if err != nil {
		return err
	}

	// COMMONS mailer
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(mailer.Path), mailer.Content, templateVar)
	if err != nil {
		return err
	}

	return nil
}

func generateConfig() error {
	err := gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.Path), config.Content, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.AppPath), config.AppContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(config.DatabasePath), config.DatabaseContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(config.EmailPath), config.EmailContent, templateVar)
	if err != nil {
		return err
	}

	// config certs
	err = gencommon.GenerateCertificates(rootAppPath)
	if err != nil {
		return err
	}

	// config locales
	err = gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(configlocales.EnPath), configlocales.EnContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.CreatePathAndFileFromTemplateString(prependRootAppPathToPath(configlocales.PtBrPath), configlocales.PtBrContent, templateVar)
	if err != nil {
		return err
	}

	return nil
}

func generateDb() error {
	err := gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(entities.SessionPath), entities.SessionContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(entities.UserPath), entities.UserContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(schema.Path), schema.Content, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(col.Path), col.Content, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(migrate.FacadePath), migrate.FacadeContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(adapter.Path), adapter.Content, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(postgresql.Path), postgresql.Content, templateVar)
	if err != nil {
		return err
	}

	templateVar.MigrationMetadata = gencommon.MigrationMetadata{Type: "CREATE_TABLE", Name: "CreateUsers", Version: time.Now().Format("20060102150405")}
	migrate.UserPath[len(migrate.UserPath)-1] = templateVar.MigrationMetadata.Version + "_" + migrate.UserPath[len(migrate.UserPath)-1]
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(migrate.UserPath), migrate.UserContent, templateVar)
	if err != nil {
		return err
	}
	templateVar.EntityName = gencommon.SetEntityName("user")
	gencommon.UpdateMigrate(rootAppPath, templateVar)

	for {
		if templateVar.MigrationMetadata.Version != time.Now().Format("20060102150405") {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	templateVar.MigrationMetadata = gencommon.MigrationMetadata{Type: "CREATE_TABLE", Name: "CreateSessions", Version: time.Now().Format("20060102150405")}
	migrate.SessionPath[len(migrate.SessionPath)-1] = templateVar.MigrationMetadata.Version + "_" + migrate.SessionPath[len(migrate.SessionPath)-1]
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(migrate.SessionPath), migrate.SessionContent, templateVar)
	if err != nil {
		return err
	}
	templateVar.EntityName = gencommon.SetEntityName("session")
	gencommon.UpdateMigrate(rootAppPath, templateVar)

	return nil
}

func generateRoutes() error {
	err := gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.AuthorizePath), routes.AuthorizeContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.MiddlewarePath), routes.MiddlewareContent, templateVar)
	if err != nil {
		return err
	}

	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(routes.Path), routes.Content, templateVar)
	if err != nil {
		return err
	}

	return nil
}

func Generate(options map[string]interface{}) error {
	var err error

	// Main vars
	templateVar = gencommon.TemplateVar{
		AppName:       options["app_name"].(string),
		AppRepository: options["app_repository"].(string),
		SecretKey:     gencommon.SecureRandom(128),
	}

	rootAppPath, err = gencommon.BuildRootAppPath(options["app_repository"].(string))
	if err != nil {
		return err
	}

	// APP Root path
	if err = gencommon.CreateRootAppPath(rootAppPath); err != nil {
		return err
	}

	// APP handler
	if err = generateHandlers(); err != nil {
		return err
	}

	// APP myself
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(myself.ViewPath), myself.ViewContent, templateVar)
	if err != nil {
		return err
	}

	// APP session
	if err = generateSession(); err != nil {
		return err
	}

	// APP user
	if err = generateUser(); err != nil {
		return err
	}

	// COMMONS APPs
	if err = generateCommmonsApp(); err != nil {
		return err
	}

	// COMMONS Others
	if err = generateCommonsOthers(); err != nil {
		return err
	}

	// config
	if err = generateConfig(); err != nil {
		return err
	}

	// db
	if err = generateDb(); err != nil {
		return err
	}

	// routes
	if err = generateRoutes(); err != nil {
		return err
	}

	// main
	err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(templates.MainPath), templates.MainContent, templateVar)
	if err != nil {
		return err
	}

	if options["git_ignore"].(bool) {
		err = gencommon.GeneratePathAndFileFromTemplateString(prependRootAppPathToPath(templates.GitIgnorePath), templates.GitIgnoreContent, templateVar)
		if err != nil {
			return err
		}
	}

	// Final
	gencommon.NotifyNewApp(rootAppPath)

	return nil
}
