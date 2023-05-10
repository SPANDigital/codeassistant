// SPDX-License-Identifier: MIT

package model

import "errors"

type Command struct {
	Library     *Library          `json:"-"`
	Name        string            `yaml:"name"`
	Model       string            `yaml:"model"`
	Usage       string            `yaml:"usage"`
	Inherit     string            `yaml:"inherit"`
	Params      map[string]string `yaml:"params"`
	Prompts     []Prompt          `yaml:"prompts"`
	Abstract    bool              `yaml:"abstract"`
	Temperature *float32          `yaml:"temperature"`
	TopP        *float32          `yaml:"top_p"`
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
