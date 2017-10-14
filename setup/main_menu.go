package setup

import (
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

func shortAnswerSlice(s []string) (r string) {
	if len(s) != 0 {
		r = s[0]
		if r != shortAnswer(r) {
			r = shortAnswer(r)
		}
		if len(s) > 1 {
			r += " ..."
		}
	}
	return
}

func mainMenu(answers *settings.Setup) *settings.Setup {
	submitted := false
	defaultOption := 0
	setDefaultOption := func(i int, o []string) {
		i++
		for len(o) <= i {
			i--
		}
		defaultOption = i
	}
	for !submitted {
		options := []string{
			"Name           " + shortAnswer(answers.Name),
			"Description    " + shortAnswer(answers.Description),
			"Source         " + shortAnswer(answers.Source),
			"Excludes       " + shortAnswerSlice(answers.Excludes),
		}
		if answers.Submittable() {
			options = append(options, ">> Next <<")
			setDefaultOption(len(options), options)
		}
		prompt := &survey.Select{
			Message: "Choose:",
			Options: options,
			Default: options[defaultOption],
		}
		selectedOption := ""
		util.ClearClient()
		survey.AskOne(prompt, &selectedOption, nil)
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
				default:
					submitted = true
				}
				setDefaultOption(i, options)
			}
		}
	}
	return answers
}
