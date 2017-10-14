package setup

import (
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

type Source struct {
	Result string
}

func (s *Source) WriteAnswer(source string, value interface{}) error {
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
	survey.AskOne(prompt, &source, nil)
	*answer = source.Result
}
