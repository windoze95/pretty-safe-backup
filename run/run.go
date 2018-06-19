package run

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/marcsauter/single"
	"github.com/orange-lightsaber/pretty-safe-backup/rsync"
	"github.com/orange-lightsaber/pretty-safe-backup/ssh"
	"github.com/orange-lightsaber/psb-rotatord/rotator"
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

func GetTimeSinceLastRun(rcd rotator.RunConfigData, remoteRun ssh.SSH) (timeSince time.Duration, available bool, err error) {
	var str string
	remote := len(remoteRun.Host) > 0
	if remote {
		remoteRun.Cmd = fmt.Sprintf("psb-rotatorc lastrun -name %s -dir '%s'", rcd.Name, rcd.BackupDir)
		str, err = remoteRun.Run()
		if err != nil {
			if strings.Contains(err.Error(), "connection timed out") {
				// Don't return an error if the host is not available
				// psb will try again on the next poll
				err = nil
				return
			}
			return
		}
		timeSince, err = time.ParseDuration(str)
	} else {
		str, err = rotator.TimeSinceLastRun(rcd)
		if err != nil {
			return
		}
		timeSince, err = time.ParseDuration(str)
	}
	available = true
	return
}

func initRun(rcd rotator.RunConfigData, remoteRun ssh.SSH) (string, error) {
	remote := len(remoteRun.Host) > 0
	if remote {
		remoteRun.Cmd = fmt.Sprintf("psb-rotatorc init -name %s -compkey %s -freq %d -delay %d -year %d -month %d -day %d -initial %d -dir '%s'",
			rcd.Name,
			rcd.CompatibilityKey,
			rcd.Frequency,
			rcd.RotationDelay,
			rcd.Year.Duration,
			rcd.Month.Duration,
			rcd.Day.Duration,
			rcd.Initial.Duration,
			rcd.BackupDir)
		return remoteRun.Run()
	}
	return rotator.InitRun(rcd)
}

func runRsync(rc RunConfig, to string) (string, error) {
	rsync := rsync.Rsync{
		From:     rc.Source,
		To:       to,
		User:     rc.Username,
		Port:     rc.Port,
		Host:     rc.RemoteHost,
		Key:      rc.PrivateKeyPath,
		Flags:    "aAx",
		Includes: rc.Includes,
		Excludes: rc.Excludes,
	}
	return rsync.Run(config.configDir, rc.Name)
}

func rotate(name string, remoteRun ssh.SSH) (string, error) {
	remote := len(remoteRun.Host) > 0
	if remote {
		remoteRun.Cmd = fmt.Sprintf("psb-rotatorc rotate -name %s", name)
		return remoteRun.Run()
	}
	return rotator.Rotate(name)
}

func waitToStart(frequency time.Duration, timeSince time.Duration) {
	timeUntil := frequency - timeSince
	if timeUntil > frequency {
		return
	}
	if timeSince < frequency {
		time.Sleep(timeUntil)
	}
}

func Run(rc RunConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	runErr := func(e error) {
		logErr(fmt.Errorf("Run error: %s", e.Error()))
	}
	rcd, err := rc.GetRotatorData()
	remoteRun := ssh.SSH{
		Username: rc.Username,
		Host:     rc.RemoteHost,
		Port:     rc.Port,
		KeyPath:  rc.PrivateKeyPath,
	}
	timeSince, available, err := GetTimeSinceLastRun(rcd, remoteRun)
	if err != nil {
		runErr(err)
		return
	}
	if !available {
		return
	}
	frequency := time.Duration(rc.Rotations.Frequency) * time.Minute
	waitToStart(frequency, timeSince)
	to, err := initRun(rcd, remoteRun)
	if err != nil {
		runErr(err)
		return
	}
	_, err = runRsync(rc, to)
	if err != nil {
		runErr(err)
		return
	}
	_, err = rotate(rc.Name, remoteRun)
	if err != nil {
		runErr(err)
		return
	}
}

func Daemonize() {
	// Allow only a single instance
	instance := single.New("psb_run_daemon")
	instance.Lock()
	defer instance.Unlock()
	// Set error logger
	errorLogger = setLogOutput("run_error")
	// Start enabled operations
	runConfigs, err := GetEnabledRunConfigs()
	if err != nil {
		logErr(err)
		log.Fatal(err)
	}
	pollrate := getPollrate(runConfigs)
	var wg sync.WaitGroup
	for {
		for _, rc := range runConfigs {
			wg.Add(1)
			go Run(rc, &wg)
		}
		wg.Wait()
		time.Sleep(pollrate * time.Minute)
	}
}

func logErr(err error) {
	errorLogger.Println(err)
}
