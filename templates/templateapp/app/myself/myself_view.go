package myself

var ViewPath = []string{"app", "myself", "myself_view.go"}

var ViewContent = `package myself

import (
	"time"
	"{{ .AppRepository }}/db/entities"
)

type Json struct {
	ID        uint      ` + "`" + `json:"id"` + "`" + `
	Name      string    ` + "`" + `json:"name"` + "`" + `
	Email     string    ` + "`" + `json:"email"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}

func SetJson(userMyself entities.User) Json {
	return Json{
		ID:        userMyself.ID,
		Name:      userMyself.Name,
		Email:     userMyself.Email,
		CreatedAt: userMyself.CreatedAt,
		UpdatedAt: userMyself.UpdatedAt,
	}
}`
