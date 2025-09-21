package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
    JsonLogger *logrus.Logger
    TermLogger *logrus.Logger
)

func InitLogger() {

    JsonLogger = logrus.New()
    JsonLogger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
    })
    file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        logrus.Fatalf("Failed to open log file: %v", err)
    }
    JsonLogger.SetOutput(file)

 
    TermLogger = logrus.New()
    TermLogger.SetFormatter(&logrus.TextFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
        FullTimestamp:   true,
        ForceColors:     true,
        DisableLevelTruncation: true,
        PadLevelText:    true,
    })
    TermLogger.SetOutput(os.Stdout)
}