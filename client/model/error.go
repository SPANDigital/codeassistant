package model

import "fmt"

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Message: %v Type: %v Param: %v Code: %v", e.Message, e.Type, e.Param, e.Code)
}
