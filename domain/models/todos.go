package models

import (
	"syncova-todo/domain/enums"
	"time"
)

type TodosStatus struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid      string    `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	ProjectId uint      `gorm:"not null" json:"project_id"`
	Status    string    `gorm:"type:varchar(50);not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Labels struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid      string    `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	ProjectId uint      `gorm:"not null" json:"project_id"`
	Label     string    `gorm:"type:varchar(50);not null" json:"label"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Todos struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid        string         `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	UserId      string         `gorm:"type:varchar(255);not null" json:"user_id"`
	ProjectId   *uint          `gorm:"default:null" json:"project_id"`
	LabelId     *uint          `gorm:"default:null" json:"label_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	DueDate     string         `gorm:"type:date" json:"due_date"`
	StartDate   string         `gorm:"type:date" json:"start_date"`
	ParentID    *uint          `gorm:"default:null" json:"parent_id"`
	Children    []Todos        `gorm:"foreignKey:ParentID" json:"children"`
	Priority    enums.Priority `gorm:"type:enum('low','medium','high');default:'medium';not null" json:"priority"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
}

type TodosStatusTodos struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid          string    `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	TodosStatusID uint      `gorm:"not null" json:"todos_status_id"`
	TodosID       uint      `gorm:"not null" json:"todos_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TodosLabel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid      string    `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	LabelID   uint      `gorm:"not null" json:"label_id"`
	TodosID   uint      `gorm:"not null" json:"todos_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
