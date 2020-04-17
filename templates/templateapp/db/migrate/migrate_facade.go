package migrate

var FacadePath = []string{"db", "migrate", "migrate_facade.go"}

var FacadeContent = `package migrate

import (
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/db/migrate/adapter"
	"{{ .AppRepository }}/db/migrate/adapter/postgresql"
	"{{ .AppRepository }}/db/schema/data/col"
	"time"
  "fmt"
)

var sql string
var Ddl adapter.Adapter

func init() {
	Ddl = postgresql.Ddl{}
}

func CreateTable(table string, columns []col.Info) error {
	sql = Ddl.CreateTable(table, columns)
	return PrintAndExec("CreateTable(\""+table+"\")", sql)
}

func DropTable(table string) error {
	sql = Ddl.DropTable(table)
	return PrintAndExec("DropTable(\""+table+"\")", sql)
}

func AddIndex(table string, columns []string, options col.Index) error {
	sql = Ddl.AddIndex(table, columns, options)
	return PrintAndExec("AddIndex(\""+table+"\")", sql)
}

// options are: "columns" or "column" (slice or string), "index" (string), "concurrently" (bool) and "option" ("cascade" or "restrict")
func RemoveIndex(table string, options map[string]interface{}) error {
	sql = Ddl.RemoveIndex(table, options)
	return PrintAndExec("RemoveIndex(\""+table+"\")", sql)
}

func AddColumn(table string, column col.Info) error {
	sql = Ddl.AddColumn(table, column)
	return PrintAndExec("AddColumn(\""+table+"\", \""+column.Name+"\")", sql)
}

func RenameColumn(table string, column string, newColumnName string) error {
	sql = Ddl.RenameColumn(table, column, newColumnName)
	return PrintAndExec("RenameColumn(\""+table+"\", \""+column+"\", \""+newColumnName+"\")", sql)
}

func ChangeColumnType(table string, column string, newColumnType string) error {
	sql = Ddl.ChangeColumnType(table, column, newColumnType)
	return PrintAndExec("ChangeColumnType(\""+table+"\", \""+column+"\", \""+newColumnType+"\")", sql)
}

func ChangeColumnNull(table string, column string, isNull bool) error {
	sql = Ddl.ChangeColumnNull(table, column, isNull)
	return PrintAndExec("ChangeColumnNull(\""+table+"\", \""+column+"\")", sql)
}

func ChangeColumnDefault(table string, column string, defaultValue interface{}) error {
	sql = Ddl.ChangeColumnDefault(table, column, defaultValue)
	return PrintAndExec("ChangeColumnDefault(\""+table+"\", \""+column+"\")", sql)
}

func RemoveColumn(table string, column string) error {
	sql = Ddl.RemoveColumn(table, column)
	return PrintAndExec("RemoveColumn(\""+table+"\", \""+column+"\")", sql)
}

// options are: "column", "on_delete", "on_update", "name" and "primary_key" (all are strings)
// "on_delete" and "on_update" available values are "nullify", "cascade" and "restrict"
func AddForeignKey(fromTable string, toTable string, options map[string]string) error {
	sql = Ddl.AddForeignKey(fromTable, toTable, options)
	return PrintAndExec("AddForeignKey(\""+fromTable+"\", \""+toTable+"\")", sql)
}

// options are: "to_table", "column" and "name"
func RemoveForeignKey(table string, options map[string]string) error {
	sql = Ddl.RemoveForeignKey(table, options)
	return PrintAndExec("RemoveForeignKey(\""+table+"\")", sql)
}

func PrintAndExec(message string, sql string) error {
	fmt.Printf("|   ─> %s\n", message)
	return Exec(sql)
}

func Exec(sql string) error {
	t0 := time.Now()
	err := model.Db.Exec(sql).Error
	if err != nil {
		return err
	} else {
		t1 := time.Now()
		f := ((float64(t1.Sub(t0)) / 1e6) / float64(1000))
		fmt.Printf("|   ─> %.4fs\n", f)
		fmt.Printf("|\n")
		return nil
	}
}`
