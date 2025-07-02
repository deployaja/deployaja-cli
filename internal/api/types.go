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

type DependencyInstanceResponse struct {
	ID        string      `json:"id"`
	UserID    string      `json:"userId"`
	Type      string      `json:"type"`
	Config    interface{} `json:"config"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
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

type DescribeResponse struct {
	Pod    map[string]interface{}   `json:"pod"`
	Events []map[string]interface{} `json:"events"`
}

// Marketplace types
type InstallResponse struct {
	DeploymentID   string `json:"deploymentId"`
	AppName        string `json:"appName"`
	DeploymentName string `json:"deploymentName"`
	Config         string `json:"config"` // Base64 encoded YAML config
	Status         string `json:"status"`
	Message        string `json:"message"`
	EstimatedTime  string `json:"estimatedTime,omitempty"`
	URL            string `json:"url,omitempty"`
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

type AppResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Message     string `json:"message"`
	Status      string `json:"status"`
	PublishedAt string `json:"publishedAt,omitempty"`
}

// Validate types
type ValidateResponse struct {
	Valid    bool     `json:"valid"`
	Message  string   `json:"message"`
	Warnings []string `json:"warnings,omitempty"`
}

type ValidateErrorResponse struct {
	Valid  bool              `json:"valid"`
	Error  ErrorResponse     `json:"error"`
	Errors []ValidationError `json:"errors,omitempty"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// Gen types
type GenRequest struct {
	Prompt string `json:"prompt"`
}

type GenResponse struct {
	Content string `json:"content"`
}

// Restart types
type RestartResponse struct {
	Success bool        `json:"success"`
	Data    RestartData `json:"data"`
}

type RestartData struct {
	Status        string        `json:"status"`
	Message       string        `json:"message"`
	Method        string        `json:"method"`
	RolloutStatus RolloutStatus `json:"rolloutStatus"`
}

type RolloutStatus struct {
	Generation         int `json:"generation"`
	ObservedGeneration int `json:"observedGeneration"`
	Replicas           int `json:"replicas"`
	ReadyReplicas      int `json:"readyReplicas"`
	UpdatedReplicas    int `json:"updatedReplicas"`
}

// Progress types
type DeploymentProgress struct {
	DeploymentName string `json:"deploymentName"`
	Status         string `json:"status"`
	Progress       struct {
		Percentage  int    `json:"percentage"`
		CurrentStep string `json:"currentStep"`
		IsComplete  bool   `json:"isComplete"`
	} `json:"progress"`
	Queue struct {
		Status   string  `json:"status"`
		Progress float64 `json:"progress"`
		Message  string  `json:"message"`
		Error    string  `json:"error"`
	} `json:"queue"`
	LastUpdated string `json:"lastUpdated"`
}
