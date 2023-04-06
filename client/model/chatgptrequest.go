package model

import "github.com/spandigitial/codeassistant/model"

type ChatGPTRequest struct {
	Model    string          `json:"model"`
	Messages []model.Message `json:"messages"`
}
