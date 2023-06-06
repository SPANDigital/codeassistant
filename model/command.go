// SPDX-License-Identifier: MIT

package model

import (
	"errors"
)

type Command struct {
	Library        *Library          `json:"-"`
	Name           string            `yaml:"-"`
	DisplayName    string            `yaml:"display-name"`
	BuiltFromPaths []string          `yaml:"-" json:"-"`
	Usage          string            `yaml:"usage"`
	Inherit        string            `yaml:"inherit"`
	Params         map[string]string `yaml:"params"`
	UiHints        map[string]UiHint `yaml:"ui-hints"`
	Prompts        []Prompt          `yaml:"prompts" json:"-"`
	Abstract       bool              `yaml:"abstract"`
	OpenAIConfig   OpenAIConfig      `yaml:"open-ai-config"`
	VertexAIConfig VertexAIConfig    `yaml:"vertex-ai-config"`
	Script         string            `yaml:"-"`
}

func (c *Command) AllPrompts() ([]Prompt, error) {
	prompts, _, err := c.allPrompts()
	return prompts, err
}

func (c *Command) allPrompts() ([]Prompt, []string, error) {
	if c.Inherit != "" {
		command, found := c.Library.Commands[c.Inherit]
		if found {
			prompts, commands, err := command.allPrompts()
			if err != nil {
				return nil, nil, err
			}
			for _, name := range commands {
				if name == c.Name {
					return nil, nil, errors.New("inherit loop detected")
				}
			}
			return append(prompts, c.Prompts...), append(commands, c.Name), nil
		}
	}
	return c.Prompts, []string{c.Name}, nil
}
