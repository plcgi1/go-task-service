package logger

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger
var instanceID string

func Init() {
	Log = logrus.New()
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
	instanceID = os.Getenv("INSTANCE_ID")
	if instanceID == "" {
		if host, err := os.Hostname(); err == nil {
			instanceID = host
		} else {
			instanceID = "unknown-instance"
		}
	}
}

func WithContext(ctx context.Context) *logrus.Entry {
	traceID, _ := ctx.Value("traceID").(string)
	workerID, _ := ctx.Value("workerID").(int)
	taskType, _ := ctx.Value("taskType").(string)
	errVal, _ := ctx.Value("error").(string)

	return Log.WithFields(logrus.Fields{
		"traceID":   traceID,
		"workerID":  workerID,
		"taskType":  taskType,
		"instance":  instanceID,
		"timestamp": time.Now().Format(time.RFC3339),
		"error":     errVal,
	})
}
