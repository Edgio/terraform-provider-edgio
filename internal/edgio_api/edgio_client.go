package edgio_api

import (
	"context"
	"fmt"
	"net/http"
	"terraform-provider-edgio/internal/edgio_api/dtos/configuration"
	"terraform-provider-edgio/internal/edgio_api/dtos/environments"
	"terraform-provider-edgio/internal/edgio_api/dtos/properties"
	"terraform-provider-edgio/internal/edgio_api/dtos/purge"
	"terraform-provider-edgio/internal/edgio_api/dtos/tls"
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

// TokenCache represents a cached token. The token is stored along
// with its expiry time. Because different endpoints require different
// scopes, we store the token with the scope as the key, so that we
// can fetch the token from the cache based on the scope.
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
		tokenCache:   make(map[string]TokenCache),
	}
}

func (c *EdgioClient) getToken(scope string) (string, error) {
	if cachedToken, exists := c.tokenCache[scope]; exists && time.Now().Before(cachedToken.Expiry) {
		return cachedToken.AccessToken, nil
	}

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

	c.tokenCache[scope] = TokenCache{
		AccessToken: tokenResp.AccessToken,
		Expiry:      time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}

	return tokenResp.AccessToken, nil
}

func (c *EdgioClient) GetProperty(ctx context.Context, propertyID string) (*properties.Property, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/properties/%s", c.apiURL, propertyID)

	var property properties.Property
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(&property).
		Get(url)

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

	url := fmt.Sprintf("%s/accounts/v0.1/properties", c.apiURL)

	var propertiesResp properties.Properties
	resp, err := c.client.R().
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"page":            fmt.Sprintf("%d", page),
			"page_size":       fmt.Sprintf("%d", pageSize),
			"organization_id": organizationID,
		}).
		SetResult(&propertiesResp).
		Get(url)

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

	url := fmt.Sprintf("%s/accounts/v0.1/properties", c.apiURL)

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
		Post(url)

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
		return fmt.Errorf("error deleting property: status code %d", resp.StatusCode())
	}

	return nil
}

func (c *EdgioClient) UpdateProperty(ctx context.Context, propertyID string, slug string) (*properties.Property, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/properties/%s", c.apiURL, propertyID)

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
		return nil, fmt.Errorf("unexpected status code for updateProperty: %d", resp.StatusCode())
	}

	return &updatedProperty, nil
}

func (c *EdgioClient) GetEnvironments(page, pageSize int, propertyID string) (*environments.EnvironmentsResponse, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments", c.apiURL)

	resp, err := c.client.R().
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"page":        fmt.Sprintf("%d", page),
			"page_size":   fmt.Sprintf("%d", pageSize),
			"property_id": propertyID,
		}).
		SetResult(&environments.EnvironmentsResponse{}).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return resp.Result().(*environments.EnvironmentsResponse), nil
}

func (c *EdgioClient) GetEnvironment(environmentID string) (*environments.Environment, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments/%s", c.apiURL, environmentID)

	resp, err := c.client.R().
		SetPathParams(map[string]string{
			"environment_id": environmentID,
		}).
		SetAuthToken(token).
		SetResult(&environments.Environment{}).
		Get(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return resp.Result().(*environments.Environment), nil
}

func (c *EdgioClient) CreateEnvironment(propertyID, name string, canMembersDeploy, onlyMaintainersCanDeploy, httpRequestLogging bool) (*environments.Environment, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments", c.apiURL)

	body := map[string]interface{}{
		"property_id":                 propertyID,
		"name":                        name,
		"can_members_deploy":          canMembersDeploy,
		"only_maintainers_can_deploy": onlyMaintainersCanDeploy,
		"http_request_logging":        httpRequestLogging,
	}

	resp, err := c.client.R().
		SetBody(body).
		SetAuthToken(token).
		SetResult(&environments.Environment{}).
		Post(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return resp.Result().(*environments.Environment), nil
}

func (c *EdgioClient) UpdateEnvironment(environmentID, name string, canMembersDeploy, httpRequestLogging, preserveCache bool) (*environments.Environment, error) {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments/%s", c.apiURL, environmentID)

	body := map[string]interface{}{
		"name":                 name,
		"can_members_deploy":   canMembersDeploy,
		"http_request_logging": httpRequestLogging,
		"preserve_cache":       preserveCache,
	}

	resp, err := c.client.R().
		SetPathParams(map[string]string{
			"environment_id": environmentID,
		}).
		SetBody(body).
		SetAuthToken(token).
		SetResult(&environments.Environment{}).
		Patch(url)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return resp.Result().(*environments.Environment), nil
}

func (c *EdgioClient) DeleteEnvironment(environmentID string) error {
	token, err := c.getToken("app.accounts")
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/accounts/v0.1/environments/%s", c.apiURL, environmentID)

	resp, err := c.client.R().
		SetPathParams(map[string]string{
			"environment_id": environmentID,
		}).
		SetAuthToken(token).
		SetResult(&environments.Environment{}).
		Delete(url)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("error response: %s", resp.String())
	}

	return nil
}

func (c *EdgioClient) PurgeCache(purgeRequest *purge.PurgeRequest) (*purge.PurgeResponse, error) {
	token, err := c.getToken("app.cache.purge")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/cache/v0.1/purge-requests", c.apiURL)
	var purgeResponse purge.PurgeResponse

	resp, err := c.client.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(purgeRequest).
		SetResult(&purgeResponse).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("error response: %s", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return &purgeResponse, nil
}

func (c *EdgioClient) GetPurgeStatus(requestId string) (*purge.PurgeResponse, error) {
	token, err := c.getToken("app.cache.purge")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/cache/v0.1/purge-requests/%s", c.apiURL, requestId)
	var purgeStatus purge.PurgeResponse

	resp, err := c.client.R().
		SetAuthToken(token).
		SetResult(&purgeStatus).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error response: %s", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return &purgeStatus, nil
}

func (c *EdgioClient) GetTlsCert(tlsCertId string) (*tls.TLSCertResponse, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/tls-certs/%s", c.apiURL, tlsCertId)

	var tlsCertResponse tls.TLSCertResponse
	resp, err := c.client.R().
		SetAuthToken(token).
		SetResult(&tlsCertResponse).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("error response: %s", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error response: %s", resp.String())
	}

	return &tlsCertResponse, nil
}

func (c *EdgioClient) UploadTlsCert(req tls.UploadTlsCertRequest) (*tls.TLSCertResponse, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/tls-certs", c.apiURL)
	response := &tls.TLSCertResponse{}

	resp, err := c.client.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(response).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to upload TLS certificate: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API responded with error: %s", resp.String())
	}

	return response, nil
}

func (c *EdgioClient) GenerateTlsCert(environmentId string) (*tls.TLSCertResponse, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/tls-certs/generate", c.apiURL)
	request := map[string]interface{}{
		"environment_id": environmentId,
	}
	response := &tls.TLSCertResponse{}

	resp, err := c.client.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(response).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to upload TLS certificate: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API responded with error: %s", resp.String())
	}

	return response, nil
}

func (c *EdgioClient) GetTlsCerts(page int, pageSize int, environmentID string) (*tls.TLSCertSResponse, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/tls-certs", c.apiURL)

	var tlsCertsResponse tls.TLSCertSResponse
	resp, err := c.client.R().
		SetAuthToken(token).
		SetQueryParams(map[string]string{
			"page":           fmt.Sprintf("%d", page),
			"page_size":      fmt.Sprintf("%d", pageSize),
			"environment_id": environmentID,
		}).
		SetResult(&tlsCertsResponse).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for getTlsCerts: %d", resp.StatusCode())
	}

	return &tlsCertsResponse, nil
}

func (c *EdgioClient) UploadCdnConfiguration(config *configuration.CDNConfiguration) (*configuration.CDNConfiguration, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("%s/config/v0.1/configs", c.apiURL)
	var response configuration.CDNConfiguration

	resp, err := c.client.R().
		SetAuthToken(token).
		SetHeader("Content-Type", "application/json").
		SetBody(config).
		SetResult(&response).
		Post(url)

	if err != nil {
		return nil, fmt.Errorf("failed to upload CDN configuration: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for uploadCdnConfiguration: %d, %s", resp.StatusCode(), resp.Body())
	}

	return &response, nil
}

func (c *EdgioClient) GetCDNConfiguration(configID string) (*configuration.CDNConfiguration, error) {
	token, err := c.getToken("app.config")
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	url := fmt.Sprintf("https://edgioapis.com/config/v0.1/configs/%s", configID)
	var response configuration.CDNConfiguration

	resp, err := c.client.R().
		SetAuthToken(token).
		SetResult(&response).
		Get(url)

	if err != nil {
		return nil, fmt.Errorf("failed to get CDN configuration: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("unexpected status code for GetCDNConfiguration: %d", resp.StatusCode())
	}

	return &response, nil
}
