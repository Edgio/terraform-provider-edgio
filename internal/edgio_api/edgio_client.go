package edgio_api

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-edgio/internal/edgio_api/dtos/properties"
	"time"

	"github.com/go-resty/resty/v2"
)

// AccessTokenResponse represents the response from the token endpoint.
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type TokenCache struct {
	AccessToken string
	Expiry      time.Time
}

// EdgioClient manages communication with the Edgio API.
type EdgioClient struct {
	client       *resty.Client
	clientID     string
	clientSecret string
	tokenURL     string
	apiURL       string
	tokenCache   map[string]TokenCache
}

// NewEdgioClient creates a new EdgioClient instance with necessary details.
func NewEdgioClient(clientID, clientSecret, tokenURL, apiURL string) *EdgioClient {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(20 * time.Second)

	return &EdgioClient{
		client:       client,
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenURL:     tokenURL,
		apiURL:       apiURL,
	}
}

func (c *EdgioClient) getToken(scope string) (string, error) {
	// Check if we have a valid cached token for the given scope
	if cachedToken, exists := c.tokenCache[scope]; exists && time.Now().Before(cachedToken.Expiry) {
		return cachedToken.AccessToken, nil
	}

	// If not cached or expired, fetch a new token
	var tokenResp AccessTokenResponse
	resp, err := c.client.R().
		SetFormData(map[string]string{
			"client_id":     c.clientID,
			"client_secret": c.clientSecret,
			"grant_type":    "client_credentials",
			"scope":         scope,
		}).
		SetResult(&tokenResp).
		Post(c.tokenURL)

	if err != nil {
		return "", fmt.Errorf("failed to request token: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("unexpected status code for getToken: %d", resp.StatusCode())
	}

	// Cache the new token
	c.tokenCache[scope] = TokenCache{
		AccessToken: tokenResp.AccessToken,
		Expiry:      time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	// Return the access token
	return tokenResp.AccessToken, nil
}

func (c *EdgioClient) GetProperty(ctx context.Context, propertyID string) (*properties.Property, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var property properties.Property
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(&property).
		Get(fmt.Sprintf("%s/properties/%s", c.apiURL, propertyID))

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code for getSpecificProperty: %d, %s", resp.StatusCode(), resp.Request.URL)
	}

	return &property, nil
}

func (c *EdgioClient) GetProperties(page int, pageSize int, organizationID string) (*properties.Properties, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var propertiesResp properties.Properties
	resp, err := c.client.R().
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"page":            fmt.Sprintf("%d", page),
			"page_size":       fmt.Sprintf("%d", pageSize),
			"organization_id": organizationID,
		}).
		SetResult(&propertiesResp).
		Get(fmt.Sprintf("%s/properties", c.apiURL))

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code for getProperties: %d", resp.StatusCode())
	}

	return &propertiesResp, nil
}

func (c *EdgioClient) CreateProperty(ctx context.Context, organizationID, slug string) (*properties.Property, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var createdProperty properties.Property
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{
			"organization_id": organizationID,
			"slug":            slug,
		}).
		SetResult(&createdProperty).
		Post("https://edgioapis.com/accounts/v0.1/properties")

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != 201 {
		return nil, fmt.Errorf("unexpected status code for createProperty: %d, response: %s", resp.StatusCode(), resp.String())
	}

	return &createdProperty, nil
}

func (c *EdgioClient) DeleteProperty(propertyID string) error {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/properties/%s", c.apiURL, propertyID)

	resp, err := c.client.R().
		SetAuthToken(token).
		Delete(url)

	if err != nil {
		return fmt.Errorf("error sending DELETE request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("error deleting property: status code %d, response body: %s", resp.StatusCode(), resp.Body())
	}

	return nil
}

func (c *EdgioClient) UpdateProperty(ctx context.Context, propertyID string, slug string) (*properties.Property, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/properties/%s", c.apiURL, propertyID)

	requestBody := map[string]interface{}{
		"slug": slug,
	}

	var updatedProperty properties.Property
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(requestBody).
		SetResult(&updatedProperty).
		Patch(url)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code for updateProperty: %d, response body: %s", resp.StatusCode(), resp.Body())
	}

	return &updatedProperty, nil
}
