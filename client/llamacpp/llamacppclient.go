package llamacpp

import (
	"fmt"
	"github.com/spandigitial/codeassistant/client"
	"github.com/spandigitial/codeassistant/model"
	"os/exec"
)

type Client struct {
	binaryPath        string
	modelPath         string
	promptContextSize int
	extraArguments    []string
}

type Option func(client *Client)

func New(binaryPath string, modelPath string, promptContextSize int, options ...Option) *Client {
	c := &Client{
		binaryPath:        binaryPath,
		modelPath:         modelPath,
		promptContextSize: promptContextSize,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithExtraArguments(arguments ...string) Option {
	return func(client *Client) {
		client.extraArguments = arguments
	}
}

func (c *Client) Models(models chan<- client.LanguageModel) error {
	close(models)
	return nil
}

func (c *Client) Completion(ci *model.CommandInstance, messageParts chan<- client.MessagePart) error {
	args := append([]string{"-m", c.modelPath, "-n", fmt.Sprintf("%d", c.promptContextSize)}, c.extraArguments...)
	out, err := exec.Command(c.binaryPath, args...).Output()
	if err != nil {
		return err
	}
	messageParts <- client.MessagePart{Delta: "", Type: "Start"}
	messageParts <- client.MessagePart{Delta: string(out), Type: "Part"}
	messageParts <- client.MessagePart{Delta: "", Type: "Done"}
}
