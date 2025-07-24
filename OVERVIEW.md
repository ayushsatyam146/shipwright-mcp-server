# Shipwright Build MCP Server - Project Overview

## What We Built

This is a Model Context Protocol (MCP) server specifically designed for the Shipwright Build Kubernetes project. The server exposes Shipwright Build functionality through standardized MCP tools that can be consumed by AI assistants and other MCP clients.

## Architecture

```
┌─────────────────────┐    ┌──────────────────────┐    ┌─────────────────────┐
│   MCP Client        │◄──►│ Shipwright Build     │◄──►│ Kubernetes Cluster  │
│ (Claude, AI tools)  │    │ MCP Server           │    │ (with Shipwright)   │
└─────────────────────┘    └──────────────────────┘    └─────────────────────┘
```

## Key Features

### 📦 Build Management
- **list_builds** - List and filter builds in namespaces
- **get_build** - Get detailed build information
- **create_build** - Create new Build resources with source, strategy, and output configuration

### 🚀 BuildRun Management  
- **list_buildruns** - List and filter buildruns with status
- **get_buildrun** - Get detailed buildrun information including logs and status
- **create_buildrun** - Create new BuildRuns from existing Builds or with inline specifications
- **restart_buildrun** - Restart failed or completed buildruns

### 📋 Strategy Management
- **list_buildstrategies** - List namespace-scoped build strategies
- **list_clusterbuildstrategies** - List cluster-scoped build strategies

## What Makes This Special

1. **Native Kubernetes Integration** - Directly communicates with Shipwright Build APIs
2. **Rich Filtering** - Supports label selectors and prefix filtering
3. **Detailed Information** - Provides comprehensive status, timing, and failure details
4. **Easy Deployment** - Can run in-cluster or with local kubeconfig
5. **Standard MCP Protocol** - Works with any MCP-compatible client
6. **Full CRUD Operations** - Create, read, list, and restart operations for core resources
7. **Flexible BuildRun Creation** - Support both existing Build references and inline specifications

## Project Structure

```
server/
├── main.go              # Main server implementation
├── go.mod               # Go module with dependencies  
├── README.md            # Detailed usage documentation
├── Dockerfile           # Container image for deployment
├── Makefile             # Build and deployment automation
├── config.example.json  # Sample MCP client configuration
├── .gitignore          # Git ignore rules
└── OVERVIEW.md         # This file
```

## Supported Resources

- **Build** (`builds.shipwright.io/v1beta1`) - Container image build definitions
- **BuildRun** (`buildruns.shipwright.io/v1beta1`) - Build execution instances
- **BuildStrategy** (`buildstrategies.shipwright.io/v1beta1`) - Namespace-scoped strategies
- **ClusterBuildStrategy** (`clusterbuildstrategies.shipwright.io/v1beta1`) - Cluster-scoped strategies

## How It Integrates with Shipwright Build

This MCP server is designed as a companion tool to the main Shipwright Build project. It:

- Uses the same v1beta1 APIs as the main project
- Shares the same Go module dependencies (via replace directive)
- Follows the same conventions and patterns
- Can be deployed alongside Shipwright Build controllers

## Use Cases

1. **AI-Assisted DevOps** - AI assistants can help create builds, troubleshoot issues, check status, restart failed runs
2. **Automated CI/CD** - MCP clients can programmatically create and monitor builds
3. **Developer Tools** - IDEs and editors can integrate build creation and monitoring
4. **Chatbots** - Slack/Teams bots can provide build operations and status
5. **GitOps Integration** - Automated build creation based on repository changes

## Enhanced Creation Capabilities

### Build Creation
- Support for Git and OCI source types
- Configurable build strategies (BuildStrategy or ClusterBuildStrategy)
- Parameter passing and timeout configuration
- Flexible context directory and revision specification

### BuildRun Creation
- **Reference Mode** - Create BuildRuns from existing Build resources
- **Inline Mode** - Create BuildRuns with complete build specifications
- Service account configuration
- Parameter overrides and timeout settings
- Auto-generated names when not specified

## Quick Start

```bash
# Build the server
make build

# Run locally (requires kubeconfig)
./shipwright-build-mcp-server

# Or run in development mode
make run-dev

# Create a build via MCP client
{
  "tool": "create_build",
  "arguments": {
    "name": "my-app",
    "source-type": "Git",
    "source-url": "https://github.com/my-org/my-app",
    "strategy": "buildah",
    "output-image": "quay.io/my-org/my-app:latest"
  }
}

# Create and run a build in one step
{
  "tool": "create_buildrun",
  "arguments": {
    "source-type": "Git",
    "source-url": "https://github.com/my-org/my-app",
    "strategy": "buildah",
    "output-image": "quay.io/my-org/my-app:v1.0.0"
  }
}
```

## Version History

- **v1.0.0** - Initial release with read-only operations
- **v1.1.0** - Added creation capabilities for Build and BuildRun resources

## Next Steps

Future enhancements could include:
- Support for deleting builds and buildruns
- Integration with build logs streaming
- Support for build triggers and webhooks  
- Metrics and observability endpoints
- Support for build artifacts and image management
- Build template creation and management
- Advanced validation and dry-run capabilities

This server provides a comprehensive foundation for AI-powered interactions with Shipwright Build, making container image building more accessible, automated, and developer-friendly. 