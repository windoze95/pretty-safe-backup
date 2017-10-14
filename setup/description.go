package setup

import (
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

type Description struct {
	Result string
}

func (d *Description) WriteAnswer(description string, value interface{}) error {
	d.Result = strings.Trim(value.(string), " ")
	return nil
}

func setDescription(answer *string) {
	description := Description{*answer}
	prompt := &survey.Input{
		Message: "Write a brief description (optional)",
		Default: description.Result,
	}
	survey.AskOne(prompt, &description, nil)
	*answer = description.Result
}
