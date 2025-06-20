package api

// API Response structures
type AuthCheckResponse struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"`
}

type CostResponse struct {
	EstimatedCost struct {
		Monthly  float64 `json:"monthly"`
		Daily    float64 `json:"daily"`
		Currency string  `json:"currency"`
	} `json:"estimatedCost"`
	Breakdown struct {
		Compute      float64            `json:"compute"`
		Storage      float64            `json:"storage"`
		Network      float64            `json:"network"`
		Dependencies map[string]float64 `json:"dependencies,omitempty"`
	} `json:"breakdown"`
}

type DeployResponse struct {
	DeploymentID  string `json:"deploymentId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	EstimatedTime string `json:"estimatedTime,omitempty"`
	URL           string `json:"url,omitempty"`
}

type StatusResponse struct {
	Deployments []DeploymentStatus `json:"deployments"`
}

type DeploymentStatus struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	URL          string `json:"url,omitempty"`
	LastDeployed string `json:"lastDeployed"`
	CreatedAt    string `json:"createdAt,omitempty"`
	Replicas     struct {
		Desired   int `json:"desired"`
		Available int `json:"available"`
	} `json:"replicas"`
}

type DependenciesResponse struct {
	Dependencies []DependencyInfo `json:"dependencies"`
}

type DependencyInfo struct {
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	Versions       []string `json:"versions"`
	DefaultVersion string   `json:"defaultVersion"`
	Pricing        struct {
		Base    float64 `json:"base"`
		Storage float64 `json:"storage,omitempty"`
	} `json:"pricing"`
	Specs struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	} `json:"specs"`
}

type ErrorResponse struct {
	Error struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Source    string `json:"source"`
}

// Marketplace types
type InstallResponse struct {
	AppName    string `json:"appName"`
	Config     string `json:"config"` // Base64 encoded YAML config
	Message    string `json:"message"`
	Status     string `json:"status"`
	InstallURL string `json:"installUrl,omitempty"`
}

type SearchResponse struct {
	Apps  []MarketplaceApp `json:"apps"`
	Total int              `json:"total"`
}

type MarketplaceApp struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Author      string   `json:"author"`
	Version     string   `json:"version"`
	Downloads   int      `json:"downloads"`
	Rating      float64  `json:"rating"`
	Image       string   `json:"image,omitempty"`
	Repository  string   `json:"repository,omitempty"`
}
