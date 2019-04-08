// Package cmd provides ...
package cmd

import (
	"log"

	"gopkg.in/AlecAivazis/survey.v1"
)

//Cmd return user Select value
func Cmd() string {
	var cmdQs = []*survey.Question{
		{
			Name: "Language",
			Prompt: &survey.Select{
				Message: "Choose your Template Language",
				Options: []string{"English", "Chinese"},
			},
			Validate: survey.Required,
		},
	}
	var answer = struct {
		Language string
	}{}
	err := survey.Ask(cmdQs, &answer)
	if err != nil {
		log.Println(err)
	}
	return answer.Language
}
