# Shipwright Build MCP Server - Project Overview

## What We Built

This is a Model Context Protocol (MCP) server specifically designed for the Shipwright Build Kubernetes project. The server exposes Shipwright Build functionality through standardized MCP tools that can be consumed by AI assistants and other MCP clients.

## Use Cases

1. **AI-Assisted DevOps** - AI assistants can help create builds, troubleshoot issues, check status, restart failed runs, and clean up resources
2. **Automated CI/CD** - MCP clients can programmatically create, monitor, and manage builds
3. **Developer Tools** - IDEs and editors can integrate build creation, monitoring, and cleanup
4. **Chatbots** - Slack/Teams bots can provide build operations and status
5. **GitOps Integration** - Automated build creation and lifecycle management based on repository changes

## Architecture

```
┌─────────────────────┐    ┌──────────────────────┐    ┌─────────────────────┐
│   MCP Client        │◄──►│ Shipwright Build     │◄──►│ Kubernetes Cluster  │
│ (Claude, AI tools)  │    │ MCP Server           │    │ (with Shipwright)   │
└─────────────────────┘    └──────────────────────┘    └─────────────────────┘
```

## Key Features

### Build Management
- **list_builds** - List and filter builds in namespaces
- **get_build** - Get detailed build information
- **create_build** - Create new Build resources with source, strategy, and output configuration
- **delete_build** - Delete Build resources safely with validation

### BuildRun Management  
- **list_buildruns** - List and filter buildruns with status
- **get_buildrun** - Get detailed buildrun information including logs and status
- **create_buildrun** - Create new BuildRuns from existing Builds or with inline specifications
- **restart_buildrun** - Restart failed or completed buildruns
- **delete_buildrun** - Delete BuildRun resources safely with validation

### Strategy Management
- **list_buildstrategies** - List namespace-scoped build strategies
- **list_clusterbuildstrategies** - List cluster-scoped build strategies


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


