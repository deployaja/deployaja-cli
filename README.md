# DeployAja CLI 🚀

> Deploy applications with managed dependencies in seconds, not hours.

DeployAja is a powerful CLI tool that simplifies container deployment with managed dependencies like PostgreSQL, Redis, RabbitMQ, and more. Get your app running in the cloud with auto-injected environment variables and zero configuration overhead.

[![Go Version](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/badge/release-beta-orange.svg)](https://github.com/deployaja/deployaja-cli/releases)

## ✨ Features

- ⚡ **Quick Deploy from Marketplace** - Instantly install and launch popular apps with `aja install`
- 🤖 **AI Configuration Generation** - Generate deployment configs with natural language prompts using `aja gen`
- 🎯 **Managed Configuration** - Auto-inject connection strings for all dependencies
- 💰 **Cost Forecasting** - See deployment costs before you deploy with `aja plan`
- 🔧 **Managed Dependencies** - PostgreSQL, Redis, MySQL, RabbitMQ, MongoDB, and more
- 🚀 **One Command Deploy** - From code to production in seconds
- 📊 **Real-time Monitoring** - Status, logs, and health checks
- 🔄 **Configuration Overrides** - Override any config value using `--set` flags
- 🔍 **Pod Inspection** - Describe pod details, containers, and events

## 🚀 Quick Start

### Installation

#### Download Binary (Recommended)
```bash
# macOS/Linux
curl -sSL https://deployaja.id/setup.sh | bash

# Windows
iwr -useb https://deployaja.id/setup.bat | iex
```

#### Build from Source
```bash
git clone https://github.com/deployaja/deployaja-cli.git
cd deployaja-cli
go build -o aja main.go
```

#### Using Docker
```bash
docker pull ghcr.io/deployaja/deployaja-cli/aja
```

## Deploy APP from marketplace

```bash
# Install n8n instantly
$ aja install n8n 

📦 Installing n8n from marketplace...
✅ Configuration saved to: /path/to/n8n.yaml
💡 Review the configuration and run 'aja deploy' to deploy

# Install with custom domain
$ aja install n8n --domain my-n8n.example.com

# Install with custom deployment name
$ aja install n8n --name my-workflow-tool

# You can edit the generated configuration
$ vim n8n.yaml

# Deploy the app
$ aja deploy -f n8n.yaml
```

### Deploy Your Own App

```bash
# 1. Initialize configuration
$ aja init

# 2. Edit deployaja.yaml for your app
$ vim deployaja.yaml

# 3. Login to DeployAja
$ aja login

# 4. See costs plan
$ aja plan

# 5. Deploy
$ aja deploy
```

## GitHub Action Usage

This repository provides a GitHub Action that you can use in your workflows to deploy applications using DeployAja CLI.

### Quick Start

```yaml
name: Deploy with DeployAja
on:
  push:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Deploy Application
        uses: deployaja/deployaja-cli@v1
        with:
          command: 'deploy'
          api-token: ${{ secrets.DEPLOYAJA_API_TOKEN }}
          environment: 'production'
          project-name: 'my-app'
```

### Action Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `command` | DeployAja command to execute | Yes | `status` |
| `api-token` | DeployAja API token for authentication | No | - |
| `config-file` | Path to DeployAja configuration file | No | `./deployaja.yaml` |
| `environment` | Target environment for deployment | No | `production` |
| `project-name` | Name of the project to deploy | No | - |
| `additional-args` | Additional arguments to pass to the CLI | No | - |

### Action Outputs

| Output | Description |
|--------|-------------|
| `deployment-id` | ID of the created deployment |
| `deployment-url` | URL of the deployed application |
| `status` | Status of the deployment operation |

### Example Workflows

#### Basic Deployment
```yaml
- name: Deploy to Production
  uses: deployaja/deployaja-cli@v1
  with:
    command: 'deploy'
    api-token: ${{ secrets.DEPLOYAJA_API_TOKEN }}
    environment: 'production'
```

#### Check Deployment Status
```yaml
- name: Check Deployment Status
  uses: deployaja/deployaja-cli@v1
  with:
    command: 'status'
    project-name: 'my-app'
    environment: 'production'
```

#### View Application Logs
```yaml
- name: View Logs
  uses: deployaja/deployaja-cli@v1
  with:
    command: 'logs'
    project-name: 'my-app'
    additional-args: '--tail 100'
```

## 📖 Usage Examples

### Basic Web Application

```yaml
# deployaja.yaml
name: "bima-42-app"
description: "Simple web application with nginx and postgres"

container:
  image: "nginx:latest"
  port: 80

resources:
  cpu: "500m"
  memory: "1Gi"
  replicas: 2

dependencies:
  - name: "postgres"
    type: "database"
    version: "15"
    storage: "1Gi"

env:
  - name: "NODE_ENV"
    value: "production"
  - name: "LOG_LEVEL"
    value: "info"

healthCheck:
  path: "/api/health"
  port: 8080
  initialDelaySeconds: 60
  periodSeconds: 30

domain: "bima42.deployaja.id"

volumes:
  - name: "app-storage"
    size: "1Gi"
    mountPath: "/usr/share/nginx/html"
```

## Auto-Injected Environment Variables

When you deploy with dependencies, your application automatically receives connection strings and configuration variables for all managed services.

### Cost Planning

```bash
$ aja plan

📋 Deployment Plan
Application: my-web-app
Image: nginx:latest
Replicas: 2

Dependencies:
  - postgres (database 15)

💰 Cost Estimate
Monthly: $45.50
Daily: $1.50

Breakdown:
  Compute: $25.00
  Storage: $15.00
  Network: $5.50
```

## 🔧 Commands

### Core Commands

| Command | Description |
|---------|-------------|
| `aja init` | Create deployaja.yaml configuration with random Wayang-inspired name |
| `aja gen PROMPT` | Generate deployment configuration using AI based on natural language prompt |
| `aja validate` | Validate configuration file |
| `aja plan` | Show deployment plan and costs |
| `aja deploy` | Deploy application |
| `aja status` | Check deployment health and status |
| `aja describe NAME` | Describe deployment pod details (status, containers, events, etc.) |
| `aja logs NAME` | View application logs |

### Management Commands

| Command | Description |
|---------|-------------|
| `aja env [edit\|set\|get]` | Manage environment variables |
| `aja drop NAME` | Delete deployment |
| `aja rollback NAME` | Rollback to previous version |

### Utility Commands

| Command | Description |
|---------|-------------|
| `aja deps [instance]` | List available dependencies and versions |
| `aja login` | Authenticate with platform using browser OAuth |
| `aja config` | Show configuration |
| `aja search QUERY` | Search for apps in the marketplace |
| `aja install APPNAME` | Install an app from the marketplace |
| `aja publish` | Publish your app to the marketplace |
| `aja version` | Show CLI version |

### Command Examples

```bash
# Deploy with configuration overrides
aja deploy --set container.image=nginx:alpine --set resources.replicas=3

# Deploy with custom config file
aja deploy --file my-custom-config.yaml

# Deploy with custom name override
aja deploy --name my-production-app

# Dry run deployment
aja deploy --dry-run

# Follow logs in real-time
aja logs my-app --follow

# Follow logs with specific tail count
aja logs my-app --tail 50 -f

# Check all deployments
aja status

# Describe pod details with events
aja describe my-app

# List dependencies with pricing
aja deps --type database

# Get specific dependency instance details
aja deps my-postgres-instance

# Edit environment variables in vim
aja env edit

# Set environment variable
aja env set DEBUG=true

# Get all environment variables
aja env get

# Get specific environment variable
aja env get DEBUG

# Generate configuration with AI
aja gen "create a nodejs api with postgresql database"
aja gen "docker configuration for wordpress with mysql"

# Search for apps in marketplace
aja search wordpress
aja search "node.js api"

# Install app from marketplace
aja install wordpress
aja install react-app --domain myapp.example.com --name my-react-app

# Publish your app to marketplace
aja publish
```

### Logs Command Options

The `aja logs` command supports several options for viewing application logs:

```bash
# Basic usage - show last 100 lines
aja logs my-app

# Show specific number of lines
aja logs my-app --tail 50

# Follow logs in real-time
aja logs my-app --follow
aja logs my-app -f

# Combine options - follow last 20 lines
aja logs my-app --tail 20 -f
```

**Available Flags:**
- `--tail <number>`: Number of lines to show (default: 100)
- `-f, --follow`: Follow log output in real-time

### Environment Variables Management

The `aja env` command provides comprehensive environment variable management:

```bash
# Interactive editing in vim
aja env edit

# Set a single variable
aja env set API_KEY=your-secret-key

# Get all variables
aja env get

# Get specific variable
aja env get API_KEY
```

### Describe Command

The `aja describe` command provides detailed information about your deployment's pod, including:

- **Pod Information**: Name, namespace, node, phase, IP addresses, and start time
- **Pod Conditions**: Ready, initialized, scheduled status with reasons  
- **Container Details**: Image, ready state, restart count, ports, environment variables, and volume mounts
- **Pod Events**: Recent events like pulling images, starting containers, or error conditions

```bash
# Get detailed pod information
aja describe my-app
```

### Dependencies Command

The `aja deps` command allows you to explore available dependencies and their pricing:

```bash
# List all available dependencies
aja deps

# Filter dependencies by type
aja deps --type database

# Get detailed information about a specific dependency instance
aja deps my-postgres-instance
```

## 🏪 Marketplace

The DeployAja marketplace provides pre-configured applications that you can deploy with a single command.

### Searching Apps

```bash
# Search by name
aja search wordpress

# Search by description or tags
aja search "node.js api"
```

Example output:
```
🔍 Searching for: wordpress

✅ Found 3 apps

1 WordPress
   A popular content management system
   Category: CMS
   Author: WordPress.org
   Version: 6.4
   Downloads: 15420
   Rating: 4.8/5.0
   Tags: cms, blog, php, mysql

💡 Use 'aja install <app-name>' to install an app
```

### Installing Apps

```bash
# Install an app from marketplace
aja install wordpress

# Install with custom domain
aja install wordpress --domain mywordpress.example.com

# Install with custom deployment name  
aja install wordpress --name my-blog --domain blog.example.com

# Dry run installation
aja install wordpress --dry-run
```

This will:
1. Download the app configuration from the marketplace
2. Save it as `wordpress.yaml` in your current directory
3. Configure custom domain and name if specified
4. Display installation instructions and estimated deployment time

## 🗃️ Supported Dependencies

Dependencies are automatically configured with connection strings and environment variables:

| Service | Type | Versions | Auto-Injected Variables |
|---------|------|----------|------------------------|
| **PostgreSQL** | `database` | 13, 14, 15, 16 | Connection strings and credentials |
| **MySQL** | `database` | 5.7, 8.0 | Connection strings and credentials |
| **Redis** | `cache` | 6, 7 | Connection URLs and endpoints |
| **RabbitMQ** | `queue` | 3.11, 3.12 | Connection URLs and credentials |
| **MongoDB** | `database` | 5.0, 6.0, 7.0 | Connection strings and credentials |

### Dependency Configuration Examples

```yaml
dependencies:
  # PostgreSQL database
  - name: "postgres"  
    type: "database"
    version: "15"
    storage: "5Gi"

  # Redis cache
  - name: "cache"
    type: "cache" 
    version: "7"
    storage: "1Gi"

  # RabbitMQ message queue
  - name: "queue"
    type: "queue"
    version: "3.12"
    storage: "2Gi"
```

## ⚙️ Configuration

### CLI Configuration

Location: `~/.deployaja/config.yaml`

Authentication token is stored in: `~/.deployaja/token`

### Authentication

DeployAja uses browser-based OAuth for secure authentication:

```bash
aja login
# Opens browser for authentication
# Token stored securely in ~/.deployaja/token
```

You can also set the token via environment variable:
```bash
export DEPLOYAJA_TOKEN=your-token-here
```

## 🏗️ deployaja.yaml Reference

### Complete Configuration Example

```yaml
# Application metadata
name: "arjuna-23-app"              # Required: Application name
description: "My awesome app"       # Optional: Description

# Container configuration
container:
  image: "node:18-alpine"          # Required: Docker image
  port: 3000                       # Required: Container port

# Resource requirements  
resources:
  cpu: "500m"                      # CPU request (millicores)
  memory: "1Gi"                    # Memory request
  replicas: 2                      # Number of instances

# Dependencies (managed services)
dependencies:
  - name: "postgres"
    type: "database"
    version: "15"
    storage: "2Gi"

# Environment variables
env:
  - name: "NODE_ENV"
    value: "production"
  - name: "LOG_LEVEL"
    value: "info"

# Health checks
healthCheck:
  path: "/api/health"              # Health check endpoint
  port: 8080                       # Port for health checks
  initialDelaySeconds: 60          # Delay before first check
  periodSeconds: 30                # Check interval

# Optional: Custom domain
domain: "arjuna23.deployaja.id"

# Optional: Persistent storage
volumes:
  - name: "app-storage"
    size: "1Gi"
    mountPath: "/app/data"
# Optional    
envMap:
  db_host: WORDPRESS_DB_HOST
  db_name: WORDPRESS_DB_NAME
  db_user: WORDPRESS_DB_USER
  db_password: WORDPRESS_DB_PASSWORD
  db_pass: WORDPRESS_DB_PASSWORD

```

### Validation Rules

- `name`: Required, automatically generated with Wayang mythology names if using `aja init`
- `container.image`: Required, valid Docker image reference
- `container.port`: Required, valid port number (1-65535)
- `resources.cpu`: Valid CPU request (e.g., "100m", "0.5", "1")
- `resources.memory`: Valid memory request (e.g., "128Mi", "1Gi")
- `dependencies[].type`: Must be supported dependency type

## 🔍 Troubleshooting

### Common Issues

**Authentication Errors**
```bash
# Clear stored token and re-authenticate
rm ~/.deployaja/token
aja login
```

**Deployment Failures**
```bash
# Check deployment status
aja status

# Get detailed pod information and events
aja describe my-app

# View logs for errors
aja logs my-app --tail=100

# Follow logs in real-time for debugging
aja logs my-app -f

# Validate configuration
aja validate
```

**Configuration Issues**
```bash
# Validate deployaja.yaml
aja validate

# Check CLI configuration
aja config

# Test with dry run
aja deploy --dry-run
```

### Debug Mode

Enable verbose logging:
```bash
export AJA_DEBUG=true
aja deploy
```

### Getting Help

```bash
# General help
aja --help

# Command-specific help
aja deploy --help
aja env --help
aja logs --help
```

## 🌍 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DEPLOYAJA_TOKEN` | API token for authentication | - |
| `DEPLOYAJA_API_TOKEN` | Alternative API token variable | - |
| `AJA_DEBUG` | Enable debug logging | `false` |
| `NO_COLOR` | Disable colored output | `false` |

## 💰 Cost Optimization

### Tips for Reducing Costs

1. **Right-size Resources**
   ```yaml
   resources:
     cpu: "200m"     # Start small
     memory: "256Mi" # Scale up as needed
     replicas: 1     # Single instance for dev
   ```

2. **Optimize Dependencies**
   ```yaml
   dependencies:
     - name: "cache"
       type: "cache"
       storage: "256Mi"  # Minimal storage for cache
   ```

3. **Use Cost Planning**
   ```bash
   aja plan  # Always check costs first
   ```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

### Development Setup

```bash
# Clone repository
git clone https://github.com/deployaja/deployaja-cli.git
cd deployaja-cli

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build
go build -o aja main.go

# Run
./aja --help
```

### Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Color](https://github.com/fatih/color) - Colored terminal output
- [Browser](https://github.com/pkg/browser) - Browser launching for OAuth
- [UUID](https://github.com/google/uuid) - UUID generation
- [YAML](https://gopkg.in/yaml.v3) - YAML parsing

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Website**: [deployaja.id](https://deployaja.id)
- **Documentation**: [docs.deployaja.id](https://docs.deployaja.id)
- **API Reference**: [api.deployaja.id](https://api.deployaja.id)
- **Support**: [support@deployaja.id](mailto:support@deployaja.id)
- **GitHub**: [github.com/deployaja/deployaja-cli](https://github.com/deployaja/deployaja-cli)

## ⭐ Support

If you find DeployAja helpful, please:
- ⭐ Star this repository
- 🐛 Report bugs via [GitHub Issues](https://github.com/deployaja/deployaja-cli/issues)
- 💡 Request features via [GitHub Discussions](https://github.com/deployaja/deployaja-cli/discussions)
- 📢 Share with your team

---

**Made with ❤️ by the DeployAja Team**

*Deploy applications, not infrastructure.*
