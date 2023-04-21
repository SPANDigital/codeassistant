package model

type Command struct {
	Library *Library
	Name    string            `yaml:"name"`
	Model   string            `yaml:"model"`
	Usage   string            `yaml:"usage"`
	Inherit string            `yaml:"inherit"`
	Params  map[string]string `yaml:"params"`
	Prompts []Prompt          `yaml:"prompts"`
	System  bool              `yaml:"system"`
}

func (c *Command) AllPrompts() []Prompt {
	if c.Inherit != "" {
		command, found := c.Library.Commands[c.Inherit]
		if found {
			return append(command.AllPrompts(), c.Prompts...)
		}
	}
	return c.Prompts
}

func (c *Command) Run(args []string) {

}
