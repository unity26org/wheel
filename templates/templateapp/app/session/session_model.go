package session

var ModelPath = []string{"app", "session", "session_model.go"}

var ModelContent = `package session

import (
	"errors"
	"time"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/db/entities"
)

const NotFound = "session was not found"

func Find(id interface{}) (entities.Session, error) {
	var session entities.Session
	var err error

	model.Db.First(&session, id)
	if model.Db.NewRecord(session) {
		err = errors.New(NotFound)
	}

	return session, err
}

func IsValid(session *entities.Session) (bool, []error) {
	return true, []error{}
}

func Update(session *entities.Session) (bool, []error) {
	var newValue, currentValue interface{}
	var valid bool
	var errs []error

	mapUpdate := make(map[string]interface{})

	currentSession, findErr := Find(session.ID)
	if findErr != nil {
		return false, []error{findErr}
	}

	valid, errs = IsValid(session)

	if valid {
		columns := model.ColumnsFromTable(session, false)
		for _, column := range columns {
			newValue, _ = model.GetColumnValue(session, column)
			currentValue, _ = model.GetColumnValue(currentSession, column)

			if newValue != currentValue {
				mapUpdate[column] = newValue
			}
		}

		if len(mapUpdate) > 0 {
			model.Db.Model(&session).Updates(mapUpdate)
		}

	}

	return valid, errs
}

func Create(session *entities.Session) (bool, []error) {
	valid, errs := IsValid(session)
	if valid && model.Db.NewRecord(session) {
		model.Db.Create(&session)

		if model.Db.NewRecord(session) {
			errs = append(errs, errors.New("database error"))
			return false, errs
		}
	}

	return valid, errs
}

func Save(session *entities.Session) (bool, []error) {
	if model.Db.NewRecord(session) {
		return Create(session)
	} else {
		return Update(session)
	}
}

func Destroy(session *entities.Session) bool {
	if model.Db.NewRecord(session) {
		return false
	} else {
		model.Db.Delete(&session)
		return true
	}
}

func FindByJti(jti string) (entities.Session, error) {
	var session entities.Session
	var err error

	model.Db.Where("jti = ?", jti).First(&session)
	if model.Db.NewRecord(session) {
		session = entities.Session{}
		err = errors.New(NotFound)
	}

	return session, err
}

func Deactivate(session *entities.Session) (bool, []error) {
	session.Active = false
	return Save(session)
}

func IncrementStats(session *entities.Session) (bool, []error) {
	t := time.Now()
	session.LastRequestAt = &t
	session.Requests = session.Requests + 1
	return Save(session)
}`
