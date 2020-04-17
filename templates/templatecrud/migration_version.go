package templatecrud

var MigrationVersionPath = []string{"db", "migrate", "create.go"}

var MigrationVersionContent = `package migrate

import (
	"{{ .AppRepository }}/db/schema/data/col"
)

type Version{{ .MigrationMetadata.Version }} struct {
}

func (m *Version{{ .MigrationMetadata.Version }}) {{ .MigrationMetadata.Name }}(direction string) error {
  var err error
  
	if direction == "up" {
		err = m.up()
	} else if direction == "down" {
		err = m.down()
	}
  
  return err
}

func (m *Version{{ .MigrationMetadata.Version }}) up() error {
  var err error
  
  {{- if checkMigrationType "CREATE_TABLE" }}
	err = CreateTable("{{ .EntityName.SnakeCasePlural }}", []col.Info{
    {{- $filteredEntityColumns := filterEntityColumnsNotRelations .EntityColumns }}
    {{- range $index, $element := $filteredEntityColumns }}
    col.{{ .MigrateType }}("{{ .NameSnakeCase }}", {{ .MigrateExtra }}),
    {{- end }}
  })
  
  if (err != nil) {
    return err
  }
  
  {{- else if checkMigrationType "ADD_COLUMN" }}
  
  {{- $EntityNameSnakeCasePlural := .EntityName.SnakeCasePlural }}
  {{- $filteredEntityColumns := filterEntityColumnsNotRelations .EntityColumns }}
  {{- range $index, $element := $filteredEntityColumns }} 
	err = AddColumn("{{ $EntityNameSnakeCasePlural }}", col.{{ .MigrateType }}("{{ .NameSnakeCase }}", {{ .MigrateExtra }}))
  if err != nil {
    return err
  }
  
  {{- end }}
  
  {{- else if checkMigrationType "REMOVE_COLUMN" }}
  
  {{- $EntityNameSnakeCasePlural := .EntityName.SnakeCasePlural }}
  {{- $filteredEntityColumns := filterEntityColumnsNotForeignKeys .EntityColumns }}
  {{- range $index, $element := .EntityColumns }}
	err = RemoveColumn("{{ $EntityNameSnakeCasePlural }}", "{{ .NameSnakeCase }}")
  if err != nil {
    return err
  }

  {{- end }}
  
  {{- end }}

  return nil
}

func (m *Version{{ .MigrationMetadata.Version }}) down() error {
  var err error
  
  {{- if checkMigrationType "CREATE_TABLE" }}
	err = DropTable("{{ .EntityName.SnakeCasePlural }}")
  if (err != nil) {
    return err
  }  

  {{- else if checkMigrationType "ADD_COLUMN" }}
  
  {{- $EntityNameSnakeCasePlural := .EntityName.SnakeCasePlural }}
  {{- $filteredEntityColumns := filterEntityColumnsNotForeignKeys .EntityColumns }}
  {{- range $index, $element := $filteredEntityColumns }} 

	err = RemoveColumn("{{ $EntityNameSnakeCasePlural }}", "{{ .NameSnakeCase }}")
  if err != nil {
    return err
  }
  
  {{- end }}
  
  return nil
  {{- else if checkMigrationType "REMOVE_COLUMN" }}

  {{- $EntityNameSnakeCasePlural := .EntityName.SnakeCasePlural }}
  {{- $filteredEntityColumns := filterEntityColumnsNotRelations .EntityColumns }}
  {{- range $index, $element := $filteredEntityColumns }} 
	err = AddColumn("{{ $EntityNameSnakeCasePlural }}", col.{{ .MigrateType }}("{{ .NameSnakeCase }}", {{ .MigrateExtra }}))
  if err != nil {
    return err
  }
  
  {{- end }}
  
  {{- end }}

  return nil
}`
