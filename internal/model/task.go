package model

import "time"

type TaskStatus string

const (
	StatusNew        TaskStatus = "NEW"
	StatusProcessing TaskStatus = "PROCESSING"
	StatusProcessed  TaskStatus = "PROCESSED"
	StatusFailed     TaskStatus = "FAILED"
)

type Task struct {
	ID             uint   `gorm:"primaryKey"`
	Status         string `gorm:"type:varchar(20);default:NEW;index"`
	CountOfTryings int    `gorm:"column:count_of_tryings"`
	ErrorMessage   string `gorm:"type:varchar(255)"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (Task) TableName() string {
	return "task"
}
