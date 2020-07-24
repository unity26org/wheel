package myself

var ViewPath = []string{"app", "models", "myself", "myself_view.go"}

var ViewContent = `package myself

import (
	"time"
	"{{ .AppRepository }}/app/entities"
)

type Json struct {
	ID        uint64    ` + "`" + `json:"id"` + "`" + `
	Name      string    ` + "`" + `json:"name"` + "`" + `
	Email     string    ` + "`" + `json:"email"` + "`" + `
	Locale    string    ` + "`" + `json:"locale"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func SetJson(userMyself entities.User) Json {
	return Json{
		ID:        userMyself.ID,
		Name:      userMyself.Name,
		Email:     userMyself.Email,
		Locale:    userMyself.Locale,
		CreatedAt: userMyself.CreatedAt,
		UpdatedAt: userMyself.UpdatedAt,
	}
}`
