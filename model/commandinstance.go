// SPDX-License-Identifier: MIT

package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"io"
	"os"
	"strings"
	"text/template"
)

type CommandInstance struct {
	Command     *Command
	Params      map[string]string
	Prompts     []Prompt
	EnableStdin bool
	Data        map[string]interface{}
}

type TemplateFunc func() (string, error)

func stdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func fakeStdin() (string, error) {
	return "", nil
}

func (ci *CommandInstance) combineParamsData() map[string]interface{} {
	combo := make(map[string]interface{})
	for k, v := range ci.Params {
		combo[k] = v
	}
	for k, v := range ci.Data {
		combo[k] = v
	}
	return combo
}

func (ci *CommandInstance) runParamTemplate(input string) (string, error) {

	tmpl, err := template.New("paramTemplate").Funcs(sprig.FuncMap()).Funcs(map[string]any{
		"stdin": (map[bool]TemplateFunc{true: stdin, false: fakeStdin})[ci.EnableStdin],
	}).Parse(input)
	if err != nil {
		return "", err
	}
	buff := bytes.Buffer{}
	err = tmpl.Execute(&buff, ci.combineParamsData())
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (ci *CommandInstance) runContentTemplate(content string) (string, error) {
	tmpl, err := template.New("runContentTemplate").Funcs(sprig.FuncMap()).Parse(content)
	if err != nil {
		return "", err
	}
	buff := bytes.Buffer{}
	err = tmpl.Execute(&buff, ci.combineParamsData())
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (ci *CommandInstance) buildParams(defaultParams map[string]string, args []string) {
	params := make(map[string]string)
	for k, v := range ci.Command.Params {
		value, found := defaultParams[k]
		if found {
			params[k] = value
		} else {
			value, err := ci.runParamTemplate(v)
			if err == nil {
				params[k] = value
			}
		}
	}
	for _, arg := range args {
		parts := strings.Split(arg, ":")
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		}
	}
	ci.Params = params
}

func (ci *CommandInstance) buildPrompts() ([]Prompt, error) {
	prompts := make([]Prompt, 0)
	allPrompts, err := ci.Command.AllPrompts()
	if err != nil {
		return nil, err
	}
	for _, prompt := range allPrompts {
		role, err := ci.runContentTemplate(prompt.Role)
		content, err := ci.runContentTemplate(prompt.Content)
		if err != nil {
			return nil, err
		}
		prompts = append(prompts, Prompt{
			Role:    role,
			Content: content,
		})

	}
	return prompts, nil
}

func stringArguments(call *otto.FunctionCall) []string {
	ret := make([]string, len(call.ArgumentList))
	for _, argument := range call.ArgumentList {
		ret = append(ret, argument.String())
	}
	return ret
}

func (ci *CommandInstance) runScript() error {
	paramsJson, err := json.Marshal(ci.Params)
	if err != nil {
		return err
	}
	dataJson, err := json.Marshal(ci.Data)
	if err != nil {
		return err
	}
	if ci.Command.Script != "" {
		vm := otto.New()

		vm.Set("addPrompt", func(call otto.FunctionCall) otto.Value {
			ci.Prompts = append(ci.Prompts, Prompt{
				Role:    call.Argument(0).String(),
				Content: call.Argument(1).String(),
			})
			return otto.Value{}
		})
		_, err := vm.Run(fmt.Sprintf("var data=%s;\nvar params=%s;\n%s",
			string(dataJson),
			string(paramsJson),
			ci.Command.Script))
		if err != nil {
			return err
		}
	}
	return nil
}

func CommandInstanceParams(args ...string) ([]string, error) {
	libraryName, commandName := args[0], args[1]
	library, found := BuildLibraries()[libraryName]
	if !found {
		return nil, fmt.Errorf("library: '%s' not found", libraryName)
	}
	command, found := library.Commands[commandName]
	if !found {
		return nil, fmt.Errorf("command: '%s' not found in library: '%s", commandName, libraryName)
	}
	var params []string
	for param, _ := range command.Params {
		params = append(params, param)
	}
	return params, nil
}

func NewCommandInstance(enableStdin bool, defaultParams map[string]string, args ...string) (*CommandInstance, error) {
	if len(args) < 2 {
		return nil, errors.New("at least two arguments are required to construct a command instance")
	}
	libraryName, commandName := args[0], args[1]
	library, found := BuildLibraries()[libraryName]
	if !found {
		return nil, fmt.Errorf("library: '%s' not found", libraryName)
	}
	command, found := library.Commands[commandName]
	if !found {
		return nil, fmt.Errorf("command: '%s' not found in library: '%s", commandName, libraryName)
	}
	commandInstance := &CommandInstance{
		EnableStdin: enableStdin,
		Data:        library.Data,
		Command:     command,
	}
	commandInstance.buildParams(defaultParams, args[2:])
	prompts, err := commandInstance.buildPrompts()
	if err != nil {
		return nil, err
	}
	commandInstance.Prompts = prompts
	err = commandInstance.runScript()
	if err != nil {
		return nil, err
	}

	return commandInstance, nil
}
