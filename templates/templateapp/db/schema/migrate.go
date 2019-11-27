package schema

var Path = []string{"db", "schema", "migrate.go"}

var Content = `package schema

import (
	"{{ .AppRepository }}/app/user"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/commons/crypto"
	"{{ .AppRepository }}/db/entities"
)

func Migrate() {
	model.Db.AutoMigrate(&entities.User{})

	_, err := user.FindByEmail("user@example.com")
	if err != nil {
 		model.Db.Create(&entities.User{Name: "User Name", Email: "user@example.com", Password: crypto.SetPassword("!Secret.123!"), Locale: "en", Admin: true})
	}

	model.Db.AutoMigrate(&entities.Session{})
	model.Db.Model(&entities.Session{}).AddForeignKey("user_id", "users(id)", "NO ACTION", "NO ACTION")
}`
