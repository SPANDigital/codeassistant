package debugger

import (
	"fmt"
	"os"
	"strings"
)

type Debugger struct {
	detailsMap map[string]struct{}
}

func New(details ...string) *Debugger {
	var detailsMap = make(map[string]struct{})
	for _, detail := range details {
		detailsMap[detail] = struct{}{}
	}
	return &Debugger{
		detailsMap: detailsMap,
	}
}

func (d *Debugger) IsRecording(detail string) bool {
	_, found := d.detailsMap[detail]
	return found
}

func (d *Debugger) Message(detail string, message string) {
	if _, found := d.detailsMap[detail]; found {
		if !strings.HasSuffix(message, "\n") {
			message = message + "\n"
		}
		fmt.Fprintf(os.Stderr, "%s >>> %s", detail, message)
	}
}
