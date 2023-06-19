package llamagoremote

import "github.com/google/uuid"

type promptRequest struct {
	ID     uuid.UUID `json:"id"`
	Prompt string    `json:"prompt"`
}
