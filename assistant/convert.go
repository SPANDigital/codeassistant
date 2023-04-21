package assistant

import (
	"bytes"
	"fmt"
	"github.com/spandigitial/codeassistant/model"
	"io"
	"path"
)

func (a ChatGPTBasedCodeAssistant) Convert(rubyCode io.Reader, rubyType string, nestJsType string, targetFileName string, codeHandlers []model.SourceCodeHandler) error {

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, rubyCode)
	rubyCodeValue := buf.String()

	messages := []model.Prompt{
		{
			Role:    "system",
			Content: "You are a typescript code generating bot. Format all output in markdown. For every question answer with one block of code which is a class in typescript. Do not return code snippets. Explain assumptions in comments.",
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Convert ruby on rails %s to NestJS %s: ```%s```.", rubyType, nestJsType, rubyCodeValue),
		},
	}

	choices, err := a.chatGPT.Completion(messages...)

	if err == nil {
		handleCodeBlockFromChoices(choices, func(block model.FencedCodeBlock) string {
			return path.Base(targetFileName)
		}, codeHandlers)
		return nil
	} else {
		return err
	}

}
