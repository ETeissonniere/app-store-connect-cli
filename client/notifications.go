package client

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotificationEndpointNotFound = errors.New("notification endpoint not found, it may not have been configured in the App Store Connect dashboard")
)

// Call the StoreKit API to request a test server to server notification.
func (c *Client) RequestTestNotification() (string, error) {
	endpoint := endpoints[storeKit].getEndpoint(c.config.UseSandbox)
	req, err := http.NewRequest("POST", "https://"+endpoint+"/inApps/v1/notifications/test", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return "TBD", nil
	case http.StatusNotFound:
		return "", ErrNotificationEndpointNotFound
	default:
		return "", handleGeneralError(resp.StatusCode, "request test notification")
	}
}
