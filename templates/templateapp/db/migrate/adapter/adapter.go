package adapter

var Path = []string{"db", "migrate", "adapter", "adapter.go"}

var Content = `package adapter

import (
	"{{ .AppRepository }}/db/schema/data/col"
)

type Adapter interface {
	CreateTable(table string, columns []col.Info) string
	DropTable(table string) string
	AddIndex(table string, columns []string, options col.Index) string
	RemoveIndex(table string, options map[string]interface{}) string
	AddColumn(table string, column col.Info) string
	RenameColumn(table string, column string, newColumnName string) string
	ChangeColumnType(table string, column string, newColumnType string) string
	ChangeColumnNull(table string, column string, isNull bool) string
	ChangeColumnDefault(table string, column string, defaultValue interface{}) string
	RemoveColumn(table string, column string) string
	AddForeignKey(fromTable string, toTable string, options map[string]string) string
	RemoveForeignKey(table string, options map[string]string) string
}`
