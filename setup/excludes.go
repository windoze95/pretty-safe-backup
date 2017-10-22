package setup

import (
	"log"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

type Excludes struct {
	Result []string
}

func (e *Excludes) WriteAnswer(excludes string, value interface{}) error {
	var formatedList []string
	list := strings.Split(value.(string), "\n")
	for _, str := range list {
		str = strings.Trim(str, " ")
		if str != "" {
			formatedList = append(formatedList, str)
		}
	}
	e.Result = formatedList
	return nil
}

func (e Excludes) generateDefault() string {
	list := "\n"
	for _, oneLine := range e.Result {
		list += oneLine + "\n"
	}
	return list
}

func setExcludes(answer *[]string) {
	excludes := Excludes{*answer}
	prompt := &survey.Editor{
		Message: "Excluded directories and files (optional)",
		Default: excludes.generateDefault(),
	}
	err := survey.AskOne(prompt, &excludes, nil)
	if err != nil {
		log.Fatal(err)
	}
	*answer = excludes.Result
}
