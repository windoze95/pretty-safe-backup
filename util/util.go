package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"

	"gopkg.in/natefinch/lumberjack.v2"
)

func ErrorMsg(msg string, err error) error {
	err = fmt.Errorf("%s: %s", msg, err.Error())
	return err
}

func ClearClient() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func IsEmptyString(str string) bool {
	return str == ""
}

func SetLogOutput(logName string) *log.Logger {
	logger := log.New(nil, log.Prefix(), log.LstdFlags)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/psb/" + logName + ".log",
		MaxSize:    256,
		MaxBackups: 1,
		MaxAge:     365,
		Compress:   true,
	})
	return logger
}

func GetUser() (usr *user.User, err error) {
	if usr, err = user.Lookup(os.Getenv("SUDO_USER")); err != nil {
		usr, err = user.Lookup(os.Getenv("USER"))
	}
	return
}
