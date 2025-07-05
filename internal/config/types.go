package config

// DeploymentConfig represents the deployaja.yaml structure
type DeploymentConfig struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description,omitempty"`

	Container struct {
		Image string `yaml:"image"`
		Port  int    `yaml:"port"`
	} `yaml:"container"`

	Resources struct {
		CPU      string `yaml:"cpu"`
		Memory   string `yaml:"memory"`
		Replicas int    `yaml:"replicas"`
	} `yaml:"resources"`

	Dependencies []Dependency `yaml:"dependencies,omitempty"`

	Env []EnvVar `yaml:"env,omitempty"`

	HealthCheck struct {
		Path                string `yaml:"path"`
		Port                int    `yaml:"port"`
		InitialDelaySeconds int    `yaml:"initialDelaySeconds"`
		PeriodSeconds       int    `yaml:"periodSeconds"`
	} `yaml:"healthCheck,omitempty"`

	Domain  string            `yaml:"domain,omitempty"`
	Volumes []Volume          `yaml:"volumes,omitempty"`
	EnvMap  map[string]string `yaml:"envMap,omitempty"`
	DockerConfig *DockerConfig `yaml:"dockerConfig,omitempty"`
}

type DockerConfig struct {
	Auths map[string]DockerAuth `yaml:"auths"`
}

type DockerAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
	Auth     string `yaml:"auth"`
}

type Dependency struct {
	Name    string                 `yaml:"name"`
	Type    string                 `yaml:"type"`
	Version string                 `yaml:"version"`
	Config  map[string]interface{} `yaml:"config,omitempty"`
	Storage string                 `yaml:"storage,omitempty"`
}

type EnvVar struct {
	Name        string `yaml:"name"`
	Value       string `yaml:"value"`
	UserManaged bool   `yaml:"userManaged"`
}

type Volume struct {
	Name      string `yaml:"name"`
	Size      string `yaml:"size"`
	MountPath string `yaml:"mountPath"`
}
