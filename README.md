# [Draft - Not Live Yet] Aja (DeployAja) 🚀

> Deploy applications with managed dependencies in seconds, not hours.

Aja is a powerful CLI tool that simplifies container deployment with managed dependencies like PostgreSQL, Redis, RabbitMQ, and more. Get your app running in the cloud with auto-injected environment variables and zero configuration overhead.

[![Go Version](https://img.shields.io/badge/go-1.20+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/badge/release-beta-orange.svg)](https://github.com/deployaja/deployaja-cli/releases)

## ✨ Features

- ⚡ **Quick Deploy from Marketplace** - Instantly install and launch popular apps with `aja install`
- 🎯 **Managed Configuration** - Auto-inject connection strings for all dependencies
- 💰 **Cost Forecasting** - See deployment costs before you deploy
- 🔧 **Managed Dependencies** - PostgreSQL, Redis, MySQL, RabbitMQ, MongoDB, and more
- 🚀 **One Command Deploy** - From code to production in seconds
- 📊 **Real-time Monitoring** - Status, logs, and health checks
- 🔄 **Instant Rollbacks** - Rollback to previous version with one command

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

## Deploy APP from marketplace

```bash
# deploy n8n instantly for you
# and will create n8n.yaml in current dir
$ aja install n8n 

📊 Deployment Status

NAME          STATUS      REPLICAS   URL                                LAST DEPLOYED      
-----------   ---------   --------   --------------------------------   -------------------
n8n           deploying   1/1        https://02342.n8n.deployaja.id     2025-06-20 11:00:00

# you can edit n8n deployment spec
$ vim n8n.yaml

# redeploy the update
$ aja deploy -f n8n.yaml
```

### Deploy Your Own App

```bash
# 1. Initialize configuration
$ aja init

# 2. Edit deployaja.yaml for your app
$ vim deployaja.yaml

# 3. Login to Aja
$ aja login

# 4. See costs plan
$ aja plan

# 5. Deploy
$ aja deploy
```

## 🤖 AI-Powered Configuration Generation

The `aja gen` command uses AI to generate deployment configurations based on natural language prompts. This makes it easier to create complex configurations without manually writing YAML.

```bash
# Generate a configuration for a Node.js API with PostgreSQL
$ aja gen "create a nodejs api with postgresql database"

🤖 Generating content...
✅ Content written to /tmp/deployaja-gen-42.yaml
📝 Opening vim...
```

The generated configuration will be saved to a temporary YAML file with a unique filename and automatically opened in vim for editing. You can then:

1. Review and modify the generated configuration
2. Save the file to your project directory
3. Deploy using `aja deploy -f <filename>`

### Example Prompts

```bash
# Simple web applications
aja gen "n8n"
aja gen "vault"
aja gen "minecraft server"
aja gen "nodejs express api with cors enabled"
aja gen "python flask app with gunicorn"

# Applications with databases
aja gen "wordpress with mysql database"
aja gen "django app with postgresql and redis"
aja gen "rails app with postgres and sidekiq"

```

## 📖 Usage Examples

### Basic Web Application

```yaml
# deployaja.yaml
name: "my-web-app"
version: "1.0.0"
description: "My awesome web application"

container:
  image: "node:18-alpine"
  port: 3000

resources:
  cpu: "200m"
  memory: "256Mi"
  replicas: 2

dependencies:
  - name: "postgres"
    type: "postgresql"
    version: "15"
    config:
      database: "myapp_db"
      storage: "2Gi"

  - name: "redis"
    type: "redis"
    version: "7"
    config:
      storage: "512Mi"

env:
  - name: "NODE_ENV"
    value: "production"
  - name: "PORT"
    value: "3000"

healthCheck:
  path: "/health"
  port: 3000
  initialDelaySeconds: 30
  periodSeconds: 10
```

### Auto-Injected Environment Variables

When you deploy the above configuration, your application automatically receives:

```bash
# Database connection
POSTGRES_URL=postgresql://user:pass@postgres-myapp.aja.id:5432/myapp_db
POSTGRES_HOST=postgres-myapp.aja.id
POSTGRES_PORT=5432
POSTGRES_DB=myapp_db
POSTGRES_USER=auto_generated_user
POSTGRES_PASSWORD=auto_generated_password

# Redis connection  
REDIS_URL=redis://redis-myapp.aja.id:6379/0
REDIS_HOST=redis-myapp.aja.id
REDIS_PORT=6379

# Your custom variables
NODE_ENV=production
PORT=3000
```

### Cost Planning

```bash
$ aja plan

📋 Deployment Plan
Application: my-web-app
Image: node:18-alpine
Replicas: 2

Dependencies:
  - postgres (postgresql 15)
  - redis (redis 7)

💰 Cost Estimate
Monthly: IDR 45.500
Daily: IDR 1.500
```

## 🔧 Commands

### Core Commands

| Command | Description |
|---------|-------------|
| `aja init` | Create deployaja.yaml configuration |
| `aja gen PROMPT` | Generate aja configuration based on a prompt |
| `aja validate` | Validate configuration file |
| `aja plan` | Show deployment plan and costs |
| `aja deploy` | Deploy application |
| `aja status` | Check deployment health |
| `aja describe NAME` | Describe deployment pod details (status, containers, events, etc.) |
| `aja logs NAME` | View application logs |

### Management Commands

| Command | Description |
|---------|-------------|
| `aja env edit` | Edit environment variables in vim |
| `aja env set KEY=VALUE` | Set environment variable |
| `aja env get [KEY]` | Get environment variables |
| `aja drop NAME` | Delete deployment |

### Utility Commands

| Command | Description |
|---------|-------------|
| `aja deps [instance]` | List available dependencies and versions |
| `aja login` | Authenticate with platform |
| `aja config` | Show configuration |
| `aja search QUERY` | Search for apps in the marketplace |
| `aja install APPNAME` | Install an app from the marketplace |
| `aja install APPNAME -d DOMAIN` | Install an app with custom domain |

### Command Examples

```bash
# Deploy with dry run
aja deploy --dry-run

# Follow logs in real-time
aja logs my-app -f

# Follow logs with specific tail count
aja logs my-app --follow --tail=50

# Check deployment 
aja status

# Describe pod details with events
aja describe my-app

# List dependencies with pricing
aja deps --type postgresql

# List specific dependency instance details
aja deps instance-name

# Set environment variable
aja env set DEBUG=true

# Force delete without confirmation
aja drop my-app --force

# Search for apps in marketplace
aja search wordpress
aja search "node.js api"

# Install app from marketplace
aja install wordpress
aja install react-app

# Install app with custom domain
aja install wordpress --domain myblog.example.com
aja install react-app -d myapp.example.com

# Generate configuration using AI
aja gen "create a nodejs api with postgresql database"
aja gen "docker configuration for wordpress with mysql"
aja gen "deploy react app with redis cache"
```

### Logs Command Options

The `aja logs` command supports several options for viewing application logs:

```bash
# Basic usage - show last 100 lines
aja logs my-app

# Show specific number of lines
aja logs my-app --tail 50

# Follow logs in real-time (short form)
aja logs my-app -f

# Follow logs in real-time (long form)
aja logs my-app --follow

# Combine options - follow last 20 lines
aja logs my-app --tail 20 -f
```

**Available Flags:**
- `--tail <number>`: Number of lines to show (default: 100)
- `-f, --follow`: Follow log output in real-time

### Describe Command

The `aja describe` command provides detailed information about your deployment's pod, including its current state, containers, and events. This is useful for debugging deployment issues and understanding the current state of your application.

```bash
# Get detailed pod information
aja describe my-app
```

The describe command shows:
- **Pod Information**: Name, namespace, node, phase, IP addresses, and start time
- **Pod Conditions**: Ready, initialized, scheduled status with reasons
- **Container Details**: Image, ready state, restart count, and current state
- **Pod Events**: Recent events like pulling images, starting containers, or error conditions

Example output:
```
🔍 Fetching pod details for my-app...

📦 Pod Description

Name:        my-app-7d4f8b5c6d-xyz12
Namespace:   default
Node:        worker-node-1
Phase:       Running
Pod IP:      10.244.1.15
Host IP:     192.168.1.10
Start Time:  2024-01-15T10:30:00Z

Conditions:
  - Type: Ready, Status: True, Reason: , Message: 
  - Type: Initialized, Status: True, Reason: , Message: 
  - Type: PodScheduled, Status: True, Reason: , Message: 

Containers:
  - Name: my-app
    Image: node:18-alpine
    Ready: true
    Restarts: 0
    State: map[running:map[startedAt:2024-01-15T10:30:15Z]]

📅 Pod Events

- [Normal] Scheduled: Successfully assigned default/my-app-7d4f8b5c6d-xyz12 to worker-node-1
- [Normal] Pulling: Pulling image "node:18-alpine"
- [Normal] Pulled: Successfully pulled image "node:18-alpine"
- [Normal] Created: Created container my-app
- [Normal] Started: Started container my-app
```

### Dependencies Command

The `aja deps` command allows you to explore available dependencies and their pricing, or get detailed information about specific dependency instances.

```bash
# List all available dependencies
aja deps

# Filter dependencies by type
aja deps --type postgresql

# Get detailed information about a specific dependency instance
aja deps my-postgres-instance
```

**Available Flags:**
- `--type <type>`: Filter dependencies by type (postgresql, redis, mysql, etc.)

The command shows:
- **Dependency Information**: Name, type, available versions, and default version
- **Pricing Details**: Base cost per month, storage costs, and other pricing information
- **Specifications**: CPU and memory requirements
- **Instance Details**: When querying a specific instance, shows ID, user ID, configuration, and timestamps

## 🏪 Marketplace

The Aja marketplace provides pre-configured applications that you can deploy with a single command. Browse, search, and install applications from the community.

### Searching Apps

```bash
# Search by name
aja search wordpress

# Search by description
aja search "node.js api"

# Search by category
aja search "blog"
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

2 WordPress with WooCommerce
   WordPress with e-commerce capabilities
   Category: E-commerce
   Author: Automattic
   Version: 8.5
   Downloads: 8920
   Rating: 4.6/5.0
   Tags: cms, ecommerce, woocommerce, php

💡 Use 'aja install <app-name>' to install an app
```

### Installing Apps

```bash
# Install an app from marketplace
aja install wordpress

# Install an app with custom domain
aja install wordpress --domain mywordpress.example.com
aja install wordpress -d mywordpress.example.com
```

This will:
1. Download the app configuration from the marketplace
2. Save it as `wordpress.yaml` in your current directory
3. Configure custom domain for ingress URL (if specified)
4. Display installation instructions

Example output:
```
📦 Installing wordpress from marketplace...
✅ Configuration saved to: /path/to/wordpress-install.json
💡 Review the configuration and run 'aja deploy' to deploy
🔗 Install URL: https://marketplace.aja.id/apps/wordpress
```

The generated JSON file contains:
- App configuration in YAML format
- Installation instructions
- Metadata about the app

### Available App Categories

- **CMS**: Content Management Systems (WordPress, Drupal, etc.)
- **E-commerce**: Online stores (WooCommerce, Shopify, etc.)
- **Blog**: Blogging platforms (Ghost, Jekyll, etc.)
- **API**: Backend APIs (Node.js, Python, Go, etc.)
- **Frontend**: Single Page Applications (React, Vue, Angular, etc.)
- **Database**: Database applications (phpMyAdmin, pgAdmin, etc.)
- **Monitoring**: Monitoring tools (Grafana, Prometheus, etc.)
- **Development**: Development tools (GitLab, Jenkins, etc.)

## 🗃️ Supported Dependencies

| Service | Versions | Auto-Injected Variables |
|---------|----------|------------------------|
| **PostgreSQL** | 13, 14, 15, 16 | `POSTGRES_URL`, `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD` |
| **MySQL** | 5.7, 8.0 | `MYSQL_URL`, `MYSQL_HOST`, `MYSQL_PORT`, `MYSQL_DATABASE`, `MYSQL_USER`, `MYSQL_PASSWORD` |
| **Redis** | 6, 7 | `REDIS_URL`, `REDIS_HOST`, `REDIS_PORT` |
| **RabbitMQ** | 3.11, 3.12 | `RABBITMQ_URL`, `RABBITMQ_HOST`, `RABBITMQ_PORT`, `RABBITMQ_USER`, `RABBITMQ_PASSWORD` |
| **MongoDB** | 5.0, 6.0, 7.0 | `MONGODB_URL`, `MONGODB_HOST`, `MONGODB_PORT`, `MONGODB_DATABASE`, `MONGODB_USER`, `MONGODB_PASSWORD` |
| **Elasticsearch** | 7.17, 8.8 | `ELASTICSEARCH_URL`, `ELASTICSEARCH_HOST`, `ELASTICSEARCH_PORT` |
| **Memcached** | 1.6 | `MEMCACHED_HOST`, `MEMCACHED_PORT` |

### Dependency Configuration Examples

```yaml
dependencies:
  # PostgreSQL with custom database
  - name: "postgres"
    type: "postgresql"
    version: "15"
    config:
      database: "myapp_production"
      username: "myapp_user"
      storage: "5Gi"

  # Redis with persistence
  - name: "cache"
    type: "redis"
    version: "7"
    config:
      storage: "1Gi"
      maxMemory: "512mb"

  # RabbitMQ with management interface
  - name: "queue"
    type: "rabbitmq"
    version: "3.12"
    config:
      username: "admin"
      storage: "2Gi"

  # MongoDB cluster
  - name: "mongodb"
    type: "mongodb"
    version: "7.0"
    config:
      database: "myapp_db"
      storage: "10Gi"
```

## ⚙️ Configuration

### CLI Configuration

Location: `~/.deployaja/config.yaml`

```yaml
api:
  baseUrl: "https://aja.id/api/v1"
  timeout: 30s

output:
  format: "table"  # table, json, yaml
  colorEnabled: true

defaults:
  region: "us-east-1"
```

### Authentication

Aja uses browser-based OAuth for secure authentication:

```bash
aja login
# Opens browser for authentication
# Token stored in ~/.deployaja/token
```

## 🏗️ deployaja.yaml Reference

### Complete Configuration Example

```yaml
# Application metadata
name: "my-application"           # Required: Application name
version: "1.0.0"                # Required: Version
description: "My awesome app"    # Optional: Description

# Container configuration
container:
  image: "node:18-alpine"        # Required: Docker image
  port: 3000                     # Required: Container port

# Resource requirements  
resources:
  cpu: "200m"                    # CPU request (millicores)
  memory: "256Mi"                # Memory request
  replicas: 2                    # Number of instances

# Dependencies (managed services)
dependencies:
  - name: "postgres"
    type: "postgresql" 
    version: "15"
    config:
      database: "myapp_db"
      username: "myapp_user"
      storage: "2Gi"

# Environment variables
env:
  - name: "NODE_ENV"
    value: "production"
  - name: "LOG_LEVEL"
    value: "info"

# Health checks
healthCheck:
  path: "/health"                # Health check endpoint
  port: 3000                     # Port for health checks
  initialDelaySeconds: 30        # Delay before first check
  periodSeconds: 10              # Check interval

# Optional: Custom domain
domain: "myapp.example.com"

# Optional: Persistent storage
volumes:
  - name: "uploads"
    size: "5Gi"
    mountPath: "/app/uploads"
  - name: "logs"
    size: "1Gi" 
    mountPath: "/var/log"
```

### Validation Rules

- `name`: Required, alphanumeric with hyphens
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
```

## 🌍 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `AJA_CONFIG_DIR` | Override config directory | `~/.deployaja` |
| `AJA_DEBUG` | Enable debug logging | `false` |
| `NO_COLOR` | Disable colored output | `false` |

## 📊 Cost Optimization

### Tips for Reducing Costs

1. **Right-size Resources**
   ```yaml
   resources:
     cpu: "100m"     # Start small
     memory: "128Mi" # Scale up as needed
     replicas: 1     # Single instance for dev
   ```

2. **Optimize Dependencies**
   ```yaml
   dependencies:
     - name: "redis"
       type: "redis"
       config:
         storage: "256Mi"  # Minimal storage for cache
   ```

3. **Use Cost Planning**
   ```bash
   aja plan  # Always check costs first
   ```

## 🔐 Security

- **Secure Authentication**: Browser-based OAuth with JWT tokens
- **Encrypted Transit**: All API calls use HTTPS
- **Secret Management**: Auto-generated credentials for dependencies
- **Network Isolation**: Dependencies isolated per deployment
- **Regular Updates**: Dependencies automatically patched

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
- [Browser](https://github.com/pkg/browser) - Browser launching
- [UUID](https://github.com/google/uuid) - UUID generation

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **Website**: [aja.id](https://aja.id)
- **Documentation**: [docs.aja.id](https://docs.aja.id)
- **API Reference**: [api.aja.id](https://api.aja.id)
- **Support**: [support@aja.id](mailto:support@aja.id)
- **GitHub**: [github.com/deployaja/deployaja-cli](https://github.com/deployaja/deployaja-cli)

## ⭐ Support

If you find Aja helpful, please:
- ⭐ Star this repository
- 🐛 Report bugs via [GitHub Issues](https://github.com/deployaja/deployaja-cli/issues)
- 💡 Request features via [GitHub Discussions](https://github.com/deployaja/deployaja-cli/discussions)
- 📢 Share with your team

---

**Made with ❤️ by the Aja Team**

*Deploy applications, not infrastructure.*