package openai

import "fmt"

type LanguageModel struct {
	Id      string `json:"id"`
	OwnedBy string `json:"owned_by"`
}

func (lm LanguageModel) String() string {
	return fmt.Sprintf("Id -> %s OwnedBy -> %s\n", lm.Id, lm.OwnedBy)
}
