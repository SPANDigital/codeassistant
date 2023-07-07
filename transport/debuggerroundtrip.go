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

func New(debugger *debugger.Debugger) *DebuggerRoundtrip {
	return &DebuggerRoundtrip{
		transport: http.DefaultTransport,
		debugger:  debugger,
	}
}

func (d *DebuggerRoundtrip) RoundTrip(request *http.Request) (*http.Response, error) {
	d.debugger.MessageCalculatedF(debugger.RequestHeader, "%s", func() any {
		var bytes bytes.Buffer
		request.Header.Write(&bytes)
		return bytes.String()
	})
	response, err := d.transport.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	d.debugger.MessageCalculatedF(debugger.ResponseHeader, "%s", func() any {
		var bytes bytes.Buffer
		response.Header.Write(&bytes)
		return bytes.String()
	})
	return response, nil
}
