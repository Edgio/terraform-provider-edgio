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

// PropertiesResponse represents the entire response from the properties endpoint.
type PropertiesResponse struct {
	Type       string                `json:"@type"`
	ID         string                `json:"@id"`
	Links      properties.Links      `json:"@links"`
	TotalItems int                   `json:"total_items"`
	Items      []properties.Property `json:"items"`
}

// EdgioClient manages communication with the Edgio API.
type EdgioClient struct {
	client       *resty.Client
	token        string
	tokenExpiry  time.Time
	clientID     string
	clientSecret string
	tokenURL     string
	apiURL       string
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

func (e *EdgioClient) getToken() error {
	var tokenResp AccessTokenResponse
	resp, err := e.client.R().
		SetFormData(map[string]string{
			"client_id":     e.clientID,
			"client_secret": e.clientSecret,
			"grant_type":    "client_credentials",
			"scope":         "app.accounts",
		}).
		SetResult(&tokenResp).
		Post(e.tokenURL)

	if err != nil {
		return fmt.Errorf("failed to request token: %w", err)
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("unexpected status code for getToken: %d", resp.StatusCode())
	}

	e.token = tokenResp.AccessToken
	e.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	return nil
}

func (e *EdgioClient) GetSpecificProperty(ctx context.Context, propertyID string) (*properties.Property, error) {
	if e.token == "" || time.Now().After(e.tokenExpiry) {
		if err := e.getToken(); err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}
	}

	var property properties.Property
	resp, err := e.client.R().
		SetContext(ctx).
		SetAuthToken(e.token).
		SetResult(&property).
		Get(fmt.Sprintf("%s/properties/%s", e.apiURL, propertyID))

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code for getSpecificProperty: %d, %s", resp.StatusCode(), resp.Request.URL)
	}

	return &property, nil
}

func (e *EdgioClient) GetProperties(page int, pageSize int, organizationID string) (*PropertiesResponse, error) {
	if e.token == "" || time.Now().After(e.tokenExpiry) {
		if err := e.getToken(); err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}
	}

	var propertiesResp PropertiesResponse
	resp, err := e.client.R().
		SetAuthToken(e.token).
		SetQueryParams(map[string]string{
			"page":            fmt.Sprintf("%d", page),
			"page_size":       fmt.Sprintf("%d", pageSize),
			"organization_id": organizationID,
		}).
		SetResult(&propertiesResp).
		Get(fmt.Sprintf("%s/properties", e.apiURL))

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code for getProperties: %d", resp.StatusCode())
	}

	return &propertiesResp, nil
}

func (e *EdgioClient) CreateProperty(ctx context.Context, organizationID, slug string) (*properties.Property, error) {
	if e.token == "" || time.Now().After(e.tokenExpiry) {
		if err := e.getToken(); err != nil {
			return nil, fmt.Errorf("failed to get token: %w", err)
		}
	}

	var createdProperty properties.Property
	resp, err := e.client.R().
		SetContext(ctx).
		SetAuthToken(e.token).
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
	url := fmt.Sprintf("%s/accounts/v0.1/properties/%s", c.apiURL, propertyID)

	resp, err := c.client.R().
		SetAuthToken(c.token).
		Delete(url)

	if err != nil {
		return fmt.Errorf("error sending DELETE request: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("error deleting property: status code %d", resp.StatusCode())
	}

	return nil
}
