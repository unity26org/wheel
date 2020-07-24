package model

var Path = []string{"commons", "app", "model", "model.go"}

var Content = `package model

import (
	"errors"
	"github.com/jinzhu/gorm"
  {{- if eq .Database "postgres" }}
	_ "github.com/jinzhu/gorm/dialects/postgres"
  {{- else if eq .Database "mysql" }}
  _ "github.com/go-sql-driver/mysql"
  {{- end }}
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
	"strings"
	"{{ .AppRepository }}/commons/conversor"
	"{{ .AppRepository }}/commons/log"
)

var Db *gorm.DB

var Errors []string

type Query struct {
	Db    *gorm.DB
	Table interface{}
}

type DbConfig struct {
	Host    string ` + "`" + `yaml:"host"` + "`" + `
	User    string ` + "`" + `yaml:"user"` + "`" + `
	Dbname  string ` + "`" + `yaml:"dbname"` + "`" + `
	Port    int    ` + "`" + `yaml:"port"` + "`" + `
	Adapter string ` + "`" + `yaml:"adapter"` + "`" + `
	SslMode string ` + "`" + `yaml:"sslmode"` + "`" + `
	Pool    int    ` + "`" + `yaml:"pool"` + "`" + `
}

var CurrentDbConfig DbConfig

func Connect() {
	var err error

	dbConfig := loadDatabaseConfigFile()
	Db, err = gorm.Open(dbConfig["adapter"], stringfyDatabaseConfigFile(dbConfig))

	if err != nil {
		log.Fatal.Println(err)
		panic("failed connect to database")
	} else {
		log.Info.Println("connected to the database successfully")
	}

	pool, err := strconv.Atoi(dbConfig["pool"])
	if err != nil {
		log.Fatal.Println(err)
	} else {
		log.Info.Printf("database pool of connections: %d", pool)
	}

	Db.DB().SetMaxIdleConns(pool)
}

func Disconnect() {
	defer Db.Close()
}

func TableName(table interface{}) string {
	return Db.NewScope(table).GetModelStruct().TableName(Db)
}

func GetColumnType(table interface{}, columnName string) (string, error) {
	field, ok := Db.NewScope(table).FieldByName(columnName)

	if ok {
		return field.Field.Type().String(), nil
	} else {
		return "", errors.New("column was not found")
	}
}

func GetColumnValue(table interface{}, columnName string) (interface{}, error) {
	field, ok := Db.NewScope(table).FieldByName(columnName)

	if ok {
		return field.Field.Interface(), nil
	} else {
		return "", errors.New("column was not found")
	}
}

func SetColumnValue(table interface{}, columnName string, value string) error {
	field, ok := Db.NewScope(table).FieldByName(columnName)

	if ok {
		columnType, _ := GetColumnType(table, columnName)
		valueInterface, _ := conversor.StringToInterface(columnType, value)
		return field.Set(valueInterface)
	} else {
		return errors.New("column was not found")
	}
}

func ColumnsFromTable(table interface{}, all bool) []string {
	var columns []string
	fields := Db.NewScope(table).Fields()

	for _, field := range fields {
		if !all && ((field.Names[0] == "Model") || (field.Relationship != nil)) {
			continue
		}
		columns = append(columns, field.DBName)
	}

	return columns
}

// PACKAGE METHODS

func loadDatabaseConfigFile() map[string]string {
	config := make(map[string]string)

	err := yaml.Unmarshal(readDatabaseConfigFile(), &config)
	if err != nil {
		log.Fatal.Printf("error: %v\n", err)
	}

	if config["pool"] == "" {
		config["pool"] = "5"
	}

	return config
}

func readDatabaseConfigFile() []byte {
	data, err := ioutil.ReadFile("./config/database.yml")
	if err != nil {
		log.Fatal.Println(err)
	}

	return data
}

func stringfyDatabaseConfigFile(mapped map[string]string) string {
	var arr []string

	if mapped["adapter"] == "mysql" {
		return mapped["user"] + ` + "`" + `:` + "`" + ` + mapped["password"] + ` + "`" + `@/` + "`" + ` + mapped["dbname"] + ` + "`" + `?charset=utf8&parseTime=True&loc=Local` + "`" + `
	} else {
		for key, value := range mapped {
			if key != "pool" && key != "adapter" {
				arr = append(arr, key+"='"+value+"'")
			}
		}

		return strings.Join(arr, " ")
	}
}

func init() {
	err := yaml.Unmarshal(readDatabaseConfigFile(), &CurrentDbConfig)
	if err != nil {
		log.Error.Fatalf("error: %v", err)
	}
}`
