package entities

var UserPath = []string{"app", "entities", "user_entity.go"}

var UserContent = `package entities

import (
	"time"
)

type User struct {
	ID                  uint64     ` + "`" + `gorm:"primary_key"` + "`" + `
	Name                string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	Email               string     ` + "`" + `gorm:"type:varchar(255);unique_index"` + "`" + `
	Admin               bool       ` + "`" + `gorm:"default:false"` + "`" + `
	Password            string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	ResetPasswordToken  string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	ResetPasswordSentAt *time.Time ` + "`" + `gorm:"default:null"` + "`" + `
	Locale              string     ` + "`" + `gorm:"type:varchar(255);default:'en'"` + "`" + `
	Sessions            []Session
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}`
