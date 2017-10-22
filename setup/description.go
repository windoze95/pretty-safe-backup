package setup

import (
	"log"
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
	err := survey.AskOne(prompt, &description, nil)
	if err != nil {
		log.Fatal(err)
	}
	*answer = description.Result
}
