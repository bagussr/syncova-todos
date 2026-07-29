package models

import "time"

type LabelsStatuses struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid      string    `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	ProjectId uint      `gorm:"not null" json:"project_id"`
	Label     string    `gorm:"type:varchar(50);not null" json:"label"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TodosLabel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Uuid      string    `gorm:"type:uuid;default:gen_random_uuid()" json:"uuid"`
	LabelID   uint      `gorm:"not null" json:"label_id"`
	TodosID   uint      `gorm:"not null" json:"todos_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
