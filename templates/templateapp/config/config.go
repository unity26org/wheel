package config

var Path = []string{"config", "config.go"}

var Content = `package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"{{ .AppRepository }}/commons/log"
)

type AppConfig struct {
	AppName                        string   ` + "`" + `yaml:"app_name"` + "`" + `
	AppRepository                  string   ` + "`" + `yaml:"app_repository"` + "`" + `
	SecretKey                      string   ` + "`" + `yaml:"secret_key"` + "`" + `
	ResetPasswordExpirationSeconds int      ` + "`" + `yaml:"reset_password_expiration_seconds"` + "`" + `
	ResetPasswordUrl               string   ` + "`" + `yaml:"reset_password_url"` + "`" + `
	TokenExpirationSeconds         int      ` + "`" + `yaml:"token_expiration_seconds"` + "`" + `
	Locales                        []string ` + "`" + `yaml:"locales"` + "`" + `
}

var App AppConfig

func readAppConfigFile() []byte {
	data, err := ioutil.ReadFile("./config/app.yml")
	if err != nil {
		log.Error.Fatal(err)
	}

	return data
}

func init() {
	err := yaml.Unmarshal(readAppConfigFile(), &App)
	if err != nil {
		log.Error.Fatalf("error: %v", err)
	}
}`
