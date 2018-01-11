package run

import (
	"log"
	"path/filepath"

	"github.com/marcsauter/single"
	"gopkg.in/natefinch/lumberjack.v2"
)

var errorLogger *log.Logger

func setLogOutput(logName string) *log.Logger {
	logger := log.New(nil, log.Prefix(), log.LstdFlags)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   filepath.Join(config.logDir, logName+".log"),
		MaxSize:    256,
		MaxBackups: 1,
		MaxAge:     365,
		Compress:   true,
	})
	return logger
}

func Daemonize() {
	// Allow only a single instance
	instance := single.New("psb_run_daemon")
	instance.Lock()
	defer instance.Unlock()
	// Set error logger
	errorLogger = setLogOutput("run_error")
}

func logErr(err error) {
	errorLogger.Println(err)
}
