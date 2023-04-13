package assistant

import (
	"fmt"
	clientModel "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
)

func handleCodeBlockFromChoices(choices []clientModel.Choice, filenameGenerator func(block model.FencedCodeBlock) string, handlers []model.SourceCodeHandler) {
	for _, choice := range choices {
		fmt.Println("Handling choice: ", choice)
		for _, codeBlock := range choice.Message.FencedCodeBlocks() {
			fmt.Println("Code black: ", codeBlock)
			code := codeBlock.ToSourceCode(filenameGenerator)
			fmt.Println("As source code: ", code)
			for _, handler := range handlers {
				code = handler(code)
			}
		}
	}
}
