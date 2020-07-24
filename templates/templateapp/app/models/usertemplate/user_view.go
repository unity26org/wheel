package usertemplate

var ViewPath = []string{"app", "models", "user", "user_view.go"}

var ViewContent = `package user

import (
	"time"
	"{{ .AppRepository }}/app/entities"
	"{{ .AppRepository }}/commons/app/view"
)

type PaginationJson struct {
	Pagination view.MainPagination ` + "`" + `json:"pagination"` + "`" + `
	Users      []Json              ` + "`" + `json:"users"` + "`" + `
}

type SuccessfullySavedJson struct {
	SystemMessage view.SystemMessage ` + "`" + `json:"system_message"` + "`" + `
	User          Json               ` + "`" + `json:"user"` + "`" + `
}

type Json struct {
	ID        uint64     ` + "`" + `json:"id"` + "`" + `
	Name      string     ` + "`" + `json:"name"` + "`" + `
	Email     string     ` + "`" + `json:"email"` + "`" + `
	Admin     bool       ` + "`" + `json:"admin"` + "`" + `
	Locale    string     ` + "`" + `json:"locale"` + "`" + `
	CreatedAt time.Time  ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time  ` + "`" + `json:"updated_at"` + "`" + `
	DeletedAt *time.Time ` + "`" + `json:"deleted_at"` + "`" + `
}

func SetJson(user entities.User) Json {
	return Json{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Admin:     user.Admin,
		Locale:    user.Locale,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}`
