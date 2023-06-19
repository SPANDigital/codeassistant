package llamagoremote

import (
	"github.com/google/uuid"
	"time"
)

type promptResponse struct {
	ID      uuid.UUID `json:"id"`
	Prompt  string    `json:"prompt"`
	Created time.Time `json:"created"`
	Status  status    `json:"status"`
}
