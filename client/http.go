package client

import (
	"crypto/ecdsa"
	"net/http"

	"github.com/eteissonniere/app-store-connect-cli/helpers"
)

// A wrapper on `*http.Client` which adds the required headers for the App Store Connect API
// and offer QoL methods for making requests.
type Client struct {
	*http.Client
}

// Metadata necessary to authenticate with the App Store Connect API.
type ClientConfig struct {
	IssuerId   string
	BundleId   string
	KeyId      string
	PrivateKey *ecdsa.PrivateKey
}

type transport struct {
	config ClientConfig
	rt     http.RoundTripper
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	jwt, err := helpers.GenerateJWT(t.config.IssuerId, t.config.BundleId, t.config.KeyId, t.config.PrivateKey)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Auth", "Authorization: Bearer "+jwt)
	return t.rt.RoundTrip(r)
}

func New(config ClientConfig) *Client {
	return &Client{
		Client: &http.Client{
			Transport: &transport{
				config: config,
				rt:     http.DefaultTransport,
			},
		},
	}
}
