package postgresql

var Path = []string{"db", "migrate", "adapter", "postgresql", "postgresql.go"}

var Content = `package postgresql

import (
	"fmt"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/commons/crypto"
	"{{ .AppRepository }}/db/schema/data/col"
	"github.com/jinzhu/inflection"
	"regexp"
	"strconv"
	"strings"
)

type Ddl struct {
}

func (ddl Ddl) CreateTable(table string, columns []col.Info) string {
	var createTable string
	var newColumns []string
	var newIndexes []string
	var newForeignKeys []string

	createTable = "BEGIN;\n"
	createTable = createTable + "CREATE TABLE " + table + " (\n"
	newColumns = append(newColumns, "id BIGSERIAL PRIMARY KEY")

	for _, value := range columns {
		newColumn, newIndex, newForeignKey := ddl.generateColumn(table, value)
		if newIndex != "" {
			newIndexes = append(newIndexes, newIndex)
		}

		if newForeignKey != "" {
			newForeignKeys = append(newForeignKeys, newForeignKey)
		}

		newColumns = append(newColumns, newColumn)
	}

	createTable = createTable + strings.Join(newColumns, ",\n") + "\n);"

	if len(newIndexes) > 0 {
		createTable = createTable + strings.Join(newIndexes, "\n")
	}

	if len(newForeignKeys) > 0 {
		createTable = createTable + strings.Join(newForeignKeys, "\n")
	}

	createTable = createTable + "COMMIT;"

	return createTable
}

func (ddl Ddl) DropTable(table string) string {
	return "DROP TABLE " + table + ";"
}

func (ddl Ddl) AddIndex(table string, columns []string, options col.Index) string {
	var sql string

	if options.Name == "" {
		options.Name = "index_" + table + "_on_" + strings.Join(columns, "_")
	}

	sql = "CREATE" + ddl.checkUnique(options.Unique) + " INDEX " + options.Name + " ON " + table + " USING btree (" + strings.Join(columns, ", ") + ");"

	return sql
}

// options are: "columns" or "column" (slice or string), "index" (string), "concurrently" (bool) and "option" ("cascade" or "restrict")
func (ddl Ddl) RemoveIndex(table string, options map[string]interface{}) string {
	var removeIndexSql string

	if index, ok := options["index"]; ok {
		removeIndexSql = fmt.Sprintf("%v", index)
	} else if columns, ok := options["columns"]; ok {
		removeIndexSql = ddl.getIndexNameByColumns(table, columns)
	} else if column, ok := options["column"]; ok {
		removeIndexSql = ddl.getIndexNameByColumns(table, column)
	} else {
		removeIndexSql = ""
	}

	if removeIndexSql != "" {
		if concurrently, ok := options["concurrently"]; ok {
			removeIndexSql = ddl.checkRemoveIndexConcurrently(concurrently) + removeIndexSql
		}

		if option, ok := options["option"]; ok {
			removeIndexSql = removeIndexSql + ddl.checkRemoveIndexOption(option)
		}
	}

	return "DROP INDEX " + removeIndexSql + ";"
}

func (ddl Ddl) AddColumn(table string, column col.Info) string {
	var addColumn string

	newColumn, newIndex, newForeignKey := ddl.generateColumn(table, column)
	addColumn = "ALTER TABLE " + table + " ADD COLUMN " + newColumn

	if newIndex != "" {
		addColumn = addColumn + newIndex
	}

	if newForeignKey != "" {
		addColumn = addColumn + newForeignKey
	}

	return addColumn + ";"
}

func (ddl Ddl) RenameColumn(table string, column string, newColumnName string) string {
	return "ALTER TABLE " + table + " RENAME " + column + " TO " + newColumnName + ";"
}

func (ddl Ddl) ChangeColumnType(table string, column string, newColumnType string) string {
	return "ALTER TABLE " + table + " ALTER COLUMN " + column + " TYPE " + ddl.translateToSqlType(newColumnType) + ";"
}

func (ddl Ddl) ChangeColumnNull(table string, column string, isNull bool) string {
	var nullSql string

	if isNull {
		nullSql = "DROP NOT NULL"
	} else {
		nullSql = "SET NOT NULL"
	}

	return "ALTER TABLE " + table + " ALTER COLUMN " + column + " " + nullSql + ";"
}

func (ddl Ddl) ChangeColumnDefault(table string, column string, defaultValue interface{}) string {
	var defaultSql string

	if defaultValue == nil {
		defaultSql = "DROP DEFAULT"
	} else {
		columnType := ddl.getColumnType(table, column)
		defaultSql = " SET " + ddl.checkDefault(columnType, defaultValue)
	}

	return "ALTER TABLE " + table + " ALTER COLUMN " + column + " " + defaultSql + ";"
}

func (ddl Ddl) RemoveColumn(table string, column string) string {
	return "ALTER TABLE " + table + " DROP COLUMN " + column + ";"
}

// options are: "column", "on_delete", "on_update", "name" and "primary_key" (all are strings)
// "on_delete" and "on_update" available values are "nullify", "cascade" and "restrict"
func (ddl Ddl) AddForeignKey(fromTable string, toTable string, options map[string]string) string {
	var addForeignKeySql string

	addForeignKeySql = "ALTER TABLE " + fromTable
	addForeignKeySql = addForeignKeySql + " ADD CONSTRAINT " + ddl.checkForeignKeyName(options["name"])
	addForeignKeySql = addForeignKeySql + " FOREIGN KEY (" + ddl.checkForeignKeyColumn(toTable, options["column"]) + ")"
	addForeignKeySql = addForeignKeySql + " REFERENCES " + toTable + "(" + ddl.checkForeignKeyPrimaryKey(options["primary_key"]) + ")"

	if onDelete, ok := options["on_delete"]; ok {
		addForeignKeySql = addForeignKeySql + ddl.checkForeignKeyOnDeleteOrUpdate("DELETE", onDelete)
	}

	if onUpdate, ok := options["on_update"]; ok {
		addForeignKeySql = addForeignKeySql + ddl.checkForeignKeyOnDeleteOrUpdate("UPDATE", onUpdate)
	}

	addForeignKeySql = addForeignKeySql + ";"

	return addForeignKeySql
}

// options are: "to_table", "column" and "name"
func (ddl Ddl) RemoveForeignKey(table string, options map[string]string) string {
	type Result struct {
		ConstraintName    string
		TableName         string
		ColumnName        string
		ForeignTableName  string
		ForeignColumnName string
	}
	var result Result
	var fkName string
	var sql string
	var queryValue string
	var regexpColumn = regexp.MustCompile(` + "`" + `XXXXXX` + "`" + `)

	sql = ` + "`" + `SELECT 
            tc.constraint_name, tc.table_name, kcu.column_name, 
            ccu.table_name AS foreign_table_name, 
            ccu.column_name AS foreign_column_name
         FROM
            information_schema.table_constraints AS tc
            JOIN information_schema.key_column_usage AS kcu
              ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
            JOIN information_schema.constraint_column_usage AS ccu
              ON ccu.constraint_name = tc.constraint_name AND ccu.table_schema = tc.table_schema
         WHERE tc.constraint_type = ? AND tc.table_name = ? AND XXXXXX = ?;` + "`" + `

	if _, ok := options["name"]; ok {
		fkName = options["name"]
	} else {
		if _, ok := options["to_table"]; ok {
			sql = regexpColumn.ReplaceAllString(sql, "foreign_table_name")
			queryValue = options["to_table"]
		} else if _, ok := options["column"]; ok {
			sql = regexpColumn.ReplaceAllString(sql, "kcu.column_name")
			queryValue = options["column"]
		} else {
			sql = ""
			queryValue = ""
		}

		model.Db.Raw(sql, "FOREIGN KEY", table, queryValue).Scan(&result)

		if model.Db.Error == nil {
			fkName = result.ConstraintName
		}
	}

	return "ALTER TABLE " + table + " DROP CONSTRAINT " + fkName + ";"
}

func (ddl Ddl) translateToSqlType(inputType string) string {
	var regexpText = regexp.MustCompile(` + "`" + `(?i)text` + "`" + `)
	var regexpString = regexp.MustCompile(` + "`" + `(?i)string` + "`" + `)
	var regexpDecimal = regexp.MustCompile(` + "`" + `(?i)(float|double|decimal|numeric)` + "`" + `)
	var regexpSmallInt = regexp.MustCompile(` + "`" + `(?i)smallint` + "`" + `)
	var regexpBigInt = regexp.MustCompile(` + "`" + `(?i)bigint` + "`" + `)
	var regexpInteger = regexp.MustCompile(` + "`" + `(?i)(int|integer|uint)` + "`" + `)
	var regexpDatetime = regexp.MustCompile(` + "`" + `(?i)datetime` + "`" + `)
	var regexpBoolean = regexp.MustCompile(` + "`" + `(?i)bool` + "`" + `)
	var regexpReference = regexp.MustCompile(` + "`" + `(?i)reference` + "`" + `)

	if regexpText.MatchString(inputType) {
		return "TEXT"
	} else if regexpString.MatchString(inputType) {
		return "VARCHAR"
	} else if regexpDecimal.MatchString(inputType) {
		return "NUMERIC"
	} else if regexpSmallInt.MatchString(inputType) {
		return "SMALLINT"
	} else if regexpBigInt.MatchString(inputType) {
		return "BIGINT"
	} else if regexpInteger.MatchString(inputType) {
		return "INT"
	} else if regexpDatetime.MatchString(inputType) {
		return "TIMESTAMP"
	} else if regexpBoolean.MatchString(inputType) {
		return "BOOLEAN"
	} else if regexpReference.MatchString(inputType) {
		return "BIGINT"
	} else {
		return "VARCHAR"
	}
}

func (ddl Ddl) checkLimitForVarchar(inputType string, limit interface{}) string {
	var rLimit string
	var regexpNotNumbers = regexp.MustCompile(` + "`" + `[^\d]` + "`" + `)

	if inputType == "VARCHAR" {
		switch v := limit.(type) {
		case int:
			rLimit = "(" + strconv.Itoa(v) + ")"
		case int64:
			rLimit = "(" + strconv.FormatInt(int64(v), 10) + ")"
		case string:
			rLimit = "(" + regexpNotNumbers.ReplaceAllString(v, "") + ")"
		default:
			rLimit = ""
		}
	} else {
		rLimit = ""
	}

	return rLimit
}

func (ddl Ddl) checkNull(isNull interface{}) string {
	var rNull string

	switch v := isNull.(type) {
	case bool:
		if !v {
			rNull = " NOT NULL"
		}
	case int:
		if v == 0 {
			rNull = " NOT NULL"
		}
	case string:
		isNull = strings.ToUpper(v)
		if isNull != "T" || isNull != "TRUE" {
			rNull = " NOT NULL"
		}
	default:
		rNull = ""
	}

	return rNull
}

func (ddl Ddl) checkDefault(inputType string, defaultValue interface{}) string {
	var rDefault string
	var regexpIsCharGroupType = regexp.MustCompile(` + "`" + `(?i)CHAR` + "`" + `)

	switch v := defaultValue.(type) {
	case bool:
		if v {
			rDefault = "TRUE"
		} else {
			rDefault = "FALSE"
		}
	case int:
		rDefault = strconv.Itoa(v)
	case int64:
		rDefault = strconv.FormatInt(v, 10)
	case float64:
		rDefault = strconv.FormatFloat(v, 'G', -1, 64)
	case string:
		rDefault = v
	default:
		rDefault = ""
	}

	if regexpIsCharGroupType.MatchString(inputType) {
		rDefault = "'" + rDefault + "'"
	}

	rDefault = " DEFAULT " + rDefault

	return rDefault
}

func (ddl Ddl) checkUnique(isUnique interface{}) string {
	var rUnique string

	switch v := isUnique.(type) {
	case bool:
		if v {
			rUnique = " UNIQUE"
		}
	case int:
		if v > 0 {
			rUnique = " UNIQUE"
		}
	case string:
		isUnique = strings.ToUpper(v)
		if isUnique == "T" || isUnique == "TRUE" {
			rUnique = " UNIQUE"
		}
	default:
		rUnique = ""
	}

	return rUnique
}

func (ddl Ddl) checkIndex(table string, column string, isIndexed interface{}) string {
	var rIndex string

	switch v := isIndexed.(type) {
	case bool:
		if v {
			rIndex = ddl.AddIndex(table, []string{column}, col.Index{Name: "", Unique: false})
		}
	case int:
		if v > 0 {
			rIndex = ddl.AddIndex(table, []string{column}, col.Index{Name: "", Unique: false})
		}
	case string:
		isIndexed = strings.ToUpper(v)
		if isIndexed == "T" || isIndexed == "TRUE" {
			rIndex = ddl.AddIndex(table, []string{column}, col.Index{Name: "", Unique: false})
		}
	default:
		rIndex = ""
	}

	return rIndex
}

func (ddl Ddl) checkForeignKey(fromTable string, toTable string, isForeignKey interface{}) string {
	var rForeignKey string

	switch v := isForeignKey.(type) {
	case bool:
		if v {
			rForeignKey = ddl.AddForeignKey(fromTable, toTable, make(map[string]string))
		}
	case int:
		if v > 0 {
			rForeignKey = ddl.AddForeignKey(fromTable, toTable, make(map[string]string))
		}
	case string:
		isForeignKey = strings.ToUpper(v)
		if isForeignKey == "T" || isForeignKey == "TRUE" {
			rForeignKey = ddl.AddForeignKey(fromTable, toTable, make(map[string]string))
		}
	default:
		rForeignKey = ""
	}

	return rForeignKey
}

func (ddl Ddl) checkPrecionAndScale(inputType string, precision interface{}, scale interface{}) string {
	var regexpForValidNumber = regexp.MustCompile(` + "`" + `\A\d+\z` + "`" + `)
	var regexpForOnlyZeros = regexp.MustCompile(` + "`" + `\A0+\z` + "`" + `)
	var rPrecision string

	if inputType != "NUMERIC" {
		return ""
	}

	switch v := precision.(type) {
	case int:
		if v > 0 {
			rPrecision = strconv.Itoa(v)
		}
	case string:
		if regexpForValidNumber.MatchString(v) && !regexpForOnlyZeros.MatchString(v) {
			rPrecision = v
		}
	default:
		rPrecision = ""
	}

	if rPrecision != "" {
		switch v := scale.(type) {
		case int:
			if v > 0 {
				rPrecision = rPrecision + ", " + strconv.Itoa(v)
			}
		case string:
			if regexpForValidNumber.MatchString(v) && !regexpForOnlyZeros.MatchString(v) {
				rPrecision = rPrecision + ", " + v
			}
		default:
			rPrecision = rPrecision
		}

		return "(" + rPrecision + ")"
	} else {
		return ""
	}
}

func (ddl Ddl) generateColumn(table string, column col.Info) (string, string, string) {
	var newColumn string
	var index string
	var foreignKey string

	if column.Type == "REFERENCES" {
		referenceTable := inflection.Plural(column.Name)
		column.Name = column.Name + "_id"
		column.Options["index"] = true
		newColumn = column.Name + " " + strings.ToUpper(ddl.getColumnType(referenceTable, "id"))

		if isForeignKey, ok := column.Options["foreign_key"]; ok {
			foreignKey = ddl.checkForeignKey(table, referenceTable, isForeignKey)
		}

		if foreignKey == "" {
			newColumn = newColumn + " REFERENCES " + referenceTable + "(id)"
		}
	} else {
		newColumn = column.Name + " " + column.Type

		if precision, ok := column.Options["precision"]; ok {
			newColumn = newColumn + ddl.checkPrecionAndScale(column.Type, precision, column.Options["scale"])
		}

		if limit, ok := column.Options["limit"]; ok {
			newColumn = newColumn + ddl.checkLimitForVarchar(column.Type, limit)
		}

		if isNull, ok := column.Options["null"]; ok {
			newColumn = newColumn + ddl.checkNull(isNull)
		}

		if defaultValue, ok := column.Options["default"]; ok {
			newColumn = newColumn + ddl.checkDefault(column.Type, defaultValue)
		}

		if isIndexed, ok := column.Options["index"]; ok {
			index = ddl.checkIndex(table, column.Name, isIndexed)
		}

		if unique, ok := column.Options["unique"]; ok {
			newColumn = newColumn + ddl.checkUnique(unique)
		}
	}

	return newColumn, index, foreignKey
}

func (ddl Ddl) getColumnType(table string, column string) string {
	type Result struct {
		DataType string
	}
	var result Result
	var sql string

	sql = "SELECT data_type FROM information_schema.columns WHERE table_name = ? AND column_name = ?;"

	model.Db.Raw(sql, table, column).Scan(&result)

	if model.Db.Error == nil {
		return result.DataType
	} else {
		return ""
	}
}

func (ddl Ddl) getIndexNameByColumns(table string, columns interface{}) string {
	type Result struct {
		Indexname string
	}
	var result Result
	var columnsSearch string
	var sql string

	switch v := columns.(type) {
	case string:
		columnsSearch = "(" + strings.Trim(v, "()") + ")"
	case []string:
		columnsSearch = "(" + strings.Join(v, ",") + ")"
	default:
		return ""
	}

	sql = "SELECT indexname FROM pg_indexes WHERE tablename = ? AND indexdef LIKE ? LIMIT 1;"

	model.Db.Raw(sql, table, ` + "`" + `%` + "`" + `+columnsSearch+` + "`" + `%` + "`" + `).Scan(&result)

	if model.Db.Error == nil {
		return result.Indexname
	} else {
		return ""
	}
}

func (ddl Ddl) checkRemoveIndexConcurrently(isConcurrently interface{}) string {
	var rConcurrently string

	switch v := isConcurrently.(type) {
	case bool:
		if v {
			rConcurrently = " CONCURRENTLY "
		}
	case int:
		if v > 0 {
			rConcurrently = " CONCURRENTLY "
		}
	case string:
		isConcurrently = strings.ToUpper(v)
		if isConcurrently == "T" || isConcurrently == "TRUE" {
			rConcurrently = " CONCURRENTLY "
		}
	default:
		rConcurrently = ""
	}

	return rConcurrently
}

func (ddl Ddl) checkRemoveIndexOption(tOption interface{}) string {
	var rOption string

	switch v := tOption.(type) {
	case string:
		tOption = strings.ToUpper(v)
		if tOption == "CASCADE" {
			rOption = " CASCADE "
		} else if tOption == "RESTRICT" {
			rOption = " RESTRICT "
		}
	default:
		rOption = ""
	}

	return rOption
}

func (ddl Ddl) checkForeignKeyName(fkName interface{}) string {
	var rFkName string
	var regexpInvalidFkName = regexp.MustCompile(` + "`" + `[^\w]` + "`" + `)

	switch v := fkName.(type) {
	case string:
		if fkName == "" || regexpInvalidFkName.MatchString(v) {
			rFkName = ""
		} else {
			rFkName = v
		}
	default:
		rFkName = ""
	}

	if rFkName == "" {
		rFkName = "fk_wheel_" + crypto.RandString(10)
	}

	return rFkName
}

func (ddl Ddl) checkForeignKeyColumn(table string, fkColumn interface{}) string {
	var rFkColumn string

	switch v := fkColumn.(type) {
	case string:
		rFkColumn = v
	default:
		rFkColumn = ""
	}

	if rFkColumn == "" && table != "" {
		rFkColumn = inflection.Singular(table) + "_id"
	}

	return rFkColumn
}

func (ddl Ddl) checkForeignKeyPrimaryKey(primaryKey interface{}) string {
	var rPrimaryKeyName string
	var regexpInvalidPrimaryKeyName = regexp.MustCompile(` + "`" + `[^\w]` + "`" + `)

	switch v := primaryKey.(type) {
	case string:
		if primaryKey == "" || regexpInvalidPrimaryKeyName.MatchString(v) {
			rPrimaryKeyName = "id"
		} else {
			rPrimaryKeyName = v
		}
	default:
		rPrimaryKeyName = "id"
	}

	return rPrimaryKeyName
}

func (ddl Ddl) checkForeignKeyOnDeleteOrUpdate(trigger string, action string) string {
	var rTrigger string
	var regexpValidTrigger = regexp.MustCompile(` + "`" + `\ADELETE|UPDATE\z` + "`" + `)
	var regexpValidAction = regexp.MustCompile(` + "`" + `\ACASCADE|NULLIFY|RESTRICT\z` + "`" + `)

	trigger = strings.ToUpper(trigger)
	action = strings.ToUpper(action)

	if regexpValidTrigger.MatchString(trigger) && regexpValidAction.MatchString(action) {
		rTrigger = " ON " + trigger + " " + action
	} else {
		rTrigger = ""
	}

	return rTrigger
}`
