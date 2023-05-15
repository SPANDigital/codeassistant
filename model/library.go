// SPDX-License-Identifier: MIT

package model

type Library struct {
	Name        string
	DisplayName string `yaml:"display-name"`
	Icon        string
	FullPath    string `json:"-"`
	Index       string
	Commands    map[string]*Command
}
