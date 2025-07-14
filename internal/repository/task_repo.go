package repository

import (
	"fmt"

	"go-task-service/config"
	"go-task-service/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db}
}

func (r *TaskRepo) GetNewTasks(limit int) []model.Task {
	var tasks []model.Task
	// Debug тоьлко для отладки
	r.db.Debug().Where("status = ?", "NEW").Limit(limit).Find(&tasks)
	return tasks
}

func (r *TaskRepo) UpdateStatusTx(id uint, status string, errorMessage *string, countOfTryings int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Task{}).
			Debug().
			Where("id = ?", id).
			Updates(map[string]interface{}{
				"error_message":    errorMessage,
				"status":           status,
				"count_of_tryings": countOfTryings,
			}).
			Error
	})
}

func (r *TaskRepo) GetTasks(page, pageSize int, status string) ([]model.Task, int64, error) {
	offset := (page - 1) * pageSize
	var tasks []model.Task
	var total int64

	db := r.db.Model(&model.Task{}).Limit(pageSize).Offset(offset)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	if err := db.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	// считаем общее количество
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.DBHost, cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBName, cfg.DB.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
