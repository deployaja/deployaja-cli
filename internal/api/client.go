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

// APIClient handles all API communications
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

func NewApiClient(token string) *APIClient {
	baseUrl := os.Getenv("DEPLOYAJA_API_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:3001"
	}
	return &APIClient{
		BaseURL: baseUrl + "/api/v1",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		Token: token,
	}
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
		return nil, fmt.Errorf("API error: %s", errResp.Error.Message)
	}

	return resp, nil
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

	resp, err := c.makeRequest("POST", c.BaseURL+"/deploy", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deployResp DeployResponse
	err = json.NewDecoder(resp.Body).Decode(&deployResp)
	return &deployResp, err
}

func (c *APIClient) GetStatus(name string) (*StatusResponse, error) {
	url := c.BaseURL + "/status"
	if name != "" {
		url += "?name=" + name
	}

	resp, err := c.makeRequest("GET", url, nil)
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

	resp, err := c.makeRequest("GET", url, nil)
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
	resp, err := c.makeRequest("GET", c.BaseURL+"/list", nil)
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

func (c *APIClient) GetEnvVars() (map[string]string, error) {
	resp, err := c.makeRequest("GET", c.BaseURL+"/env", nil)
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

func (c *APIClient) UpdateEnvVars(vars map[string]string) error {
	body := map[string]interface{}{
		"variables": vars,
	}

	resp, err := c.makeRequest("PUT", c.BaseURL+"/env", body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func (c *APIClient) Rollback(name, version string) error {
	body := map[string]string{
		"name":    name,
		"version": version,
	}

	resp, err := c.makeRequest("POST", c.BaseURL+"/rollback", body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}

func (c *APIClient) Drop(name string) error {
	body := map[string]string{
		"name": name,
	}

	resp, err := c.makeRequest("POST", c.BaseURL+"/drop", body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
