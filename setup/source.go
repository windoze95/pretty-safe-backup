package setup

import (
	"log"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

type Source struct {
	Result string
}

func (s *Source) WriteAnswer(qsName string, value interface{}) error {
	s.Result = strings.Trim(value.(string), " ")
	return nil
}

func setSource(answer *string) {
	source := Source{*answer}
	prompt := &survey.Input{
		Message: "Create a name for this operation",
		Help: `ex: /home/user
If there are files, or directories you do not
wish to backup, you may define these excluded
items in the following "Excludes" step.`,
		Default: source.Result,
	}
	err := survey.AskOne(prompt, &source, nil)
	if err != nil {
		log.Fatal(err)
	}
	*answer = source.Result
}
