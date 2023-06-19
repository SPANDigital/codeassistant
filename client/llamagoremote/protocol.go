package llamagoremote

import "fmt"

type Protocol string

const (
	HttpProtocol  Protocol = "http"
	HttpsProtocol Protocol = "https"
)

var protocolMap = map[string]Protocol{
	"http":  HttpProtocol,
	"https": HttpsProtocol,
}

func ParseProtocol(protocolStr string) (Protocol, error) {
	protocol, found := protocolMap[protocolStr]
	if found {
		return protocol, nil
	} else {
		return "", fmt.Errorf("protocol: '%s' not found", protocolStr)
	}
}
