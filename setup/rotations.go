package setup

import (
	"log"
	"strconv"
	"strings"

	"github.com/orange-lightsaber/pretty-safe-backup/settings"
	"github.com/orange-lightsaber/pretty-safe-backup/util"
	"gopkg.in/AlecAivazis/survey.v1"
)

type Rotations struct {
	Frequency string
	Initial   string
	Daily     string
	Monthly   string
	Yearly    string
}

func (rota *Rotations) WriteAnswer(qsName string, value interface{}) error {
	util.ClearClient()
	trim := strings.Trim(value.(string), " ")
	switch qsName {
	case "frequency":
		rota.Frequency = trim
	case "initial":
		rota.Initial = trim
	case "daily":
		rota.Daily = trim
	case "monthly":
		rota.Monthly = trim
	case "yearly":
		rota.Yearly = trim
	}
	return nil
}

func setRotations(answer *settings.Rotations) {
	rota := convertIn(*answer)
	qs := []*survey.Question{
		{
			Name: "frequency",
			Prompt: &survey.Input{
				Message: "The amount of time between snapshots, in minutes (1 - 1440).\n",
				Default: rota.Frequency,
			},
		},
		{
			Name: "initial",
			Prompt: &survey.Input{
				Message: "Amount of days to keep initial snapshot rotations (1 or greater).\n",
				Default: rota.Initial,
			},
		},
		{
			Name: "daily",
			Prompt: &survey.Input{
				Message: `Amount of days to keep daily rotations (1 - 28).
  For amount in months, start at 300; 301 is one month, 310 is ten, etc.` + "\n",
				Default: rota.Daily,
			},
		},
		{
			Name: "monthly",
			Prompt: &survey.Input{
				Message: "Amount of months to keep monthly rotations.\n",
				Default: rota.Monthly,
			},
		},
		{
			Name: "yearly",
			Prompt: &survey.Input{
				Message: "Amount of years to keep yearly rotations.\n",
				Default: rota.Yearly,
			},
		},
	}
	err := survey.Ask(qs, &rota)
	if err != nil {
		log.Fatal(err)
	}
	*answer = convertOut(rota)
}

func convertIn(answer settings.Rotations) Rotations {
	if answer.Frequency == 0 {
		answer.Frequency = 5
	}
	if answer.Initial == 0 {
		answer.Initial = 1
	}
	return Rotations{
		strconv.Itoa(answer.Frequency),
		strconv.Itoa(answer.Initial),
		strconv.Itoa(answer.Daily),
		strconv.Itoa(answer.Monthly),
		strconv.Itoa(answer.Yearly),
	}
}

func convertOut(rota Rotations) settings.Rotations {
	Frequency, err := strconv.Atoi(rota.Frequency)
	Initial, err := strconv.Atoi(rota.Initial)
	Daily, err := strconv.Atoi(rota.Daily)
	Monthly, err := strconv.Atoi(rota.Monthly)
	Yearly, err := strconv.Atoi(rota.Yearly)
	if err != nil {
		log.Fatal(err)
	}
	return settings.Rotations{
		Frequency,
		Initial,
		Daily,
		Monthly,
		Yearly,
	}
}
