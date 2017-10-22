package setup

import (
	"log"
	"os"
	"strings"

	"github.com/orange-lightsaber/pretty-safe-backup/util"
	"gopkg.in/AlecAivazis/survey.v1"
)

type Destination struct {
	Result map[string]string
}

type SSHConfig struct {
	LocalHost     string `survey:"localHost"`
	RemoteHost    string `survey:"remoteHost"`
	Username      string `survey:"username"`
	Port          string `survey:"port"`
	PrivateKeyUrl string `survey:"privateKeyUrl"`
}

func (sc *SSHConfig) WriteAnswer(destination string, value interface{}) error {
	util.ClearClient()
	if destination == "localHost" {
		sc.LocalHost = strings.Trim(value.(string), " ")
	}
	if destination == "remoteHost" {
		sc.RemoteHost = strings.Trim(value.(string), " ")
	}
	if destination == "username" {
		sc.Username = strings.Trim(value.(string), " ")
	}
	if destination == "port" {
		sc.Port = strings.Trim(value.(string), " ")
	}
	if destination == "privateKeyUrl" {
		sc.PrivateKeyUrl = strings.Trim(value.(string), " ")
	}
	return nil
}

func remoteDirectory(dest map[string]string) map[string]string {
	sshConfig := SSHConfig{}
	qs := []*survey.Question{
		{
			Name: "localHost",
			Prompt: &survey.Input{
				Message: `If you have a hostname or IP on a local network to the back-up destination,
this will be used first when available. Otherwise, leave this blank.` + "\n",
				Default: setDefaultOption(dest, "localHost"),
			},
		},
		{
			Name: "remoteHost",
			Prompt: &survey.Input{
				Message: "Remote hostname or IP to the back-up destination.\n",
				Default: setDefaultOption(dest, "remoteHost"),
			},
		},
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Username on remote host.\n",
				Default: setDefaultOption(dest, "username"),
			},
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "SSH port, leave this blank to use the default port (22).\n",
				Default: setDefaultOption(dest, "port"),
			},
		},
		{
			Name: "privateKeyUrl",
			Prompt: &survey.Input{
				Message: "SSH port, leave this blank to use the default (" + os.Getenv("HOME") + "/.ssh/id_rsa).\n",
				Default: setDefaultOption(dest, "privateKeyUrl"),
			},
		},
	}
	err := survey.Ask(qs, &sshConfig)
	if err != nil {
		log.Fatal(err)
	}
	return map[string]string{
		"localHost":     sshConfig.LocalHost,
		"remoteHost":    sshConfig.RemoteHost,
		"username":      sshConfig.Username,
		"port":          sshConfig.Port,
		"privateKeyUrl": sshConfig.PrivateKeyUrl,
	}
}

type MountConfig struct {
	MountPoint string
}

func (mc *MountConfig) WriteAnswer(mountConfig string, value interface{}) error {
	mc.MountPoint = strings.Trim(value.(string), " ")
	return nil
}

func mountPoint(dest map[string]string) map[string]string {
	mountConfig := MountConfig{}
	prompt := &survey.Input{
		Message: "Absolute path to mount point.\n",
		Default: setDefaultOption(dest, "mountPoint"),
	}
	util.ClearClient()
	err := survey.AskOne(prompt, &mountConfig, nil)
	if err != nil {
		log.Fatal(err)
	}
	return map[string]string{
		"mountPoint": mountConfig.MountPoint,
	}
}

func setDestination(answer map[string]string) map[string]string {
	if answer == nil {
		answer = make(map[string]string)
	}
	destination := Destination{answer}
	options := []string{
		"Remote directory",
		"Mount point",
	}
	prompt := &survey.Select{
		Message: "Where is your target destination?",
		Options: options,
	}
	selectedOption := ""
	err := survey.AskOne(prompt, &selectedOption, nil)
	if err != nil {
		log.Fatal(err)
	}
	util.ClearClient()
	switch selectedOption {
	case options[0]:
		destination.Result = remoteDirectory(destination.Result)
	case options[1]:
		destination.Result = mountPoint(destination.Result)
	}
	destination.Result["destType"] = selectedOption
	return destination.Result
}

func setDefaultOption(dest map[string]string, d string) (r string) {
	if dest[d] != "" {
		r = dest[d]
	}
	return
}
