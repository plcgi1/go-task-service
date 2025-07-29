package worker

import (
	"context"
	"math/rand"
	"time"

	"go-task-service/internal/logger"
	"go-task-service/internal/metrics"
	"go-task-service/internal/model"
	"go-task-service/internal/repository"

	"github.com/sirupsen/logrus"
)

type TaskJob struct {
	ID             uint
	CountOfTryings int
	MinDelay       int
	MaxDelay       int
	SuccessRate    float64
}

type Processor struct {
	Repo           *repository.TaskRepo
	Jobs           chan TaskJob
	CountOfTryings int
	ctx            context.Context
}

func NewProcessor(ctx context.Context, repo *repository.TaskRepo, workers, countOfTryings int) *Processor {
	p := &Processor{
		Repo:           repo,
		Jobs:           make(chan TaskJob, 100),
		CountOfTryings: countOfTryings,
		ctx:            ctx,
	}
	for i := 0; i < workers; i++ {
		go p.worker(i)
	}
	return p
}

// TODO проверить код и использовать
//
//	func ProcessTask(ctx context.Context, repo *Repo, job TaskJob, workerID int) {
//	   start := time.Now()
//	   traceID := uuid.New().String()
//	   ctx = context.WithValue(ctx, "traceID", traceID)
//
//	   log := logging.WithContext(ctx).WithField("taskID", job.ID)
//
//	   err := repo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
//	       var task Task
//	       if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
//	           Where("id = ? AND status = ?", job.ID, "NEW").
//	           First(&task).Error; err != nil {
//	           metrics.TaskLockFailed.Inc()
//	           log.WithError(err).Warn("Cannot lock task")
//	           return err
//	       }
//
//	       if err := tx.Model(&task).Update("status", "PROCESSING").Error; err != nil {
//	           return err
//	       }
//
//	       select {
//	       case <-ctx.Done():
//	           log.Warn("Task cancelled during processing")
//	           metrics.TaskCancelled.Inc()
//	           return errors.New("shutdown")
//	       default:
//	           // Simulate processing
//	           delay := rand.Intn(job.MaxDelay-job.MinDelay+1) + job.MinDelay
//	           time.Sleep(time.Millisecond * time.Duration(delay))
//
//	           success := rand.Float64() < job.SuccessRate
//	           newStatus := "PROCESSED"
//	           if !success {
//	               newStatus = "NEW"
//	               metrics.TaskFailed.Inc()
//	           } else {
//	               metrics.TaskProcessed.Inc()
//	           }
//
//	           return tx.Model(&task).Update("status", newStatus).Error
//	       }
//	   })
//
//	   metrics.TaskDuration.Observe(float64(time.Since(start).Milliseconds()))
//	   if err != nil {
//	       log.WithError(err).Error("Task failed")
//	   }
//	}

func (p *Processor) worker(id int) {
	workerLogger := logger.
		WithContext(p.ctx).
		WithFields(logrus.Fields{"component": "worker", "workerID": id})
	for job := range p.Jobs {
		select {
		case <-p.ctx.Done():
			workerLogger.Warn("Task cancelled by shutdown")
		default:
			start := time.Now()

			ctx := context.WithValue(p.ctx, "traceID", job.ID)
			ctx = context.WithValue(ctx, "workerID", id)

			workerLogger.Infof("[Worker %d] Task %d started", id, job.ID)
			// проверяем количество допустимых повторных запусков одной задачи
			if job.CountOfTryings >= p.CountOfTryings {
				errorMessage := "Trying count exceeded"
				p.Repo.UpdateStatusTx(job.ID, string(model.StatusFailed), &errorMessage, job.CountOfTryings)
				workerLogger.
					WithFields(logrus.Fields{"targetCountOfTryings": p.CountOfTryings, "jobId": job.ID, "countOfTryings": job.CountOfTryings}).
					Warnf("Err with count of tryings - setting status to failed")
				continue
			}

			if job.MaxDelay > 0 {
				delay := rand.Intn(job.MaxDelay-job.MinDelay+1) + job.MinDelay
				time.Sleep(time.Millisecond * time.Duration(delay))
				workerLogger.
					WithFields(logrus.Fields{"delay": delay, "jobId": job.ID}).
					Warnf("Sleeping")
			}

			if err := p.Repo.UpdateStatusTx(job.ID, string(model.StatusProcessing), nil, job.CountOfTryings); err != nil {
				// тут и алерты накидать куданить в sentry или еще куда
				workerLogger.
					WithFields(logrus.Fields{"status": model.StatusProcessing, "jobId": job.ID}).
					WithError(err).
					Warnf("Error updating task")
				continue
			}

			workerLogger.
				WithFields(logrus.Fields{"jobId": job.ID, "status": model.StatusProcessing}).
				Infof("[Worker] Task processing")

			finalStatus, err := doSomeTaskWork(job, workerLogger)

			if err != nil {
				workerLogger.
					WithFields(logrus.Fields{"jobId": job.ID, "status": finalStatus}).
					WithError(err).
					Warnf("Error in execution for task")
				continue
			}

			err = p.Repo.UpdateStatusTx(job.ID, string(finalStatus), nil, job.CountOfTryings+1)
			if err != nil {
				workerLogger.
					WithFields(logrus.Fields{"jobId": job.ID, "status": finalStatus}).
					WithError(err).
					Warnf("Error in UpdateStatusTx for task")
				continue
			}
			duration := time.Since(start).Milliseconds()
			metrics.TaskDuration.Observe(float64(duration))

			workerLogger.
				WithFields(logrus.Fields{
					"jobId":       job.ID,
					"status":      finalStatus,
					"duration_ms": duration,
				}).
				Infof("[Worker] Task completed")

		}
	}
}

// Главный обработчик конкретной задачи
func doSomeTaskWork(job TaskJob, jobLogger *logrus.Entry) (model.TaskStatus, error) {
	// Инициализация генератора
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	min := 10
	max := 500
	// генерим случайное число - это время работы задачи
	num := r.Intn(max-min+1) + min

	jobLogger.
		WithFields(logrus.Fields{"jobId": job.ID}).
		Infof("[Worker] Do some task work")

	time.Sleep(time.Millisecond * time.Duration(num))

	finalStatus := model.StatusNew
	if rand.Float64() <= job.SuccessRate {
		finalStatus = model.StatusProcessed
	}

	return finalStatus, nil
}
