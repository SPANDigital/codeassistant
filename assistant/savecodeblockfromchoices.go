package assistant

import (
	clientModel "github.com/spandigitial/codeassistant/client/model"
	"github.com/spandigitial/codeassistant/model"
)

func handleCodeBlockFromChoices(choices []clientModel.Choice, filenameGenerator func(block model.FencedCodeBlock) string, handlers []model.SourceCodeHandler) {
	for _, choice := range choices {
		for _, codeBlock := range choice.Message.FencedCodeBlocks() {
			code := codeBlock.ToSourceCode(filenameGenerator)
			for _, handler := range handlers {
				code = handler(code)
			}
		}
	}
}
