package entities

var UserPath = []string{"db", "entities", "user_entity.go"}

var UserContent = `package entities

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name                string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	Email               string     ` + "`" + `gorm:"type:varchar(255);unique_index"` + "`" + `
	Admin               bool       ` + "`" + `gorm:"default:false"` + "`" + `
	Password            string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	ResetPasswordToken  string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	ResetPasswordSentAt *time.Time ` + "`" + `gorm:"default:null"` + "`" + `
	Locale              string     ` + "`" + `gorm:"type:varchar(255);default:'en'"` + "`" + `
	Sessions            []Session
}`
