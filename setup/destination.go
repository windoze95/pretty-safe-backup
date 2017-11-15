package setup

import (
	"log"
	"strings"

	"github.com/orange-lightsaber/pretty-safe-backup/settings"
	"github.com/orange-lightsaber/pretty-safe-backup/util"
	"gopkg.in/AlecAivazis/survey.v1"
)

type Destination struct {
	settings.Destination
}

func (dest *Destination) WriteAnswer(qsName string, value interface{}) error {
	util.ClearClient()
	trim := strings.Trim(value.(string), " ")
	switch qsName {
	case "path":
		dest.Path = trim
	case "localHost":
		dest.LocalHost = trim
	case "remoteHost":
		dest.RemoteHost = trim
	case "username":
		dest.Username = trim
	case "port":
		dest.Port = trim
	case "privateKeyUrl":
		dest.PrivateKeyUrl = trim
	}
	return nil
}

func remoteDirectory(dest *Destination) {
	qs := []*survey.Question{
		{
			Name: "path",
			Prompt: &survey.Input{
				Message: "Path to backup directory.\n",
				Default: dest.Path,
			},
		},
		{
			Name: "localHost",
			Prompt: &survey.Input{
				Message: `If you have a hostname or IP on a local network to the back-up destination,
  this will be used first when available. Otherwise, leave this blank.` + "\n",
				Default: dest.LocalHost,
			},
		},
		{
			Name: "remoteHost",
			Prompt: &survey.Input{
				Message: "Remote hostname or IP to the back-up destination.\n",
				Default: dest.RemoteHost,
			},
		},
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Username on remote host.\n",
				Default: dest.Username,
			},
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "SSH port, leave this blank to use the default.\n",
				Default: setDefaultOption(*dest, "port"),
			},
		},
		{
			Name: "privateKeyUrl",
			Prompt: &survey.Input{
				Message: "Location of SSH private key, leave this blank to use the default.\n",
				Default: setDefaultOption(*dest, "privateKeyUrl"),
			},
		},
	}
	err := survey.Ask(qs, dest)
	if err != nil {
		log.Fatal(err)
	}
	dest.Type = "remote"
}

func mountPoint(dest *Destination) {
	qs := []*survey.Question{
		{
			Name: "path",
			Prompt: &survey.Input{
				Message: "Absolute path to mount point.\n",
				Default: dest.Path,
			},
		},
	}
	err := survey.Ask(qs, dest)
	if err != nil {
		log.Fatal(err)
	}
	dest.Type = "mount"
	dest.LocalHost = ""
	dest.RemoteHost = ""
	dest.Username = ""
	dest.Port = ""
	dest.PrivateKeyUrl = ""
}

func setDestination(answer *settings.Destination) {
	dest := Destination{*answer}
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
		remoteDirectory(&dest)
	case options[1]:
		mountPoint(&dest)
	}
	*answer = dest.Destination
}

func setDefaultOption(dest Destination, d string) (r string) {
	switch d {
	case "port":
		r = dest.Port
		if util.IsEmptyString(dest.Port) {
			r = "22"
		}
	case "privateKeyUrl":
		r = dest.PrivateKeyUrl
		if util.IsEmptyString(dest.PrivateKeyUrl) {
			u, err := util.GetUser()
			if err != nil {
				log.Println("error showing users home directory:" + err.Error())
			}
			r = u.HomeDir + "/.ssh/id_rsa"
		}
	}
	return
}
