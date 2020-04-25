package migrate

var UserPath = []string{"db", "migrate", "create_users.go"}

var UserContent = `package migrate

import (
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/commons/crypto"
	"{{ .AppRepository }}/db/entities"
	"{{ .AppRepository }}/db/schema/data/col"
)

type Version{{ .MigrationMetadata.Version }} struct {
}

func (m *Version{{ .MigrationMetadata.Version }}) CreateUsers(direction string) error {
  var err error
  
	if direction == "up" {
		err = m.up()
	} else if direction == "down" {
		err = m.down()
	}
  
  return err
}

func (m *Version{{ .MigrationMetadata.Version }}) up() error {
	err := CreateTable("users", []col.Info{
		col.String("name", nil),
		col.String("email", map[string]interface{}{"unique": true}),
		col.Boolean("admin", nil),
		col.String("password", nil),
		col.String("reset_password_token", nil),
		col.Datetime("reset_password_sent_at", nil),
		col.String("locale", nil),
		col.Datetime("deleted_at", nil),
	})

	if err != nil {
		return err
	} else {
		model.Db.Create(&entities.User{Name: "User Name", Email: "user@example.com", Password: crypto.SetPassword("Secret123!"), Locale: "en", Admin: true})
	}

	return nil
}

func (m *Version{{ .MigrationMetadata.Version }}) down() error {
	err := DropTable("users")
	if err != nil {
		return err
  }
  
	return nil
}`
