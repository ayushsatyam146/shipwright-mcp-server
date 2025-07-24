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
- **delete_build** - Delete Build resources safely with validation

### 🚀 BuildRun Management  
- **list_buildruns** - List and filter buildruns with status
- **get_buildrun** - Get detailed buildrun information including logs and status
- **create_buildrun** - Create new BuildRuns from existing Builds or with inline specifications
- **restart_buildrun** - Restart failed or completed buildruns
- **delete_buildrun** - Delete BuildRun resources safely with validation

### 📋 Strategy Management
- **list_buildstrategies** - List namespace-scoped build strategies
- **list_clusterbuildstrategies** - List cluster-scoped build strategies

## What Makes This Special

1. **Native Kubernetes Integration** - Directly communicates with Shipwright Build APIs
2. **Rich Filtering** - Supports label selectors and prefix filtering
3. **Detailed Information** - Provides comprehensive status, timing, and failure details
4. **Easy Deployment** - Can run in-cluster or with local kubeconfig
5. **Standard MCP Protocol** - Works with any MCP-compatible client
6. **Complete CRUD Operations** - Create, read, list, and delete operations for core resources
7. **Flexible BuildRun Creation** - Support both existing Build references and inline specifications
8. **Safe Deletion** - Comprehensive validation and error handling for delete operations

## Project Structure

```
shipwright-mcp-server/
├── main.go              # Main server implementation
├── go.mod               # Go module with dependencies  
├── go.sum               # Go module checksum file
├── README.md            # Detailed usage documentation
├── OVERVIEW.md          # This file
├── CHANGELOG.md         # Version history and changes
├── Dockerfile           # Container image for deployment
├── Makefile             # Build and deployment automation
├── config.example.json  # Sample MCP client configuration
└── .gitignore          # Git ignore rules
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

1. **AI-Assisted DevOps** - AI assistants can help create builds, troubleshoot issues, check status, restart failed runs, and clean up resources
2. **Automated CI/CD** - MCP clients can programmatically create, monitor, and manage builds
3. **Developer Tools** - IDEs and editors can integrate build creation, monitoring, and cleanup
4. **Chatbots** - Slack/Teams bots can provide build operations and status
5. **GitOps Integration** - Automated build creation and lifecycle management based on repository changes

## Enhanced CRUD Capabilities

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

### Resource Deletion
- **Safe Deletion** - Validates resource existence before attempting deletion
- **Comprehensive Error Handling** - Clear error messages for not found or API errors
- **Namespace Support** - Delete resources from specific namespaces or default
- **Confirmation Messages** - Success confirmations with resource details

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

# Delete a build when no longer needed
{
  "tool": "delete_build",
  "arguments": {
    "name": "my-app",
    "namespace": "default"
  }
}
```

## Version History

- **v1.0.0** - Initial release with read-only operations
- **v1.1.0** - Added creation capabilities for Build and BuildRun resources
- **v1.2.0** - Added deletion capabilities for Build and BuildRun resources

## Next Steps

Future enhancements could include:
- Integration with build logs streaming
- Support for build triggers and webhooks  
- Metrics and observability endpoints
- Support for build artifacts and image management
- Build template creation and management
- Advanced validation and dry-run capabilities
- Bulk operations for multiple resources

This server provides a comprehensive foundation for AI-powered interactions with Shipwright Build, making container image building more accessible, automated, and developer-friendly. 