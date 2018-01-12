package run

import (
	"log"
	"path/filepath"
	"time"

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

func getPollrate(runConfigs []RunConfig) (result time.Duration) {
	if len(runConfigs) > 1 {
		// Euclid algo
		gcd := func(x, y int) int {
			for y != 0 {
				x, y = y, x%y
			}
			return x
		}
		// This is the result for only 2 integers
		rate := gcd(runConfigs[0].Rotations.Frequency, runConfigs[1].Rotations.Frequency)
		result = time.Duration(rate)
		// for loop in case there're more than 2 ints
		for j := 3; j <= len(runConfigs)-1; j++ {
			rate = gcd(rate, runConfigs[j].Rotations.Frequency)
			result = time.Duration(rate)
		}
	} else {
		if len(runConfigs) < 1 {
			return
		}
		result = time.Duration(runConfigs[0].Rotations.Frequency)
	}
	return
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
