package worker

import (
	"math/rand"
	"time"

	"go-task-service/internal/logger"
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
}

func NewProcessor(repo *repository.TaskRepo, workers, countOfTryings int) *Processor {
	p := &Processor{
		Repo:           repo,
		Jobs:           make(chan TaskJob, 100),
		CountOfTryings: countOfTryings,
	}
	for i := 0; i < workers; i++ {
		go p.worker(i)
	}
	return p
}

func (p *Processor) worker(id int) {
	workerLogger := logger.Log.WithFields(logrus.Fields{"component": "worker", "id": id})
	for job := range p.Jobs {
		logger.Log.Infof("[Worker %d] Task %d started", id, job.ID)

		if job.CountOfTryings >= p.CountOfTryings {
			p.Repo.UpdateStatusTx(job.ID, string(model.StatusFailed), nil, job.CountOfTryings)
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
		workerLogger.
			WithFields(logrus.Fields{"jobId": job.ID, "status": finalStatus}).
			Infof("[Worker] Task completed")
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
