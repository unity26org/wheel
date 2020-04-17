package entities

var SessionPath = []string{"db", "entities", "session_entity.go"}

var SessionContent = `package entities

import (
	"time"
)

type Session struct {
	ID            uint64     ` + "`" + `gorm:"primary_key"` + "`" + `
	UserID        uint64     ` + "`" + `gorm:"index"` + "`" + `
	Jti           string     ` + "`" + `gorm:"type:varchar(255);unique_index"` + "`" + `
	App           string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	Requests      int        ` + "`" + `gorm:"not null;default:0"` + "`" + `
	ExpiresIn     int        ` + "`" + `gorm:"not null;default:0"` + "`" + `
	Address       string     ` + "`" + `gorm:"type:varchar(255)"` + "`" + `
	Active        bool       ` + "`" + `gorm:"default:true"` + "`" + `
	LastRequestAt *time.Time ` + "`" + `gorm:"default:null"` + "`" + `
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ExpiresAt     time.Time
}`
