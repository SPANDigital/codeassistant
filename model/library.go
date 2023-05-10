// SPDX-License-Identifier: MIT

package model

type Library struct {
	Name     string
	FullPath string `json:"-"`
	Index    string
	Commands map[string]*Command
}
