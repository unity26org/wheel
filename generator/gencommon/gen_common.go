package gencommon

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/unity26org/wheel/commons/diff"
	"github.com/unity26org/wheel/commons/fileutil"
	"github.com/unity26org/wheel/commons/notify"
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
	IsReference         bool
	IsForeignKey        bool
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

type TemplateVar struct {
	AppRepository string
	AppName       string
	SecretKey     string
	EntityName    EntityName
	EntityColumns []EntityColumn
}

var yesToAll = false

func DirOrFileExists(fullPath string) bool {
	_, err := os.Stat(fullPath)
	return !os.IsNotExist(err)
}

func UpdateTextFile(content string, filePath string, fileName string) {
	persistFile(content, filePath, fileName, "a")
}

func SaveTextFile(content string, filePath string, fileName string) {
	persistFile(content, filePath, fileName, "w")
}

func persistFile(content string, filePath string, fileName string, pseudoMode string) {
	err := os.MkdirAll(filePath, 0775)
	notify.FatalIfError(err)

	fullPath := filepath.Join(filePath, fileName)

	f, err := os.Create(fullPath)
	notify.FatalIfError(err)

	defer f.Close()

	_, err = f.WriteString(content)
	notify.FatalIfError(err)

	f.Sync()

	fileutil.GoFmtFile(fullPath)

	switch pseudoMode {
	case "w":
		notify.Created(fullPath)
	case "a":
		notify.Updated(fullPath)
	case "i":
		notify.Identical(fullPath)
	case "f":
		notify.Force(fullPath)
	case "s":
		notify.Skip(fullPath)
	}

}

func DestroyDirOrFile(fullPath string) {
	err := os.Remove(fullPath)
	notify.FatalIfError(err)
}

func FmtNewContent(content string) string {
	var fileName string

	tmpfile, err := ioutil.TempFile("", "wheel-new-content*.go")
	notify.FatalIfError(err)

	fileName = strings.TrimPrefix(tmpfile.Name(), os.TempDir())
	fileName = strings.TrimPrefix(fileName, `/`)
	fileName = strings.TrimPrefix(fileName, `\`)

	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(content)
	notify.FatalIfError(err)

	err = tmpfile.Close()
	notify.FatalIfError(err)

	fileutil.GoFmtFile(tmpfile.Name())

	return fileutil.ReadTextFile(os.TempDir(), fileName)
}

func GetAppConfig() AppConfig {
	var appConfig AppConfig

	err := yaml.Unmarshal(fileutil.ReadBytesFile(filepath.Join(".", "config"), "app.yml"), &appConfig)
	notify.FatalIfError(err)

	return appConfig
}

func GenerateFromTemplateString(content string, templateVar TemplateVar) string {
	var buffContent bytes.Buffer

	FuncMap := template.FuncMap{
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
	}

	tmpl, err := template.New("T").Funcs(FuncMap).Parse(content)
	notify.FatalIfError(err)

	err = tmpl.Execute(&buffContent, templateVar)
	notify.FatalIfError(err)

	return buffContent.String()
}

func GenerateFromTemplateFile(templatePath string, templateVar TemplateVar) string {
	var content bytes.Buffer

	tmpl, err := template.ParseFiles(templatePath)
	notify.FatalIfError(err)

	err = tmpl.Execute(&content, &templateVar)
	notify.FatalIfError(err)

	return content.String()
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

func GeneratePathAndFileFromTemplateString(path []string, content string, templateVar TemplateVar) {
	fileName, filePathSliced := path[len(path)-1], path[:len(path)-1]
	filePath := sliceToPath(filePathSliced)
	pseudoMode := "w"
	fullPath := filepath.Join(filePath, fileName)

	content = GenerateFromTemplateString(content, templateVar)

	if DirOrFileExists(fullPath) {
		if FmtNewContent(content) == fileutil.ReadTextFile(filePath, fileName) {
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
		persistFile(content, filePath, fileName, pseudoMode)
	}
}

func CreatePathAndFileFromTemplateString(path []string, content string, templateVar TemplateVar) {
	fileName, filePath := path[len(path)-1], path[:len(path)-1]
	SaveTextFile(content, sliceToPath(filePath), fileName)
}

func GenerateRoutesNewCode(content string, templateVar TemplateVar) string {
	return GenerateFromTemplateString(content, templateVar)
}

func GenerateMigrateNewCode(content string, templateVar TemplateVar) string {
	return GenerateFromTemplateString(content, templateVar)
}

func GenerateAuthorizeNewCode(content string, templateVar TemplateVar) string {
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

func BuildRootAppPath(appRepository string) string {
	_, err := os.Getwd()
	notify.FatalIfError(err)

	usr, err := user.Current()
	notify.FatalIfError(err)

	path := filepath.Join(usr.HomeDir, "go", "src")

	urlPaths := strings.Split(appRepository, "/")
	for _, element := range urlPaths {
		path = filepath.Join(path, element)
	}

	if DirOrFileExists(path) {
		errDirOrFileExists := errors.New("directory \"" + path + "\" already exists\n")
		notify.FatalIfError(errDirOrFileExists)
	}

	return path
}

func CreateRootAppPath(rootAppPath string) {
	err := os.MkdirAll(rootAppPath, 0775)
	notify.FatalIfError(err)

	notify.Created(rootAppPath)
}

func NotifyNewApp(rootAppPath string) {
	notify.NewApp(rootAppPath)
}

func GenerateCertificates(rootAppPath string) {
	var out bytes.Buffer

	err := os.MkdirAll(filepath.Join(rootAppPath, "config", "keys"), 0775)
	notify.FatalIfError(err)

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
			SaveTextFile(out.String(), filepath.Join(rootAppPath, "config", "keys"), "app.key.rsa.pub")
		}
	}
}

func GetColumnInfo(columnName string, columnType string, extra string) (string, string, string, bool) {
	var regexpText = regexp.MustCompile(`text`)
	var regexpString = regexp.MustCompile(`string`)
	var regexpDecimal = regexp.MustCompile(`float|double|decimal`)
	var regexpInteger = regexp.MustCompile(`int|integer`)
	var regexpUnsignedInteger = regexp.MustCompile(`uint`)
	var regexpDatetime = regexp.MustCompile(`datetime`)
	var regexpBoolean = regexp.MustCompile(`bool`)
	var regexpReference = regexp.MustCompile(`reference`)

	isReference := false

	if regexpText.MatchString(columnType) {
		columnType = "string"
		extra = "`gorm:\"type:text\"`"
	} else if regexpString.MatchString(columnType) || regexpText.MatchString(columnType) {
		columnType = "string"
		extra = gormSpecificationForString(extra)
	} else if regexpUnsignedInteger.MatchString(columnType) {
		columnType = "uint"
		extra = gormSpecificationForIntegers(extra)
	} else if regexpInteger.MatchString(columnType) {
		columnType = "int64"
		extra = gormSpecificationForIntegers(extra)
	} else if regexpDatetime.MatchString(columnType) {
		columnType = "*time.Time"
		extra = gormSpecificationForDatetime(extra)
	} else if regexpBoolean.MatchString(columnType) {
		columnType = "bool"
		extra = gormSpecificationForBoolean(extra)
	} else if regexpDecimal.MatchString(columnType) {
		columnType = "float64"
		extra = gormSpecificationForDecimals(extra)
	} else if regexpReference.MatchString(columnType) {
		columnType = "uint"
		extra = "`gorm:\"index\"`"
		columnName = columnName + "_ID"
		isReference = true
	}

	return columnName, columnType, extra, isReference
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

func gormSpecificationForString(extra string) string {
	var index string

	if extra == "index" {
		index = ";index"
	} else {
		index = ""
	}

	return "`gorm:\"type:varchar(255)" + index + "\"`"
}

func gormSpecificationForIntegers(extra string) string {
	var index string

	if extra == "index" {
		index = "`gorm:\"index\"`"
	} else {
		index = ""
	}

	return index
}

func gormSpecificationForDecimals(extra string) string {
	var index string

	if extra == "index" {
		index = ";index"
	} else {
		index = ""
	}

	return "`gorm:\"type:decimal\"" + index + "`"
}

func gormSpecificationForDatetime(extra string) string {
	var index string

	if extra == "index" {
		index = ";index"
	} else {
		index = ""
	}

	return "`gorm:\"default:null\"" + index + "`"
}

func gormSpecificationForBoolean(extra string) string {
	var index string

	if extra == "index" {
		index = "`gorm:\"default:null\";index`"
	} else if extra == "true" || extra == "t" {
		index = "`gorm:\"default:true\"`"
	} else if extra == "false" || extra == "f" {
		index = "`gorm:\"default:false\"`"
	} else {
		index = ""
	}

	return index
}

func gormSpecificationForReference() string {
	return "`gorm:\"index\"`"
}
