package templatecrud

var MigrateContent = `model.Db.AutoMigrate(&entities.{{ .EntityName.CamelCase }}{})
{{- $filteredEntityColumns := filterEntityColumnsForeignKeysOnly .EntityColumns }}
{{- range $index, $element := $filteredEntityColumns }} 
model.Db.Model(&entities.{{ $.EntityName.CamelCase }}{}).AddForeignKey("{{ .NameSnakeCase }}_id", "{{ .NameSnakeCasePlural }}(id)", "NO ACTION", "NO ACTION")
{{- end }}`
