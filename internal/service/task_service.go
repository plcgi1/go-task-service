package service

import (
	"go-task-service/config"
	"go-task-service/internal/repository"
	"go-task-service/internal/worker"
)

type TaskService struct {
	Repo      *repository.TaskRepo
	Processor *worker.Processor
	Config    *config.Config
}

func NewTaskService(repo *repository.TaskRepo, cfg *config.Config) *TaskService {
	return &TaskService{
		Repo:      repo,
		Processor: worker.NewProcessor(repo, 5, cfg.CountOfTryings),
	}
}

func (s *TaskService) ProcessTasks(limit, minDelay, maxDelay int, successRate float64) int {
	tasks := s.Repo.GetNewTasks(limit)

	if len(tasks) == 0 {
		return 0 // нет задач
	}
	for _, task := range tasks {
		s.Processor.Jobs <- worker.TaskJob{
			ID:             task.ID,
			MinDelay:       minDelay,
			MaxDelay:       maxDelay,
			SuccessRate:    successRate,
			CountOfTryings: task.CountOfTryings,
		}
	}
	return len(tasks)
}
