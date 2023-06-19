package llamagoremote

import (
	"fmt"
)

type status string

const (
	processing status = "processing"
	finished   status = "finished"
)

var statusMap = map[string]status{
	"processing": processing,
	"finished":   finished,
}

func parseStatus(statusStr string) (status, error) {
	status, found := statusMap[statusStr]
	if found {
		return status, nil
	} else {
		return "", fmt.Errorf("status: '%s' not found", statusStr)
	}
}
