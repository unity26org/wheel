package gencommon

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/unity26org/wheel/commons/diff"
	"github.com/unity26org/wheel/commons/fileutil"
	"github.com/unity26org/wheel/commons/notify"
	"github.com/unity26org/wheel/generator/newmigrate"
	"github.com/unity26org/wheel/templates/templatecommon"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type AppConfig struct {
	AppName                        string   `yaml:"app_name"`
	AppRepository                  string   `yaml:"app_repository"`
	SecretKey                      string   `yaml:"secret_key"`
	ResetPasswordExpirationSeconds int      `yaml:"reset_password_expiration_seconds"`
	ResetPasswordUrl               string   `yaml:"reset_password_url"`
	TokenExpirationSeconds         int      `yaml:"token_expiration_seconds"`
	Locales                        []string `yaml:"locales"`
}

type EntityColumn struct {
	Name                string
	NameSnakeCase       string
	NameSnakeCasePlural string
	Type                string
	Extras              string
	IsRelation          bool
	IsForeignKey        bool
	MigrateType         string
	MigrateExtra        string
}

type EntityName struct {
	CamelCase            string
	CamelCasePlural      string
	LowerCamelCase       string
	LowerCamelCasePlural string
	SnakeCase            string
	SnakeCasePlural      string
	LowerCase            string
}

type MigrationMetadata struct {
	Type          string
	Name          string
	Version       string
	FileNameSufix string
	Entity        string
}

type TemplateVar struct {
	AppRepository     string
	AppName           string
	SecretKey         string
	EntityName        EntityName
	EntityColumns     []EntityColumn
	MigrationMetadata MigrationMetadata
}

var yesToAll = false

func FmtNewContent(content string) (string, error) {
	var fileName string

	tmpfile, err := ioutil.TempFile("", "wheel-new-content*.go")
	if err != nil {
		return "", err
	}

	fileName = strings.TrimPrefix(tmpfile.Name(), os.TempDir())
	fileName = strings.TrimPrefix(fileName, `/`)
	fileName = strings.TrimPrefix(fileName, `\`)

	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	if err != nil {
		return "", err
	}

	err = tmpfile.Close()
	if err != nil {
		return "", err
	}

	fileutil.GoFmtFile(tmpfile.Name())

	return fileutil.ReadTextFile(os.TempDir(), fileName)
}

func GetAppConfig() (AppConfig, error) {
	var appConfig AppConfig

	b, err := fileutil.ReadBytesFile(filepath.Join(".", "config"), "app.yml")
	if err != nil {
		return AppConfig{}, err
	}

	err = yaml.Unmarshal(b, &appConfig)
	if err != nil {
		return AppConfig{}, err
	}

	return appConfig, nil
}

func GenerateFromTemplateString(content string, templateVar TemplateVar) (string, error) {
	var buffContent bytes.Buffer

	FuncMap := template.FuncMap{
		// TODO: Filter for References and Not References
		"hasDateTimeType": func(tEntityColumns []EntityColumn) bool {
			for _, element := range tEntityColumns {
				if element.Type == "*time.Time" {
					return true
				}
			}
			return false
		},
		"isLastIndex": func(index int, tSlice interface{}) bool {
			return index == reflect.ValueOf(tSlice).Len()-1
		},
		"isNotLastIndex": func(index int, tSlice interface{}) bool {
			return index != reflect.ValueOf(tSlice).Len()-1
		},
		"filterEntityColumnsNotForeignKeys": func(tEntityColumns []EntityColumn) []EntityColumn {
			var notForeignKeys []EntityColumn
			for _, element := range tEntityColumns {
				if !element.IsForeignKey {
					notForeignKeys = append(notForeignKeys, element)
				}
			}
			return notForeignKeys
		},
		"filterEntityColumnsForeignKeysOnly": func(tEntityColumns []EntityColumn) []EntityColumn {
			var foreignKeys []EntityColumn
			for _, element := range tEntityColumns {
				if element.IsForeignKey {
					foreignKeys = append(foreignKeys, element)
				}
			}
			return foreignKeys
		},
		"filterEntityColumnsRelationOnly": func(tEntityColumns []EntityColumn) []EntityColumn {
			var relations []EntityColumn
			for _, element := range tEntityColumns {
				if element.IsRelation {
					relations = append(relations, element)
				}
			}
			return relations
		},
		"filterEntityColumnsNotRelations": func(tEntityColumns []EntityColumn) []EntityColumn {
			var notRelations []EntityColumn
			for _, element := range tEntityColumns {
				if !element.IsRelation {
					notRelations = append(notRelations, element)
				}
			}
			return notRelations
		},
		"checkMigrationType": func(migrationType string) bool {
			return templateVar.MigrationMetadata.Type == migrationType
		},
	}

	tmpl, err := template.New("T").Funcs(FuncMap).Parse(content)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&buffContent, templateVar)
	if err != nil {
		return "", err
	}

	return buffContent.String(), nil
}

func GenerateFromTemplateFile(templatePath string, templateVar TemplateVar) (string, error) {
	var content bytes.Buffer

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&content, &templateVar)
	if err != nil {
		return "", err
	}

	return content.String(), nil
}

func overwriteFile(content string, filePath string, fileName string) string {
	var pseudoMode string

	fullPath := filepath.Join(filePath, fileName)

	if yesToAll {
		return "f"
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		notify.Simple("Overwrite " + fullPath + "? (enter \"h\" for help) [Ynaqdph] ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		switch text {
		case "Y":
			pseudoMode = "f"
		case "n":
			pseudoMode = "s"
		case "a":
			yesToAll = true
			pseudoMode = "f"
		case "q":
			notify.Fatal("Aborting...")
		case "d":
			diff.Print(content, filePath, fileName)
		case "p":
			diff.Patch(content, filePath, fileName)
			pseudoMode = "p"
		default:
			notify.Simple(overwriteFileHelp())
			pseudoMode = ""
		}

		if pseudoMode == "f" || pseudoMode == "s" || pseudoMode == "p" {
			break
		}
	}

	return pseudoMode
}

func overwriteFileHelp() string {
	return `
        Y - yes, overwrite
        n - no, do not overwrite
        a - all, overwrite this and all others
        q - quit, abort
        d - diff, show the differences between the old and the new
        p - patch, apply patch (check "diff" first)
        h - help, show this help

`
}

func GeneratePathAndFileFromTemplateString(path []string, content string, templateVar TemplateVar) error {
	var err error

	fileName, filePathSliced := path[len(path)-1], path[:len(path)-1]
	filePath := sliceToPath(filePathSliced)
	pseudoMode := "w"
	fullPath := filepath.Join(filePath, fileName)

	content, err = GenerateFromTemplateString(content, templateVar)
	if err != nil {
		return err
	}

	if fileutil.DirOrFileExists(fullPath) {
		currentContent, err := fileutil.ReadTextFile(filePath, fileName)
		if err != nil {
			return err
		}

		newContent, err := FmtNewContent(content)
		if err != nil {
			return err
		}

		if newContent == currentContent {
			pseudoMode = "i"
		} else {
			pseudoMode = overwriteFile(content, filePath, fileName)
		}
	}

	if pseudoMode == "s" {
		notify.Skip(fullPath)
	} else if pseudoMode == "p" {
		notify.Patch(fullPath)
	} else {
		err = fileutil.PersistFile(content, filePath, fileName, pseudoMode)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreatePathAndFileFromTemplateString(path []string, content string, templateVar TemplateVar) error {
	fileName, filePath := path[len(path)-1], path[:len(path)-1]
	return fileutil.SaveTextFile(content, sliceToPath(filePath), fileName)
}

func GenerateRoutesNewCode(content string, templateVar TemplateVar) (string, error) {
	return GenerateFromTemplateString(content, templateVar)
}

func GenerateMigrateNewCode(content string, templateVar TemplateVar) (string, error) {
	return GenerateFromTemplateString(content, templateVar)
}

func GenerateAuthorizeNewCode(content string, templateVar TemplateVar) (string, error) {
	return GenerateFromTemplateString(content, templateVar)
}

func HandlePathInfo(path []string) (string, string) {
	var basePath, fileName string

	for index, value := range path {
		if index+1 != len(path) {
			basePath = filepath.Join(basePath, value)
		} else {
			fileName = value
		}
	}

	return basePath, fileName
}

func SecureRandom(size int) string {
	var letters = []rune("0123456789abcdefABCDEF")

	rand.Seed(time.Now().UnixNano())

	b := make([]rune, size)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func BuildRootAppPath(appRepository string) (string, error) {
	_, err := os.Getwd()
	if err != nil {
		return "", err
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	path := filepath.Join(usr.HomeDir, "go", "src")

	urlPaths := strings.Split(appRepository, "/")
	for _, element := range urlPaths {
		path = filepath.Join(path, element)
	}

	if fileutil.DirOrFileExists(path) {
		return "", errors.New("directory \"" + path + "\" already exists\n")
	}

	return path, nil
}

func CreateRootAppPath(rootAppPath string) error {
	err := os.MkdirAll(rootAppPath, 0775)
	if err != nil {
		return err
	}

	notify.Created(rootAppPath)

	return nil
}

func NotifyNewApp(rootAppPath string) {
	notify.NewApp(rootAppPath)
}

func GenerateCertificates(rootAppPath string) error {
	var out bytes.Buffer

	err := os.MkdirAll(filepath.Join(rootAppPath, "config", "keys"), 0775)
	if err != nil {
		return err
	}

	cmd := exec.Command("openssl", "genrsa", "-out", filepath.Join(rootAppPath, "config", "keys", "app.key.rsa"), "2048")
	cmd.Stdout = &out
	err = cmd.Run()

	if err != nil {
		notify.Warn("Could not generate certificates files. Check if openssl is installed and execute both command lines below:")
		notify.Warn("openssl genrsa -out " + filepath.Join(rootAppPath, "config", "keys", "app.key.rsa") + " 2048")
		notify.Warn("openssl rsa -in " + filepath.Join(rootAppPath, "config", "keys", "app.key.rsa") + " -pubout > " + filepath.Join(rootAppPath, "config", "keys", "app.key.rsa.pub"))
	} else {
		notify.Created(filepath.Join(rootAppPath, "config", "keys", "app.key.rsa"))

		cmd := exec.Command("openssl", "rsa", "-in", filepath.Join(rootAppPath, "config", "keys", "app.key.rsa"), "-pubout")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			notify.Error(err)
			notify.Warn("Could not generate public certificate file. Check if openssl is installed and execute the command line below:")
			notify.Warn("openssl rsa -in " + filepath.Join(rootAppPath, "config", "keys", "app.key.rsa") + " -pubout > " + filepath.Join(rootAppPath, "config", "keys", "app.key.rsa.pub"))
		} else {
			err = fileutil.SaveTextFile(out.String(), filepath.Join(rootAppPath, "config", "keys"), "app.key.rsa.pub")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func GetColumnInfo(columnName string, columnType string, extra string) EntityColumn {
	var regexpText = regexp.MustCompile(`text`)
	var regexpString = regexp.MustCompile(`string`)
	var regexpDecimal = regexp.MustCompile(`float|double|decimal`)
	var regexpBigint = regexp.MustCompile(`bigint`)
	var regexpInteger = regexp.MustCompile(`int|integer`)
	var regexpUnsignedInteger = regexp.MustCompile(`uint`)
	var regexpDatetime = regexp.MustCompile(`datetime`)
	var regexpBoolean = regexp.MustCompile(`bool`)
	var regexpReference = regexp.MustCompile(`reference`)
	var migrateType, migrateExtra string

	isRelation := false

	if regexpText.MatchString(columnType) {
		extra = "`gorm:\"type:text\"`"
		migrateExtra = `nil`
		columnType = "string"
		migrateType = "Text"
	} else if regexpString.MatchString(columnType) || regexpText.MatchString(columnType) {
		extra, migrateExtra = extraSpecificationForString(extra)
		columnType = "string"
		migrateType = "String"
	} else if regexpUnsignedInteger.MatchString(columnType) {
		extra, migrateExtra = extraSpecificationForIntegers(extra)
		columnType = "uint"
		migrateType = "Integer"
	} else if regexpBigint.MatchString(columnType) {
		extra, migrateExtra = extraSpecificationForIntegers(extra)
		columnType = "int64"
		migrateType = "Bigint"
	} else if regexpInteger.MatchString(columnType) {
		extra, migrateExtra = extraSpecificationForIntegers(extra)
		columnType = "int"
		migrateType = "Integer"
	} else if regexpDatetime.MatchString(columnType) {
		extra, migrateExtra = extraSpecificationForDatetime(extra)
		columnType = "*time.Time"
		migrateType = "Datetime"
	} else if regexpBoolean.MatchString(columnType) {
		extra, migrateExtra = extraSpecificationForBoolean(extra)
		columnType = "bool"
		migrateType = "Boolean"
	} else if regexpDecimal.MatchString(columnType) {
		extra, migrateExtra = extraSpecificationForDecimals(columnType, extra)
		columnType = "float64"
		migrateType = "Numeric"
	} else if regexpReference.MatchString(columnType) {
		extra = "`gorm:\"index\"`"
		columnType = "uint64"
		columnName = columnName + "_ID"
		migrateType = "References"
		isRelation = true
	}

	return EntityColumn{
		Name:                strcase.ToCamel(columnName),
		NameSnakeCase:       strcase.ToSnake(columnName),
		NameSnakeCasePlural: inflection.Plural(strcase.ToSnake(columnName)),
		Type:                columnType,
		Extras:              extra,
		IsRelation:          isRelation,
		IsForeignKey:        false,
		MigrateType:         migrateType,
		MigrateExtra:        migrateExtra,
	}
}

func UpdateMigrate(basePath string, templateVar TemplateVar) {
	newCode, _ := GenerateMigrateNewCode(templatecommon.MigrateContent, templateVar)
	currentFullCode, _ := fileutil.ReadTextFile(filepath.Join(basePath, "db", "schema"), "schema.go")
	newFullCode, err := newmigrate.AppendNewCode(newCode, currentFullCode)

	if err != nil {
		notify.WarnAppendToMigrate(err, newCode)
	} else if newFullCode == "" {
		notify.Identical(filepath.Join(basePath, "db", "schema", "schema.go"))
	} else {
		fileutil.UpdateTextFile(newFullCode, filepath.Join(basePath, "db", "schema"), "schema.go")
	}
}

func SetEntityName(name string) EntityName {
	nameSingular := inflection.Singular(name)
	namePlural := inflection.Plural(nameSingular)

	entityName := EntityName{
		CamelCase:            strcase.ToCamel(nameSingular),
		CamelCasePlural:      strcase.ToCamel(namePlural),
		LowerCamelCase:       strcase.ToLowerCamel(nameSingular),
		LowerCamelCasePlural: strcase.ToLowerCamel(namePlural),
		SnakeCase:            strcase.ToSnake(nameSingular),
		SnakeCasePlural:      strcase.ToSnake(namePlural),
		LowerCase:            strings.ToLower(strcase.ToCamel(nameSingular)),
	}

	return entityName
}

func SetMigrationMetadata(name string) MigrationMetadata {
	nameCamelCase := strcase.ToCamel(name)
	nameSnakeCase := strcase.ToSnake(name)

	migrationMetadata := MigrationMetadata{
		Name:          nameCamelCase,
		FileNameSufix: nameSnakeCase,
		Version:       time.Now().Format("20060102150405"),
	}

	if strings.HasPrefix(nameSnakeCase, "add") {
		migrationMetadata.Type = "ADD_COLUMN"
		pieces := strings.Split(nameSnakeCase, "_to_")
		migrationMetadata.Entity = inflection.Singular(pieces[len(pieces)-1])
	} else if strings.HasPrefix(nameSnakeCase, "remove_") {
		migrationMetadata.Type = "REMOVE_COLUMN"
		pieces := strings.Split(nameSnakeCase, "_from_")
		migrationMetadata.Entity = inflection.Singular(pieces[len(pieces)-1])
	} else if strings.HasPrefix(nameSnakeCase, "drop") {
		migrationMetadata.Type = "DROP_TABLE"
		migrationMetadata.Entity = inflection.Singular(strings.TrimPrefix(nameSnakeCase, "drop_table_"))
	} else if strings.HasPrefix(nameSnakeCase, "create") {
		migrationMetadata.Type = "CREATE_TABLE"
		migrationMetadata.Entity = inflection.Singular(strings.TrimPrefix(nameSnakeCase, "create_table_"))
	} else {
		migrationMetadata.Type = "GENERAL_CHANGE"
		migrationMetadata.Entity = ""
	}

	return migrationMetadata
}

func sliceToPath(path []string) string {
	var filePath string

	for index, element := range path {
		if index > 0 {
			filePath = filepath.Join(filePath, element)
		} else {
			filePath = element
		}
	}

	return filePath
}

func extraSpecificationForString(extra string) (string, string) {
	var index string
	var migrate string

	if extra == "index" {
		index = ";index"
		migrate = `map[string]interface{}{"index": true}`
	} else {
		migrate = "nil"
	}

	return "`gorm:\"type:varchar(255)" + index + "\"`", migrate
}

func extraSpecificationForIntegers(extra string) (string, string) {
	var index string
	var migrate string

	if extra == "index" {
		index = "`gorm:\"index\"`"
		migrate = `map[string]interface{}{"index": true}`
	} else {
		migrate = "nil"
	}

	return index, migrate
}

func extraSpecificationForDecimals(columnType string, extra string) (string, string) {
	var index string
	var migrate string
	var regexpPrecision = regexp.MustCompile(`\((\d+)(\,(\d+)){0,1}\)`)
	var subMatches [][]string
	var migrationExtras []string

	if extra == "index" {
		index = ";index"
		migrationExtras = append(migrationExtras, `"index": true`)
	}

	subMatches = regexpPrecision.FindAllStringSubmatch(columnType, -1)
	if len(subMatches) > 0 && len(subMatches[0]) > 0 {
		migrationExtras = append(migrationExtras, `"precision": `+subMatches[0][1])

		if len(subMatches[0]) == 4 {
			migrationExtras = append(migrationExtras, `"scale": `+subMatches[0][3])
		}
	}
	fmt.Println(columnType)
	fmt.Println(extra)
	fmt.Println(subMatches)
	fmt.Println(migrationExtras)

	if len(migrationExtras) > 0 {
		migrate = `map[string]interface{}{ ` + strings.Join(migrationExtras, ",") + `}`
	} else {
		migrate = "nil"
	}

	return "`gorm:\"type:decimal\"" + index + "`", migrate
}

func extraSpecificationForDatetime(extra string) (string, string) {
	var index string
	var migrate string

	if extra == "index" {
		index = ";index"
		migrate = `map[string]interface{}{"index": true}`
	} else {
		migrate = "nil"
	}

	return "`gorm:\"default:null\"" + index + "`", migrate
}

func extraSpecificationForBoolean(extra string) (string, string) {
	var index string
	var migrate string

	if extra == "index" {
		index = "`gorm:\"default:null\";index`"
		migrate = `map[string]interface{}{"index": true}`
	} else if extra == "true" || extra == "t" {
		index = "`gorm:\"default:true\"`"
		migrate = `map[string]interface{}{"default": true}`
	} else if extra == "false" || extra == "f" {
		index = "`gorm:\"default:false\"`"
		migrate = `map[string]interface{}{"default": false}`
	} else {
		migrate = "nil"
	}

	return index, migrate
}

func extraSpecificationForReference() (string, string) {
	return "`gorm:\"index\"`", `map[string]interface{}{"foreign_key": true}`
}
