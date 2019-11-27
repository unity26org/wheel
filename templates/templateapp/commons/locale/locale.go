package locale

var Path = []string{"commons", "locale", "locale.go"}

var Content = `package locale

import (
	"{{ .AppRepository }}/commons/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Keys struct {
	Welcome                      string ` + "`" + `yaml:"welcome"` + "`" + `
	PasswordRecoveryInstructions string ` + "`" + `yaml:"password_recovery_instructions"` + "`" + `
}

var I18n Keys

func Load(locale string) {
	err := yaml.Unmarshal(readLocaleFile(locale), &I18n)
	if err != nil {
		log.Error.Fatal(err)
	}
}

func readLocaleFile(locale string) []byte {
	data, err := ioutil.ReadFile("./config/locales/" + locale + ".yml")
	if err != nil {
		log.Error.Fatal(err)
	}

	return data
}`
