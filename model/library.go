// SPDX-License-Identifier: MIT

package model

import (
	"path/filepath"
	"strings"
)

type Library struct {
	Name           string
	DisplayName    string `yaml:"display-name"`
	Icon           string
	Path           string   `json:"-"`
	BuiltFromPaths []string `json:"-"`
	Index          string
	Data           map[string]interface{} `yaml:"-" json:"-"'`
	Commands       map[string]*Command
}

func (l *Library) addBuildPath(path string) {
	l.BuiltFromPaths = append(l.BuiltFromPaths, filepath.Base(path))
}

func (l *Library) getCommand(path string) *Command {
	base := filepath.Base(path)
	frontName := strings.Split(base, ".")[0]
	command, found := l.Commands[frontName]
	if !found {
		command = &Command{
			Name:    frontName,
			Library: l,
		}
		l.Commands[frontName] = command
	}
	command.BuiltFromPaths = append(command.BuiltFromPaths, base)
	return command
}
