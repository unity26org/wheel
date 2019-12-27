package templatecrud

var ModelContent = `package {{ .EntityName.LowerCase }}

import (
	"errors"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/db/entities"
)

const NotFound = "{{ .EntityName.SnakeCase }} was not found"

var Current entities.{{ .EntityName.CamelCase }}

func Find(id interface{}) (entities.{{ .EntityName.CamelCase }}, error) {
	var {{ .EntityName.LowerCamelCase }} entities.{{ .EntityName.CamelCase }}
	var err error

	model.Db.First(&{{ .EntityName.LowerCamelCase }}, id)
	if model.Db.NewRecord({{ .EntityName.LowerCamelCase }}) {
		err = errors.New(NotFound)
	}

	return {{ .EntityName.LowerCamelCase }}, err
}

func FindAll() []entities.{{ .EntityName.CamelCase }} {
	var {{ .EntityName.LowerCamelCasePlural }} []entities.{{ .EntityName.CamelCase }}

	model.Db.Find(&{{ .EntityName.LowerCamelCasePlural }})

	return {{ .EntityName.LowerCamelCasePlural }}
}

func IsValid({{ .EntityName.LowerCamelCase }} *entities.{{ .EntityName.CamelCase }}) (bool, []error) {
	var errs []error

	return (len(errs) == 0), errs
}

func Update({{ .EntityName.LowerCamelCase }} *entities.{{ .EntityName.CamelCase }}) (bool, []error) {
	var newValue, currentValue interface{}
	var valid bool
	var errs []error

	mapUpdate := make(map[string]interface{})

	current{{ .EntityName.CamelCase }}, findErr := Find({{ .EntityName.LowerCamelCase }}.ID)
	if findErr != nil {
		return false, []error{findErr}
	}

	valid, errs = IsValid({{ .EntityName.LowerCamelCase }})

	if valid {
		columns := model.ColumnsFromTable({{ .EntityName.LowerCamelCase }}, false)
		for _, column := range columns {
			newValue, _ = model.GetColumnValue({{ .EntityName.LowerCamelCase }}, column)
			currentValue, _ = model.GetColumnValue(current{{ .EntityName.CamelCase }}, column)

			if newValue != currentValue {
				mapUpdate[column] = newValue
			}
		}

		if len(mapUpdate) > 0 {
			model.Db.Model(&{{ .EntityName.LowerCamelCase }}).Updates(mapUpdate)
		}

	}

	return valid, errs
}

func Create({{ .EntityName.LowerCamelCase }} *entities.{{ .EntityName.CamelCase }}) (bool, []error) {
	valid, errs := IsValid({{ .EntityName.LowerCamelCase }})
	if valid && model.Db.NewRecord({{ .EntityName.LowerCamelCase }}) {
		model.Db.Create(&{{ .EntityName.LowerCamelCase }})

		if model.Db.NewRecord({{ .EntityName.LowerCamelCase }}) {
			errs = append(errs, errors.New("database error"))
			return false, errs
		}
	}

	return valid, errs
}

func Save({{ .EntityName.LowerCamelCase }} *entities.{{ .EntityName.CamelCase }}) (bool, []error) {
	if model.Db.NewRecord({{ .EntityName.LowerCamelCase }}) {
		return Create({{ .EntityName.LowerCamelCase }})
	} else {
		return Update({{ .EntityName.LowerCamelCase }})
	}
}

func Destroy({{ .EntityName.LowerCamelCase }} *entities.{{ .EntityName.CamelCase }}) bool {
	if model.Db.NewRecord({{ .EntityName.LowerCamelCase }}) {
		return false
	} else {
		model.Db.Delete(&{{ .EntityName.LowerCamelCase }})
		return true
	}
}

func Paginate(criteria map[string]string, order, page, perPage string) ([]entities.{{ .EntityName.CamelCase }}, int, int, int) {
	var {{ .EntityName.LowerCamelCasePlural }} []entities.{{ .EntityName.CamelCase }}
	var {{ .EntityName.LowerCamelCase }} entities.{{ .EntityName.CamelCase }}

	q := model.Query{Db: model.Db, Table: &{{ .EntityName.LowerCamelCase }}}
	q.SearchEngine(criteria)
	q.Ordering(order)
	currentPage, totalPages, totalEntries := q.Pagination(page, perPage)

	q.Db.Find(&{{ .EntityName.LowerCamelCasePlural }})

	return {{ .EntityName.LowerCamelCasePlural }}, currentPage, totalPages, totalEntries
}

func IsNil({{ .EntityName.LowerCamelCase }} *entities.{{ .EntityName.CamelCase }}) bool {
	return model.Db.NewRecord({{ .EntityName.LowerCamelCase }})
}

func Exists({{ .EntityName.LowerCamelCase }} *entities.{{ .EntityName.CamelCase }}) bool {
	return !IsNil({{ .EntityName.LowerCamelCase }})
}

func SetCurrent(id interface{}) error {
	var err error
	Current, err = Find(id)

	return err
}

func IdExists(id interface{}) bool {
	_, err := Find(id)

	return (err == nil)
}`
