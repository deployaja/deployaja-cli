package api

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"deployaja-cli/internal/config"

	"gopkg.in/yaml.v3"
)

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	Subject   string `json:"sub"`
	Email     string `json:"email"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

// APIClient handles all API communications
type APIClient struct {
	BaseURL    string
	LoginURL   string
	HTTPClient *http.Client
	Token      string
	Claims     *JWTClaims
}

func NewApiClient(token string) *APIClient {
	baseUrl := os.Getenv("DEPLOYAJA_API_URL")
	if baseUrl == "" {
		baseUrl = "https://deployaja.id"
	}

	client := &APIClient{
		BaseURL:  baseUrl + "/api/v1",
		LoginURL: baseUrl + "/login",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token: token,
	}

	// Parse JWT claims if token is provided
	if token != "" {
		claims, err := client.parseJWTClaims(token)
		if err == nil {
			client.Claims = claims
		}
	}

	return client
}

// parseJWTClaims parses JWT token and extracts claims without verification
// Note: This is for client-side token inspection only, server still validates
func (c *APIClient) parseJWTClaims(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT token format")
	}

	// Decode payload (second part)
	payload := parts[1]
	// Add padding if needed
	for len(payload)%4 != 0 {
		payload += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %v", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse JWT claims: %v", err)
	}

	return &claims, nil
}

// IsTokenExpired checks if the current token is expired
func (c *APIClient) IsTokenExpired() bool {
	if c.Claims == nil {
		return true
	}

	// Check if token expires within the next 5 minutes (buffer for refresh)
	expiryTime := time.Unix(c.Claims.ExpiresAt, 0)
	return time.Now().Add(5 * time.Minute).After(expiryTime)
}

// GetTokenInfo returns information about the current token
func (c *APIClient) GetTokenInfo() map[string]interface{} {
	if c.Claims == nil {
		return map[string]interface{}{
			"valid": false,
			"error": "No token or invalid token format",
		}
	}

	expiryTime := time.Unix(c.Claims.ExpiresAt, 0)
	issuedTime := time.Unix(c.Claims.IssuedAt, 0)

	return map[string]interface{}{
		"valid":          !c.IsTokenExpired(),
		"subject":        c.Claims.Subject,
		"email":          c.Claims.Email,
		"issued_at":      issuedTime.Format(time.RFC3339),
		"expires_at":     expiryTime.Format(time.RFC3339),
		"expired":        c.IsTokenExpired(),
		"time_to_expiry": time.Until(expiryTime).String(),
	}
}

// RefreshToken attempts to refresh the current token
func (c *APIClient) RefreshToken() error {
	if c.Token == "" {
		return fmt.Errorf("no token to refresh")
	}

	resp, err := c.makeRequest("POST", c.BaseURL+"/auth/refresh", nil)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %v", err)
	}
	defer resp.Body.Close()

	var refreshResp struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		return fmt.Errorf("failed to decode refresh response: %v", err)
	}

	// Update token and claims
	c.Token = refreshResp.Token
	claims, err := c.parseJWTClaims(refreshResp.Token)
	if err == nil {
		c.Claims = claims
	}

	// Save the new token
	if err := config.SaveToken(refreshResp.Token); err != nil {
		return fmt.Errorf("failed to save refreshed token: %v", err)
	}

	return nil
}

// ensureValidToken checks token validity and refreshes if needed
func (c *APIClient) ensureValidToken() error {
	if c.Token == "" {
		return fmt.Errorf("no authentication token")
	}

	if c.IsTokenExpired() {
		if err := c.RefreshToken(); err != nil {
			return fmt.Errorf("token expired and refresh failed: %v", err)
		}
	}

	return nil
}

// API Client methods

func (c *APIClient) makeRequest(method, url string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		resp.Body.Close()

		// Handle token expiry specifically
		if resp.StatusCode == 401 && strings.Contains(strings.ToLower(errResp.Error.Message), "token") {
			return nil, fmt.Errorf("authentication failed: %s (try running 'aja login')", errResp.Error.Message)
		}

		return nil, fmt.Errorf("API error: %s", errResp.Error.Message)
	}

	return resp, nil
}

// makeAuthenticatedRequest wraps makeRequest with token validation
func (c *APIClient) makeAuthenticatedRequest(method, url string, body interface{}) (*http.Response, error) {
	if err := c.ensureValidToken(); err != nil {
		return nil, err
	}

	return c.makeRequest(method, url, body)
}

func (c *APIClient) CheckAuth(sessionCode string) (string, error) {
	url := fmt.Sprintf("%s/check?ses=%s", c.BaseURL, sessionCode)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var authResp AuthCheckResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", err
	}

	if authResp.Status == "authenticated" {
		return authResp.Token, nil
	}

	return "", fmt.Errorf("authentication pending")
}

func (c *APIClient) GetCostEstimate(config *config.DeploymentConfig) (*CostResponse, error) {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	encodedConfig := base64.StdEncoding.EncodeToString(yamlData)

	body := map[string]string{
		"deploymentConfig": encodedConfig,
	}

	resp, err := c.makeRequest("POST", c.BaseURL+"/cost", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var costResp CostResponse
	err = json.NewDecoder(resp.Body).Decode(&costResp)
	return &costResp, err
}

func (c *APIClient) Deploy(config *config.DeploymentConfig, dryRun bool) (*DeployResponse, error) {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	encodedConfig := base64.StdEncoding.EncodeToString(yamlData)

	body := map[string]interface{}{
		"deploymentConfig": encodedConfig,
		"dryRun":           dryRun,
	}

	resp, err := c.makeAuthenticatedRequest("POST", c.BaseURL+"/deploy", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deployResp DeployResponse
	err = json.NewDecoder(resp.Body).Decode(&deployResp)
	return &deployResp, err
}

func (c *APIClient) GetStatus() (*StatusResponse, error) {
	url := c.BaseURL + "/status"

	resp, err := c.makeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var statusResp StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&statusResp)
	return &statusResp, err
}

func (c *APIClient) GetLogs(name string, tail int, follow bool) ([]LogEntry, error) {
	if follow {
		return nil, fmt.Errorf("use GetLogsStream for follow mode")
	}

	url := fmt.Sprintf("%s/logs/%s?tail=%d", c.BaseURL, name, tail)

	resp, err := c.makeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var logsResp struct {
		Logs []LogEntry `json:"logs"`
	}
	err = json.NewDecoder(resp.Body).Decode(&logsResp)
	return logsResp.Logs, err
}

// GetLogsStream streams logs in real-time using Server-Sent Events
func (c *APIClient) GetLogsStream(name string, tail int, logChan chan<- LogEntry, errorChan chan<- error) {
	url := fmt.Sprintf("%s/logs/%s/stream?tail=%d", c.BaseURL, name, tail)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errorChan <- err
		return
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	// Create a client without timeout for streaming
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errorChan <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		errorChan <- fmt.Errorf("API error: %s", errResp.Error.Message)
		return
	}

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// Stream ended normally
				close(logChan)
				close(errorChan)
				return
			}
			errorChan <- err
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse SSE format: "data: {...}"
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				// Stream ended
				close(logChan)
				close(errorChan)
				return
			}

			var logEntry LogEntry
			if err := json.Unmarshal([]byte(data), &logEntry); err != nil {
				errorChan <- fmt.Errorf("failed to parse log entry: %v", err)
				continue
			}

			logChan <- logEntry
		}
	}
}

func (c *APIClient) ListDeployments() (*StatusResponse, error) {
	resp, err := c.makeAuthenticatedRequest("GET", c.BaseURL+"/list", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var listResp StatusResponse
	err = json.NewDecoder(resp.Body).Decode(&listResp)
	return &listResp, err
}

func (c *APIClient) GetDependencies(depType string) (*DependenciesResponse, error) {
	url := c.BaseURL + "/dependencies"
	if depType != "" {
		url += "?type=" + depType
	}

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var depsResp DependenciesResponse
	err = json.NewDecoder(resp.Body).Decode(&depsResp)
	return &depsResp, err
}

func (c *APIClient) GetDependencyInstance() (*struct {
	Instances []DependencyInstanceResponse `json:"dependenciesInstances"`
}, error) {
	url := c.BaseURL + "/depInstance"

	resp, err := c.makeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var depInstanceResp struct {
		Instances []DependencyInstanceResponse `json:"dependenciesInstances"`
	}
	err = json.NewDecoder(resp.Body).Decode(&depInstanceResp)
	return &depInstanceResp, err
}

func (c *APIClient) GetEnvVars(deploymentName string) (map[string]string, error) {
	url := c.BaseURL + "/env"
	if deploymentName != "" {
		url += "?deploymentName=" + deploymentName
	}

	resp, err := c.makeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var envResp struct {
		Variables map[string]string `json:"variables"`
	}
	err = json.NewDecoder(resp.Body).Decode(&envResp)
	return envResp.Variables, err
}

func (c *APIClient) UpdateEnvVars(vars map[string]string, deploymentName string) error {
	body := map[string]interface{}{
		"variables": vars,
	}

	url := c.BaseURL + "/env"
	if deploymentName != "" {
		url += "?deploymentName=" + deploymentName
	}

	resp, err := c.makeRequest("PUT", url, body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func (c *APIClient) Rollback(name string) error {
	body := map[string]string{
		"deploymentName": name,
	}

	resp, err := c.makeRequest("POST", c.BaseURL+"/rollback", body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func (c *APIClient) Drop(name string) error {
	url := fmt.Sprintf("%s/drop/%s", c.BaseURL, name)

	resp, err := c.makeRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// InstallApp retrieves the configuration for a marketplace app
func (c *APIClient) InstallApp(appName, domain, name string, dryRun bool) (*InstallResponse, error) {
	url := fmt.Sprintf("%s/install?app=%s", c.BaseURL, appName)
	if domain != "" {
		url += "&domain=" + domain
	}
	if name != "" {
		url += "&name=" + name
	}
	if dryRun {
		url += "&dryRun=true"
	}

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var installResp InstallResponse
	err = json.NewDecoder(resp.Body).Decode(&installResp)
	return &installResp, err
}

// SearchApps searches for apps in the marketplace
func (c *APIClient) SearchApps(query string) (*SearchResponse, error) {
	url := fmt.Sprintf("%s/search?q=%s", c.BaseURL, query)

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResp SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResp)
	return &searchResp, err
}

// Validate validates a deployment configuration via the API
func (c *APIClient) Validate(config *config.DeploymentConfig) (*ValidateResponse, error) {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	encodedConfig := base64.StdEncoding.EncodeToString(yamlData)

	body := map[string]string{
		"deploymentConfig": encodedConfig,
	}

	resp, err := c.makeRequest("POST", c.BaseURL+"/validate", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		var validateErrResp ValidateErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&validateErrResp); err != nil {
			return nil, fmt.Errorf("failed to parse validation error response: %v", err)
		}
		return &ValidateResponse{
			Valid:   false,
			Message: validateErrResp.Error.Error.Message,
		}, fmt.Errorf("validation failed: %s", validateErrResp.Error.Error.Message)
	}

	var validateResp ValidateResponse
	err = json.NewDecoder(resp.Body).Decode(&validateResp)
	return &validateResp, err
}

// Gen generates aja configuration based on a prompt
func (c *APIClient) Gen(prompt string) (*GenResponse, error) {
	body := GenRequest{
		Prompt: prompt,
	}

	resp, err := c.makeAuthenticatedRequest("POST", c.BaseURL+"/gen", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var genResp GenResponse
	err = json.NewDecoder(resp.Body).Decode(&genResp)
	return &genResp, err
}

func (c *APIClient) Describe(deploymentName string) (*DescribeResponse, error) {
	url := fmt.Sprintf("%s/describe?deploymentName=%s", c.BaseURL, deploymentName)

	resp, err := c.makeAuthenticatedRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("%s", errResp.Error.Message)
	}

	var describeResp DescribeResponse
	err = json.NewDecoder(resp.Body).Decode(&describeResp)
	return &describeResp, err
}

// PublishApp publishes an app to the marketplace.
// It reads the config from deployaja.yaml or from the file specified by filePath (if not empty).
// The config is base64-encoded and sent as part of the request.
func (c *APIClient) PublishApp(
	name, description, category, author, version, repository, image string,
	tags []string,
	configFilePath string,
) (*AppResponse, error) {
	// Determine which config file to use
	filePath := "deployaja.yaml"
	if configFilePath != "" {
		filePath = configFilePath
	}

	// Read config file
	configData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %v", filePath, err)
	}

	// Base64 encode config
	configBase64 := base64.StdEncoding.EncodeToString(configData)

	// Prepare request body
	reqBody := map[string]interface{}{
		"name":        name,
		"description": description,
		"category":    category,
		"author":      author,
		"version":     version,
		"repository":  repository,
		"image":       image,
		"tags":        tags,
		"config":      configBase64,
		"downloads":   0,
		"rating":      0,
		"isActive":    true,
	}

	resp, err := c.makeAuthenticatedRequest("POST", c.BaseURL+"/publish", reqBody)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("%s", errResp.Error.Message)
	}

	var appResp AppResponse
	err = json.NewDecoder(resp.Body).Decode(&appResp)
	return &appResp, err
}
