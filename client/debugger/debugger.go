package debugger

import (
	"fmt"
	"os"
)

type Debugger struct {
	detailsMap map[Detail]struct{}
}

func New(details ...Detail) *Debugger {
	detailsMap := make(map[Detail]struct{})
	for _, detail := range details {
		detailsMap[detail] = struct{}{}
	}
	return &Debugger{detailsMap: detailsMap}
}

func (d *Debugger) IsRecording(detail Detail) bool {
	_, found := d.detailsMap[detail]
	return found
}

func (d *Debugger) Message(detail Detail, message string) {
	if _, found := d.detailsMap[detail]; !found {
		return
	}
	fmt.Fprintf(os.Stderr, "%s >>> %s\n", detail, message)
}

func (d *Debugger) MessageF(detail Detail, format string, a ...any) {
	if _, found := d.detailsMap[detail]; !found {
		return
	}
	fmt.Fprintf(os.Stderr, "%s >>> %s\n", detail, fmt.Sprintf(format, a...))
}

func (d *Debugger) MessageBytes(detail Detail, message []byte) {
	if _, found := d.detailsMap[detail]; !found {
		return
	}
	fmt.Fprintf(os.Stderr, "%s >>> %s\n", detail, string(message))
}

func (d *Debugger) MessageStringer(detail Detail, message fmt.Stringer) {
	if _, found := d.detailsMap[detail]; !found {
		return
	}
	fmt.Fprintf(os.Stderr, "%s >>> %s\n", detail, message)
}
