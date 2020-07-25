package usertemplate

var ModelPath = []string{"app", "models", "user", "user_model.go"}

var ModelContent = `package user

import (
	"errors"
	"regexp"
	"strings"
	"time"
	"{{ .AppRepository }}/app/entities"
	"{{ .AppRepository }}/commons/app/model"
	"{{ .AppRepository }}/commons/config"
	"{{ .AppRepository }}/commons/crypto"
)

const NotFound = "user was not found"

var Current entities.User

func Find(id interface{}) (entities.User, error) {
	var user entities.User
	var err error

	model.Db.First(&user, id, "deleted_at IS NULL")
	if model.Db.NewRecord(user) {
		err = errors.New(NotFound)
	}

	return user, err
}

func FindAll() []entities.User {
	var users []entities.User

	model.Db.Order("name").Find(&users, "deleted_at IS NULL")

	return users
}

func IsValid(user *entities.User) (bool, []error) {
	var count int
	var errs []error
	var validEmail = regexp.MustCompile(` + "`" + `\A[^@]+@([^@\.]+\.)+[^@\.]+\z` + "`" + `)

	if len(user.Name) == 0 {
		errs = append(errs, errors.New("name can't be blank"))
	} else if len(user.Name) > 255 {
		errs = append(errs, errors.New("name is too long"))
	}

	if len(user.Email) == 0 {
		errs = append(errs, errors.New("email can't be blank"))
	} else if len(user.Email) > 255 {
		errs = append(errs, errors.New("email is too long"))
	} else if !validEmail.MatchString(user.Email) {
		errs = append(errs, errors.New("email is invalid"))
	} else if model.Db.Model(&entities.User{}).Where("id <> ? AND email = ? AND deleted_at IS NULL", user.ID, user.Email).Count(&count); count > 0 {
		errs = append(errs, errors.New("email has already been taken"))
	}

	if len(user.Password) < 8 {
		errs = append(errs, errors.New("password is too short minimum is 8 characters"))
	} else if len(user.Password) > 255 {
		errs = append(errs, errors.New("password is too long"))
	}

	if !isLocaleValid(user.Locale) {
		errs = append(errs, errors.New("locale is invalid"))
	}

	return (len(errs) == 0), errs
}

func Update(user *entities.User) (bool, []error) {
	var newValue, currentValue interface{}
	var valid bool
	var errs []error

	mapUpdate := make(map[string]interface{})

	currentUser, findErr := Find(user.ID)
	if findErr != nil {
		return false, []error{findErr}
	}

	valid, errs = IsValid(user)

	if valid {
		columns := model.ColumnsFromTable(user, false)
		for _, column := range columns {
			newValue, _ = model.GetColumnValue(user, column)
			currentValue, _ = model.GetColumnValue(currentUser, column)

			if newValue != currentValue {
				mapUpdate[column] = newValue

				if column == "password" {
					mapUpdate[column] = crypto.SetPassword(mapUpdate[column].(string))
				}
			}
		}

		if len(mapUpdate) > 0 {
			model.Db.Model(&user).Updates(mapUpdate)
		}

	}

	return valid, errs
}

func Create(user *entities.User) (bool, []error) {
	valid, errs := IsValid(user)
	if valid && model.Db.NewRecord(user) {
		user.Password = crypto.SetPassword(user.Password)

		model.Db.Create(&user)

		if model.Db.NewRecord(user) {
			errs = append(errs, errors.New("database error"))
			return false, errs
		}
	}

	return valid, errs
}

func Save(user *entities.User) (bool, []error) {
	if model.Db.NewRecord(user) {
		return Create(user)
	} else {
		return Update(user)
	}
}

func Destroy(user *entities.User) bool {
	if model.Db.NewRecord(user) {
		return false
	} else {
		model.Db.Delete(&user)
		return true
	}
}

func FindByEmail(email string) (entities.User, error) {
	var user entities.User
	var err error

	model.Db.Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if model.Db.NewRecord(user) {
		user = entities.User{}
		err = errors.New(NotFound)
	}

	return user, err
}

func FindByResetPasswordToken(token string) (entities.User, error) {
	var user entities.User
	var err error

	enconded_token := crypto.EncryptText(token, config.App.SecretKey)
	two_days_ago := time.Now().Add(time.Second * time.Duration(config.App.ResetPasswordExpirationSeconds) * (-1))

	model.Db.Where("reset_password_token = ? AND reset_password_sent_at >= ? AND deleted_at IS NULL", enconded_token, two_days_ago).First(&user)
	if model.Db.NewRecord(user) {
		user = entities.User{}
		err = errors.New(NotFound)
	}

	return user, err
}

func Paginate(criteria map[string]string, order, page, perPage string) ([]entities.User, int, int, int) {
	var users []entities.User
	var user entities.User

	q := model.Query{Db: model.Db, Table: &user}
	q.SearchEngine(criteria)
	q.Ordering(order)
	currentPage, totalPages, totalEntries := q.Pagination(page, perPage)

	q.Db.Find(&users, "deleted_at IS NULL")

	return users, currentPage, totalPages, totalEntries
}

func Authenticate(email string, password string) (entities.User, error) {
	user, err := FindByEmail(email)

	if model.Db.NewRecord(user) || !crypto.CheckPassword(password, user.Password) {
		user = entities.User{}
		err = errors.New("invalid credentials")
	}

	return user, err
}

func IsNil(user *entities.User) bool {
	return model.Db.NewRecord(user)
}

func Exists(user *entities.User) bool {
	return !IsNil(user)
}

func SetCurrent(id interface{}) error {
	var err error
	Current, err = Find(id)

	return err
}

func IdExists(id interface{}) bool {
	_, err := Find(id)

	return (err == nil)
}

func SetRecovery(user *entities.User) (string, []error) {
	token := crypto.RandString(20)

	if model.Db.NewRecord(user) {
		return "", []error{errors.New(NotFound)}
	} else {
		t := time.Now()
		user.ResetPasswordSentAt = &t
		user.ResetPasswordToken = crypto.EncryptText(token, config.App.SecretKey)

		valid, errs := Save(user)

		if valid {
			return token, errs
		} else {
			return "", errs
		}
	}
}

func ClearRecovery(user *entities.User) (bool, []error) {
	if model.Db.NewRecord(user) {
		return false, []error{errors.New(NotFound)}
	} else {
		user.ResetPasswordToken = ""
		user.ResetPasswordSentAt = nil
		valid, errs := Save(user)

		return valid, errs
	}
}

func FirstName(user *entities.User) string {
	return strings.Split(user.Name, " ")[0]
}

// local methods

func isLocaleValid(locale string) bool {
	locales := config.App.Locales

	for _, a := range locales {
		if a == locale {
			return true
		}
	}

	return false
}`
