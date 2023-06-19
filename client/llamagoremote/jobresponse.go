package llamagoremote

import (
	"github.com/google/uuid"
	"time"
)

type jobResponse struct {
	ID      uuid.UUID `json:"id"`
	Prompt  string    `json:"prompt"`
	Output  string    `json:"output"`
	Created time.Time `json:"created"`
	Started time.Time `json:"started"`
	Model   string    `json:"model"`
	status  status    `json:"status"`
}
