// SPDX-License-Identifier: MIT

package model

type Library struct {
	Name        string
	DisplayName string `yaml:"display-name"`
	Icon        string
	FullPath    string `json:"-"`
	Index       string
	Data        map[string]interface{} `yaml:"-" json:"-"'`
	Commands    map[string]*Command
}

func (l *Library) getCommand(commandName string) *Command {
	command, found := l.Commands[commandName]
	if !found {
		command = &Command{
			Name:    commandName,
			Library: l,
		}
		l.Commands[commandName] = command
	}
	return command
}
