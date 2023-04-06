package assistant

import (
	"fmt"
	"github.com/spandigitial/codeassistant/model"
	"regexp"
)

var createTableRe = regexp.MustCompile(`(?sU)create_table "(\w+)".+end`)
var extractClassNameRe = regexp.MustCompile(`export class (\w+) {`)

func (a ChatGPTBasedCodeAssistant) RailsSchemaToEntities(railsSchema string, handlers ...func(code model.SourceCode) model.SourceCode) {

	for _, match := range createTableRe.FindAllStringSubmatch(railsSchema, -1) {

		message := model.Message{
			Role:    "user",
			Content: fmt.Sprintf("Convert ruby on rails create_table to a NestJS entity: \"%s\". Format all code examples in Markdown. No javascript examples.", match),
		}

		choices, err := a.chatGPT.Completion(message)
		if err == nil {
			for _, choice := range choices {
				for _, codeBlock := range choice.Message.FencedCodeBlocks() {
					code := codeBlock.ToSourceCode(
						func(block model.FencedCodeBlock) string {
							classNameMatch := extractClassNameRe.FindStringSubmatch(block.Content)
							if len(classNameMatch) > 0 {
								return fmt.Sprintf("%s.entity", classNameMatch[1])
							} else {
								println("Cannot extractive classname for match:", match[1])
								return fmt.Sprintf("%s.entity", match[1])
							}
						},
					)
					for _, handler := range handlers {
						code = handler(code)
					}
				}
			}
		}

	}

}
