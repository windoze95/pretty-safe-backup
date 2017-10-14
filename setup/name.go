package setup

import (
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

type Name struct {
	Result string
}

func (n *Name) WriteAnswer(name string, value interface{}) error {
	n.Result = strings.Trim(value.(string), " ")
	return nil
}

func setName(answer *string) {
	name := Name{*answer}
	prompt := &survey.Input{
		Message: "Create a name for this operation",
		Default: name.Result,
	}
	survey.AskOne(prompt, &name, nil)
	*answer = name.Result
}
