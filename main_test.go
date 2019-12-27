package main

import (
	// "fmt"
	"github.com/unity26org/wheel/commons/diff"
	"github.com/unity26org/wheel/commons/fileutil"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type WheelFileSystem struct {
	Type    string
	Path    []string
	Content *string
}

func TestMainHelp(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"wheel", "--help"}
	main()

	w.Close()

	out, _ := ioutil.ReadAll(r)

	os.Stdout = rescueStdout

	if strings.Trim(helpContent, " \n") != strings.Trim(string(out), " \n") {
		t.Errorf("Help error. \nShould \"%s\" .\nGot: \"%s\"", strings.Trim(helpContent, " \n"), strings.Trim(string(out), " \n"))
	}
}

func TestMainVersion(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"wheel", "--version"}
	main()

	w.Close()

	out, _ := ioutil.ReadAll(r)

	os.Stdout = rescueStdout

	if strings.Trim(versionContent, " \n") != strings.Trim(string(out), " \n") {
		t.Errorf("Version error. \nShould \"%s\" .\nGot: \"%s\"", strings.Trim(helpContent, " \n"), strings.Trim(string(out), " \n"))
	}
}

func TestMainInvalidArgument(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"wheel"}
	main()

	w.Close()

	out, _ := ioutil.ReadAll(r)

	os.Stdout = rescueStdout

	if !strings.Contains(string(out), "invalid argument") {
		t.Errorf("Not printing \"invalid argument\"")
	}

	if !strings.Contains(string(out), "error") {
		t.Errorf("Not printing \"error\"")
	}

	if !strings.Contains(string(out), strings.Trim(helpContent, " \n")) {
		t.Errorf("Print help wrong")
	}
}

func TestMain(t *testing.T) {
	var wheelFileSystems []WheelFileSystem

	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"app"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"app", "handlers"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appHandlersMyself, Type: "FILE", Path: []string{"app", "handlers", "myself_handler.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appHandlersSession, Type: "FILE", Path: []string{"app", "handlers", "session_handler.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appHandlersUser, Type: "FILE", Path: []string{"app", "handlers", "user_handler.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"app", "myself"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appMyselfView, Type: "FILE", Path: []string{"app", "myself", "myself_view.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"app", "session"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appSessionModel, Type: "FILE", Path: []string{"app", "session", "session_model.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appSessionView, Type: "FILE", Path: []string{"app", "session", "session_view.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"app", "session", "mailer"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appSessionMailerPasswordRecoveryEn, Type: "FILE", Path: []string{"app", "session", "mailer", "password_recovery.en.html"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appSessionMailerPasswordRecoveryPtBr, Type: "FILE", Path: []string{"app", "session", "mailer", "password_recovery.pt-BR.html"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appSessionMailerSignUpEn, Type: "FILE", Path: []string{"app", "session", "mailer", "sign_up.en.html"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appSessionMailerSignUpPtBr, Type: "FILE", Path: []string{"app", "session", "mailer", "sign_up.pt-BR.html"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"app", "user"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appUserModelContent, Type: "FILE", Path: []string{"app", "user", "user_model.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &appUserViewContent, Type: "FILE", Path: []string{"app", "user", "user_view.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "app"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "app", "handler"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsAppHandleContent, Type: "FILE", Path: []string{"commons", "app", "handler", "handler.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "app", "model"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsAppModelContent, Type: "FILE", Path: []string{"commons", "app", "model", "model.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsAppOrderingContent, Type: "FILE", Path: []string{"commons", "app", "model", "ordering.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsAppPaginationContent, Type: "FILE", Path: []string{"commons", "app", "model", "pagination.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsAppSearchEngineContent, Type: "FILE", Path: []string{"commons", "app", "model", "searchengine.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "app", "view"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsAppViewContent, Type: "FILE", Path: []string{"commons", "app", "view", "view.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "conversor"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsConversorContent, Type: "FILE", Path: []string{"commons", "conversor", "conversor.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "crypto"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsCryptoContent, Type: "FILE", Path: []string{"commons", "crypto", "crypto.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "locale"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsLocaleContent, Type: "FILE", Path: []string{"commons", "locale", "locale.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "log"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &commonsLogContent, Type: "FILE", Path: []string{"commons", "log", "log.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"commons", "mailer"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &mailerContent, Type: "FILE", Path: []string{"commons", "mailer", "mailer.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"config"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &configAppContent, Type: "FILE", Path: []string{"config", "app.yml"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &configConfigContent, Type: "FILE", Path: []string{"config", "config.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &configDatabaseContent, Type: "FILE", Path: []string{"config", "database.yml"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &configEmailContent, Type: "FILE", Path: []string{"config", "email.yml"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"config", "keys"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "FILE", Path: []string{"config", "keys", "app.key.rsa"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "FILE", Path: []string{"config", "keys", "app.key.rsa.pub"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"config", "locales"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &localeEnContent, Type: "FILE", Path: []string{"config", "locales", "en.yml"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &localePtBrContent, Type: "FILE", Path: []string{"config", "locales", "pt-BR.yml"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"db"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"db", "entities"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &dbEntitySessionContent, Type: "FILE", Path: []string{"db", "entities", "session_entity.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &dbEntityUserContent, Type: "FILE", Path: []string{"db", "entities", "user_entity.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"db", "schema"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &dbSchemaMigrateContentV1, Type: "FILE", Path: []string{"db", "schema", "migrate.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Type: "DIR", Path: []string{"routes"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &authorizeContentV1, Type: "FILE", Path: []string{"routes", "authorize.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &middlewareContent, Type: "FILE", Path: []string{"routes", "middleware.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &routesContentV1, Type: "FILE", Path: []string{"routes", "routes.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &mainContent, Type: "FILE", Path: []string{"main.go"}})
	wheelFileSystems = append(wheelFileSystems, WheelFileSystem{Content: &gitIgnoreContent, Type: "FILE", Path: []string{".gitignore"}})

	currentDir, err := os.Getwd()
	if err != nil {
		t.Errorf("could not get current directory")
	}

	currentUser, err := user.Current()
	if err != nil {
		t.Errorf("could not get current user")
	}

	goPath := filepath.Join(currentUser.HomeDir, "go", "src")

	if !fileutil.DirOrFileExists(goPath) {
		t.Errorf("go path not found. should be %s", goPath)
	}

	destroyTestRepositoryDir(filepath.Join(goPath, "test_repository_hub.com"))

	testRepositoryPath := filepath.Join(goPath, "test_repository_hub.com")
	testAccountPath := filepath.Join(testRepositoryPath, "test_account")
	testAppPath := filepath.Join(testAccountPath, "test_project")

	os.Stdout, _ = os.Open(os.DevNull)
	os.Args = []string{"wheel", "new", "test_repository_hub.com/test_account/test_project"}
	main()

	if !fileutil.DirOrFileExists(testRepositoryPath) {
		t.Errorf("directory %s was not found", testRepositoryPath)
	}

	if !fileutil.DirOrFileExists(testAccountPath) {
		t.Errorf("directory %s was not found", testAccountPath)
	}

	if !fileutil.DirOrFileExists(testAppPath) {
		t.Errorf("directory %s was not found", testAppPath)
	}

	for _, fileDir := range wheelFileSystems {
		currentPath := testAppPath
		for _, path := range fileDir.Path {
			currentPath = filepath.Join(currentPath, path)
			if !fileutil.DirOrFileExists(currentPath) {
				t.Errorf("directory %s was not found", currentPath)
			}
		}
	}

	replaceRandForConstFromConfigApp(testAppPath, wheelFileSystems)

	for _, fileDir := range wheelFileSystems {
		if fileDir.Type == "FILE" && fileDir.Content != nil {
			basePath, fileName := buildPath(testAppPath, fileDir.Path)
			fileContent := strings.Trim(fileutil.ReadTextFile(basePath, fileName), " \n")
			varContent := strings.Trim(*fileDir.Content, " \n")

			if fileContent != varContent {
				diffError := diff.Diff(varContent, basePath, fileName, true)
				t.Errorf("generated file %s is diffent from expected\nDIFF:\n\n%s", filepath.Join(basePath, fileName), diffError)
			}
		}
	}

	os.Chdir(filepath.Join(goPath, "test_repository_hub.com", "test_account", "test_project"))
	os.Stdout, _ = os.Open(os.DevNull)
	os.Args = []string{"wheel", "generate", "scaffold", "post", "title:string", "description:text", "user:reference", "rate:decimal", "views:integer", "likes:uint", "published:boolean", "published_at:datetime"}
	main()

	os.Chdir(currentDir)

	// destroyTestRepositoryDir(filepath.Join(goPath, "test_repository_hub.com"))
}

func destroyTestRepositoryDir(path string) {
	if fileutil.DirOrFileExists(path) {
		fileutil.DestroyAllDirOrFile(path)
	}
}

func buildPath(basePath string, filePath []string) (string, string) {
	if len(filePath) == 1 {
		return basePath, filePath[0]
	} else {
		fileName := filePath[len(filePath)-1]
		filePath = filePath[:len(filePath)-1]

		for _, path := range filePath {
			basePath = filepath.Join(basePath, path)
		}

		return basePath, fileName
	}
}

func replaceRandForConstFromConfigApp(basePath string, wheelFileSystems []WheelFileSystem) {
	var appConfig WheelFileSystem

	for _, fileDir := range wheelFileSystems {
		pathFile := strings.Join(fileDir.Path, "_")
		if pathFile == "config_app.yml" {
			appConfig = fileDir
			break
		}
	}

	basePath, fileName := buildPath(basePath, appConfig.Path)
	fileContent := fileutil.ReadTextFile(basePath, fileName)

	secretKeyRegexp := regexp.MustCompile(`secret\_key\:\s*\"[A-Fa-f0-9]{128}\"\s*`)
	fileContent = secretKeyRegexp.ReplaceAllString(fileContent, "secret_key: \"0B7f3892773a0be1E4C6c992f9D4BcdB96Bf2277665133d8CdE18FfED0ECcF7d4a5e7114d4c5fDD3430ad9d87A1C705aACfBCBBD928B38248aBd48Ff8bDdA5E0\"\n")

	fileutil.SaveTextFile(fileContent, basePath, fileName)
}
