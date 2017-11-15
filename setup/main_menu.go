package setup

import (
	"log"

	"github.com/orange-lightsaber/pretty-safe-backup/settings"
	"github.com/orange-lightsaber/pretty-safe-backup/util"
	"gopkg.in/AlecAivazis/survey.v1"
)

func shortAnswer(s string) (r string) {
	r = s
	if len(r) >= 20 {
		r = r[0:20] + " ..."
	}
	return
}

func mainMenu(answers *settings.RunConfig) *settings.RunConfig {
	defaultOption := 0
	setDefaultOption := func(i int, o []string) {
		i++
		for len(o) <= i {
			i--
		}
		defaultOption = i
	}
Menu:
	for {
		options := []string{
			"Name           " + shortAnswer(answers.Name),
			"Description    " + shortAnswer(answers.Description),
			"Source",
			"Excludes",
			"Destination",
			"Rotations",
		}
		if answers.Submittable() {
			options = append(options, ">> Next <<")
			setDefaultOption(len(options), options)
		}
		prompt := &survey.Select{
			Message: "Complete the required fields to setup your back-up:",
			Options: options,
			Default: options[defaultOption],
		}
		selectedOption := ""
		util.ClearClient()
		err := survey.AskOne(prompt, &selectedOption, nil)
		if err != nil {
			log.Fatal(err)
		}
		for i, o := range options {
			if o == selectedOption {
				util.ClearClient()
				switch i {
				case 0:
					setName(&answers.Name)
				case 1:
					setDescription(&answers.Description)
				case 2:
					setSource(&answers.Source)
				case 3:
					setExcludes(&answers.Excludes)
				case 4:
					setDestination(&answers.Destination)
				case 5:
					setRotations(&answers.Rotations)
				default:
					break Menu
				}
				setDefaultOption(i, options)
			}
		}
	}
	return answers
}
