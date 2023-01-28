package client

type apiId int

const (
	storeKit apiId = iota
)

type endpointConfig struct {
	Production string
	Sandbox    string
}

func (u endpointConfig) getEndpoint(useSandbox bool) string {
	if useSandbox {
		return u.Sandbox
	}
	return u.Production
}

var endpoints = map[apiId]endpointConfig{
	storeKit: {
		Production: "api.storekit.itunes.apple.com",
		Sandbox:    "api.storekit-sandbox.itunes.apple.com",
	},
}
