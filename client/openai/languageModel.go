package openai

import "fmt"

type languageModel struct {
	Id      string `json:"id"`
	OwnedBy string `json:"owned_by"`
}

func (lm languageModel) String() string {
	return fmt.Sprintf("Id -> %s OwnedBy -> %s\n", lm.Id, lm.OwnedBy)
}
