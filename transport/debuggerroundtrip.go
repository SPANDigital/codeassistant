package transport

import (
	"bytes"
	"github.com/spandigitial/codeassistant/client/debugger"
	"net/http"
)

type DebuggerRoundtrip struct {
	transport http.RoundTripper
	debugger  *debugger.Debugger
}

func New(transport http.RoundTripper, debugger *debugger.Debugger) *DebuggerRoundtrip {
	return &DebuggerRoundtrip{
		transport: transport,
		debugger:  debugger,
	}
}

func (d *DebuggerRoundtrip) RoundTrip(request *http.Request) (*http.Response, error) {
	if d.debugger.IsRecording("request-header") {
		var bytes bytes.Buffer
		request.Header.Write(&bytes)
		d.debugger.Message("request-header", bytes.String())
	}
	response, err := d.transport.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	if d.debugger.IsRecording("response-header") {
		var bytes bytes.Buffer
		request.Header.Write(&bytes)
		d.debugger.Message("response-header", bytes.String())
	}
	return response, nil
}
