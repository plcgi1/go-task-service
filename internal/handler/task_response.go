package handler

import "go-task-service/internal/model"

type TaskListResponse struct {
	Data     []model.Task `json:"data"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
}
