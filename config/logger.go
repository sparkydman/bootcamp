package config

import (
	"io"
	"os"
	"strconv"

	nested "github.com/antonfisher/nested-logrus-formatter"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog() {
	logger.SetLevel(getLoggerLevel(os.Getenv("LOG_LEVEL")))
	logger.SetReportCaller(true)
	logger.SetFormatter(&nested.Formatter{
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "2006-01-02 15:04:05",
		ShowFullLevel:   true,
		CallerFirst:     true,
	})

	writers := []io.Writer{os.Stdout}
	logToFile, err := strconv.ParseBool(os.Getenv("LOG_TO_FILE"))
	if err != nil {
		panic(err)
	}
	if logToFile {
		writers = append(writers, &lumberjack.Logger{
			Filename:   os.Getenv("LOG_FILE"),
			MaxSize:    500, //megabytes max size of log file
			MaxBackups: 20,  //number of files to retain
			LocalTime:  true,
			Compress:   true, //compress backups})
		})
	}
}

func getLoggerLevel(value string) logger.Level {
	switch value {
	case "DEBUG":
		return logger.DebugLevel
	case "TRACE":
		return logger.TraceLevel
	case "WARNING":
		return logger.WarnLevel
	case "ERROR":
		return logger.ErrorLevel
	case "PANIC":
		return logger.PanicLevel
	default:
		return logger.InfoLevel
	}
}
