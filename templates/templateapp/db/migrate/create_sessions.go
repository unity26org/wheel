package migrate

var SessionPath = []string{"db", "migrate", "create_sessions.go"}

var SessionContent = `package migrate

import (
	"{{ .AppRepository }}/db/schema/data/col"
)

type Version{{ .MigrationMetadata.Version }} struct {
}

func (m *Version{{ .MigrationMetadata.Version }}) CreateSessions(direction string) error {
  var err error
  
	if direction == "up" {
		err = m.up()
	} else if direction == "down" {
		err = m.down()
	}
  
  return err
}

func (m *Version{{ .MigrationMetadata.Version }}) up() error {
	err := CreateTable("sessions", []col.Info{
		col.References("user", map[string]interface{}{"foreign_key": true}),
		col.String("jti", map[string]interface{}{"unique": true}),
		col.String("app", nil),
		col.Integer("requests", map[string]interface{}{"default": 0}),
		col.Integer("expires_in", map[string]interface{}{"default": 0}),
		col.String("address", nil),
		col.Boolean("active", map[string]interface{}{"default": true}),
		col.Datetime("last_request_at", map[string]interface{}{"null": false}),
		col.Datetime("created_at", nil),
		col.Datetime("updated_at", nil),
		col.Datetime("expires_at", nil),
	})

	if err != nil {
		return err
  }
  
	return nil
}

func (m *Version{{ .MigrationMetadata.Version }}) down() error {
	err := DropTable("sessions")

	if err != nil {
		return err
  }
  
	return nil
}`
