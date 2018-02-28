package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/GrupaPracuj/iislog-prometheus-exporter/config"

	"gopkg.in/natefinch/lumberjack.v2"
)

const logFileName string = "log.txt"
const logFileDir string = "logs\\"
const logRotateEveryMb int = 10
const filesNumber int = 0
const maxFilesAge int = 0

func Init(cfg *config.Config, isDebug bool) (logger *log.Logger) {
	var outputFile string
	var rotateEveryMb int
	var logWriter io.Writer

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(ex)

	if cfg.Logger != (config.LoggerConfig{}) && len(cfg.Logger.OutputDir) > 0 {
		outputFile = cfg.Logger.OutputDir + "\\" + logFileName
	} else {
		outputFile = exeDir + "\\" + logFileDir + logFileName
	}
	if cfg.Logger != (config.LoggerConfig{}) && cfg.Logger.RotateEveryMb > 0 {
		rotateEveryMb = cfg.Logger.RotateEveryMb
	} else {
		rotateEveryMb = logRotateEveryMb
	}

	logWriter = &lumberjack.Logger{
		Filename:   outputFile,
		MaxSize:    rotateEveryMb,
		MaxBackups: filesNumber,
		MaxAge:     maxFilesAge, //days
	}

	if isDebug {
		logger = log.New(io.MultiWriter(os.Stdout, logWriter), "", log.Ldate|log.Ltime)
	} else {
		logger = log.New(logWriter, "", log.Ldate|log.Ltime)
	}
	return logger
}

func Info(logger *log.Logger, message string) {
	logger.Printf("INFO %s\r\n", message)
}

func Error(logger *log.Logger, message string, err error) {
	logger.Printf("ERROR %s %s\r\n", message, err)
}
