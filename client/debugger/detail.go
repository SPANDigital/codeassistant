package debugger

import (
	"encoding/json"
	"fmt"
	"github.com/spandigitial/codeassistant/maps"
	"strings"
)

type Detail uint8

const (
	Undefined = iota
	RequestHeader
	SentPrompt
	RequestPayload
	ResponseHeader
	RequestTime
	FirstResponseTime
	LastResponseTime
	Configuration
	WebServer
)

var (
	detailName = map[Detail]string{
		RequestHeader:     "request-header",
		SentPrompt:        "sent-prompt",
		RequestPayload:    "request-payload",
		ResponseHeader:    "response-header",
		RequestTime:       "request-time",
		FirstResponseTime: "first-response-time",
		LastResponseTime:  "last-response-time",
		WebServer:         "web-server",
		Configuration:     "configuration",
	}
	detailValue = maps.Inverse(detailName)
)

// String allows SettingType to implement fmt.Stringer
func (s Detail) String() string {
	return detailName[s]
}

// Convert a string to a Control, returns an error if the string is unknown.
// NOTE: for JSON marshaling this must return a Control value not a pointer, which is
// common when using integer enumerations (or any primitive type alias).
func Parse(s string) (Detail, error) {
	s = strings.TrimSpace(strings.ToLower(s))
	value, ok := detailValue[s]
	if !ok {
		return Detail(0), fmt.Errorf("%q is not a valid setting type", s)
	}
	return value, nil
}

// MarshalJSON must be a value receiver to ensure that a Control on a parent object
// does not have to be a pointer in order to have it correctly marshaled.
func (s Detail) MarshalJSON() ([]byte, error) {
	// It is assumed Control implements fmt.Stringer.
	return json.Marshal(s.String())
}

// UnmarshalJSON must be a pointer receiver to ensure that the indirect from the
// parsed value can be set on the unmarshalling object. This means that the
// Parse function must return a value and not a pointer.
func (c *Detail) UnmarshalJSON(data []byte) (err error) {
	var detailStr string
	if err := json.Unmarshal(data, &detailStr); err != nil {
		return err
	}
	if *c, err = Parse(detailStr); err != nil {
		return err
	}
	return nil
}

func (c *Detail) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var detailStr string
	err := unmarshal(&detailStr)
	if err != nil {
		return err
	}

	// make sure to dereference before assignment,
	// otherwise only the local variable will be overwritten
	// and not the value the pointer actually points to
	if *c, err = Parse(detailStr); err != nil {
		return err
	}
	return nil
}
