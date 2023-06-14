package openai

import "fmt"

type openAiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

func (e *openAiError) Error() string {
	return fmt.Sprintf("Message: %v Type: %v Param: %v Code: %v", e.Message, e.Type, e.Param, e.Code)
}
