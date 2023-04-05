package assistant

import (
	"fmt"
	"github.com/stoewer/go-strcase"
	"regexp"
)

func (a ChatGPTBasedCodeAssistant) RailsSchemaToEntities(railsSchema string) []Code {
	var re = regexp.MustCompile(`(?sU)create_table "(\w+)".+end`)
	var codex []Code
	for i, match := range re.FindAllStringSubmatch(railsSchema, -1) {
		fmt.Println(match[0], "found at index", i)

		prompt := fmt.Sprintf("convert ruby on rails create_table to a NestJS entity: \"%s\"", match)

		codex = append(codex, Code{
			Filename: fmt.Sprintf("%s.entity", strcase.UpperCamelCase(match[1])),
			Language: "Typescript",
			Content:  a.chatGPT.Completion(prompt),
		})

	}

	return codex
}
