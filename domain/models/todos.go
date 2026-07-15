package models

import (
	"syncova-todo/domain/enums"
	"time"
)

type TodosStatus struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid      string `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	ProjectId uint   `gorm:"not null" json:"project_id"`
	Status    string `gorm:"type:varchar(50);not null" json:"status"`
}

type Labels struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid      string `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	ProjectId uint   `gorm:"not null" json:"project_id"`
	Label     string `gorm:"type:varchar(50);not null" json:"label"`
}

type Todos struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid        string         `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	LabelId     *uint          `gorm:"default:null" json:"label_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	DueDate     time.Time      `gorm:"type:date" json:"due_date"`
	StartDate   time.Time      `gorm:"type:date" json:"start_date"`
	ParentID    *uint          `gorm:"default:null" json:"parent_id"`
	Children    []Todos        `gorm:"foreignKey:ParentID" json:"children"`
	Priority    enums.Priority `gorm:"type:enum('low','medium','high');default:'medium';not null" json:"priority"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type TodosStatusTodos struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid          string `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	TodosStatusID uint   `gorm:"not null" json:"todos_status_id"`
	TodosID       uint   `gorm:"not null" json:"todos_id"`
}
