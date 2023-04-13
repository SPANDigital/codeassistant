package assistant

import (
	"fmt"
	"github.com/spandigitial/codeassistant/model"
	"regexp"
	"strings"
)

var createTableRe = regexp.MustCompile(`(?sU)create_table "(\w+)".+end`)
var extractClassNameRe = regexp.MustCompile(`export class (\w+) {`)

func (a ChatGPTBasedCodeAssistant) RailsSchemaToEntities(railsSchema string, entityHandlers []model.SourceCodeHandler, serviceHandlers []model.SourceCodeHandler) error {

	for _, match := range createTableRe.FindAllStringSubmatch(railsSchema, -1) {

		messages := []model.Message{
			{
				Role:    "system",
				Content: "You are a typescript code generating bot. Format all output in markdown. For every question answer with one block of code which is a class in typescript. Do not return code snippets.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("Convert ruby on rails create_table to a NestJS entity: \"%s\".", match),
			},
		}

		var className string
		choices, err := a.chatGPT.Completion(messages...)

		if err != nil {
			println("Error:", err)
		}

		if err == nil {

			handleCodeBlockFromChoices(choices, func(block model.FencedCodeBlock) string {
				classNameMatch := extractClassNameRe.FindStringSubmatch(block.Content)
				if len(classNameMatch) > 0 {
					className = classNameMatch[1]
				} else {
					println("Cannot extractive classname for match:", match[1])
					className = match[1]
				}
				return fmt.Sprintf("%s.entity.ts", match[1])
			}, entityHandlers)

			serviceName := fmt.Sprintf("%sService", className)

			classNameLower := strings.ToLower(className)

			messages = []model.Message{
				{
					Role:    "system",
					Content: "You are a typescript code generating bot. Format all output in markdown. For every question answer with one block of code which is a class in typescript. Do not return code snippets.",
				},
				{
					Role:    "user",
					Content: fmt.Sprintf("Assuming %s entity already exists. Create a NestJs service `%s` that implements the `findAll()`, `findOne(id: number)`, `create(%s: %s)`, `update(id: number, %s: %s)`, and `delete(id: number), the service should make use of a Repository`.", className, serviceName, classNameLower, className, classNameLower, className),
				},
			}

			choices, err = a.chatGPT.Completion(messages...)
			if err != nil {
				println("Error:", err)
			}
			if err == nil {
				handleCodeBlockFromChoices(choices, func(block model.FencedCodeBlock) string {
					return match[1] + ".service.ts"
				}, serviceHandlers)
			}

		}
	}
	return nil
}
