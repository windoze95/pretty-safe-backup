package setup

import (
	"github.com/orange-lightsaber/pretty-safe-backup/util"
	"gopkg.in/AlecAivazis/survey.v1"
	"strings"
)

func newScreen(prompt survey.Prompt, selection *string) {
	util.ClearClient()
	survey.AskOne(prompt, selection, nil)
}

func mainMenu() {
	var (
		defaultOption int
		saved         bool
	)

	shortAnswer := func(s string) (r string) {
		r = s
		if len(r) >= 20 {
			r = r[0:20] + " ..."
		}
		return
	}
	firstLine := func(s []string) (r string) {
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

	for !saved {
		options := []string{
			"Name           " + shortAnswer(answers.Name),
			"Description    " + shortAnswer(answers.Description),
			"Source         " + shortAnswer(answers.Source),
			"Excludes       " + firstLine(answers.Excludes),
		}
		if answers.savable() {
			options = append(options, ">> SAVE <<")
		}
		prompt := &survey.Select{
			Message: "Choose:",
			Options: options,
			Default: options[defaultOption],
		}
		selectedOption := ""
		newScreen(prompt, &selectedOption)
		switch selectedOption {
		case options[0]:
			defaultOption = 1
			getName()
		case options[1]:
			defaultOption = 2
			getDescription()
		case options[2]:
			defaultOption = 3
			getSource()
		case options[3]:
			defaultOption = 4
			getExcludes()
		default:
			saved = true
		}
	}
}

func getName() {
	prompt := &survey.Input{
		Message: "Create a name for this operation",
		Default: answers.Name,
	}
	newScreen(prompt, &answers.Name)

	err := answers.Name == ""
	if err {
		getName()
	}
}

func getDescription() {
	prompt := &survey.Input{
		Message: "Write a brief description (optional)",
		Default: answers.Description,
	}
	newScreen(prompt, &answers.Description)
}

func getSource() {
	prompt := &survey.Input{
		Message: "Create a name for this operation",
		Help: `ex: /home/user
If there are files, or directories you do not
wish to backup, you may define these excluded
items in the following "Excludes" step.`,
		Default: answers.Source,
	}
	newScreen(prompt, &answers.Source)

	err := answers.Source == ""
	if err {
		getSource()
	}
}

func getExcludes() {
	excludes := "\n"
	for _, oneLine := range answers.Excludes {
		excludes += oneLine + "\n"
	}

	prompt := &survey.Editor{
		Message: "Excluded directories and files (optional)",
		Default: excludes,
	}
	newScreen(prompt, &excludes)

	format := func(excludeList []string) []string {
		var formatedList []string
		for _, str := range excludeList {
			str = strings.Trim(str, " ")
			if str != "" {
				formatedList = append(formatedList, str)
			}
		}
		return formatedList
	}
	answers.Excludes = format(strings.Split(excludes, "\n"))
}
