package models

import (
	"syncova-todo/domain/enums"
	"time"
)

type Project struct {
	Id          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid        string       `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	Status      enums.Status `sql:"type:enum('not_started','in_progress','testing','completed','backlog');default:'not_started';not null" json:"status"`
	DueDate     time.Time    `gorm:"type:date" json:"due_date"`
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Description string       `gorm:"type:text" json:"description"`
}
