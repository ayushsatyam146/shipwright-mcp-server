# Shipwright Build MCP Server

A Model Context Protocol (MCP) server for the Shipwright Build project. This server provides tools for interacting with Shipwright Build resources including Builds, BuildRuns, BuildStrategies, and ClusterBuildStrategies.

## Overview

This MCP server exposes Shipwright Build functionality through standardized tools that can be consumed by MCP clients. It's designed to work with the Shipwright Build Kubernetes-native CI/CD framework for building container images from source code.

## Features

### Build Management
- **list_builds** - List builds in a namespace with filtering options
- **get_build** - Get detailed information about a specific build
- **create_build** - Create a new Build resource from source
- **delete_build** - Delete a Build resource

### BuildRun Management  
- **list_buildruns** - List buildruns in a namespace with filtering options
- **get_buildrun** - Get detailed information about a specific buildrun
- **create_buildrun** - Create a new BuildRun (from existing Build or inline spec)
- **restart_buildrun** - Restart a buildrun by creating a new one
- **delete_buildrun** - Delete a BuildRun resource

### Strategy Management
- **list_buildstrategies** - List namespace-scoped build strategies with filtering options
- **list_clusterbuildstrategies** - List cluster-scoped build strategies with filtering options

## Prerequisites

- Go 1.23 or later
- Access to a Kubernetes cluster with Shipwright Build installed
- Proper Kubernetes configuration (kubeconfig or in-cluster config)

## Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd shipwright-mcp-server
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
go build -o shipwright-build-mcp-server main.go
```

## Usage

The server runs as a standard MCP server over stdin/stdout. You can use it with any MCP client.

### Running the Server

```bash
./shipwright-build-mcp-server
```

### Example Client Configuration

For Claude Desktop or other MCP clients, add this configuration:

```json
{
  "mcpServers": {
    "shipwright-build": {
      "command": "/path/to/shipwright-build-mcp-server"
    }
  }
}
```

## Available Tools

### Build Tools

#### `list_builds` – List Builds in a Namespace with Filtering Options

* `namespace`: Namespace to list builds from (string, required)
* `prefix`: Name prefix to filter builds (string, optional)  
* `label-selector`: Label selector to filter builds (string, optional)

#### `get_build` – Get a Specific Build by Name

* `name`: Name of the build to get (string, required)
* `namespace`: Namespace where the build is located (string, optional, default: "default")

#### `create_build` – Create a New Build Resource

* `name`: Name of the build to create (string, required)
* `namespace`: Namespace where the build will be created (string, optional, default: "default")
* `source-type`: Source type - "Git" or "OCI" (string, required)
* `source-url`: Source URL (Git repository or OCI image) (string, required)
* `context-dir`: Context directory within the source (string, optional)
* `revision`: Git revision (branch, tag, or commit SHA) (string, optional)
* `strategy`: Build strategy name (string, required)
* `strategy-kind`: Build strategy kind - "BuildStrategy" or "ClusterBuildStrategy" (string, optional, default: "ClusterBuildStrategy")
* `output-image`: Output container image reference (string, required)
* `parameters`: Build parameters as key-value pairs (object, optional)
* `timeout`: Build timeout duration, e.g. "30m", "1h" (string, optional)

#### `delete_build` – Delete a Build Resource

* `name`: Name of the build to delete (string, required)
* `namespace`: Namespace where the build is located (string, optional, default: "default")

### BuildRun Tools

#### `list_buildruns` – List BuildRuns in a Namespace with Filtering Options

* `namespace`: Namespace to list buildruns from (string, required)
* `prefix`: Name prefix to filter buildruns (string, optional)
* `label-selector`: Label selector to filter buildruns (string, optional)

#### `get_buildrun` – Get a Specific BuildRun by Name

* `name`: Name of the buildrun to get (string, required)
* `namespace`: Namespace where the buildrun is located (string, optional, default: "default")

#### `create_buildrun` – Create a New BuildRun Resource

This tool supports two modes:

**Mode 1: Reference Existing Build**
* `name`: Name of the buildrun (string, optional - auto-generated if not provided)
* `namespace`: Namespace where the buildrun will be created (string, optional, default: "default")
* `build-name`: Name of existing Build to run (string, required for this mode)
* `parameters`: Build parameters to override (object, optional)
* `timeout`: BuildRun timeout duration (string, optional)
* `service-account`: Service account for the buildrun (string, optional)

**Mode 2: Inline Build Specification**
* `name`: Name of the buildrun (string, optional - auto-generated if not provided)
* `namespace`: Namespace where the buildrun will be created (string, optional, default: "default")
* `source-type`: Source type - "Git" or "OCI" (string, required for this mode)
* `source-url`: Source URL (string, required for this mode)
* `context-dir`: Context directory (string, optional)
* `revision`: Git revision (string, optional)
* `strategy`: Build strategy name (string, required for this mode)
* `strategy-kind`: Build strategy kind (string, optional, default: "ClusterBuildStrategy")
* `output-image`: Output image (string, required for this mode)
* `parameters`: Build parameters (object, optional)
* `timeout`: BuildRun timeout duration (string, optional)
* `service-account`: Service account for the buildrun (string, optional)

#### `restart_buildrun` – Restart a BuildRun by Creating a New One

* `name`: Name or reference of the buildrun to restart (string, required)
* `namespace`: Namespace where the buildrun is located (string, optional, default: "default")

#### `delete_buildrun` – Delete a BuildRun Resource

* `name`: Name of the buildrun to delete (string, required)
* `namespace`: Namespace where the buildrun is located (string, optional, default: "default")

### Strategy Tools

#### `list_buildstrategies` – List BuildStrategies in a Namespace with Filtering Options

* `namespace`: Namespace to list build strategies from (string, required)
* `prefix`: Name prefix to filter build strategies (string, optional)
* `label-selector`: Label selector to filter build strategies (string, optional)

#### `list_clusterbuildstrategies` – List ClusterBuildStrategies with Filtering Options

* `prefix`: Name prefix to filter cluster build strategies (string, optional)
* `label-selector`: Label selector to filter cluster build strategies (string, optional)

## Examples

### Creating a Build

```json
{
  "name": "my-app-build",
  "namespace": "default",
  "source-type": "Git",
  "source-url": "https://github.com/my-org/my-app",
  "context-dir": ".",
  "revision": "main",
  "strategy": "buildah",
  "strategy-kind": "ClusterBuildStrategy",
  "output-image": "quay.io/my-org/my-app:latest",
  "parameters": {
    "dockerfile": "Dockerfile"
  },
  "timeout": "30m"
}
```

### Creating a BuildRun from Existing Build

```json
{
  "name": "my-app-build-run-1",
  "namespace": "default",
  "build-name": "my-app-build",
  "service-account": "build-sa"
}
```

### Creating a BuildRun with Inline Spec

```json
{
  "namespace": "default",
  "source-type": "Git",
  "source-url": "https://github.com/my-org/my-app",
  "strategy": "buildah",
  "output-image": "quay.io/my-org/my-app:v1.0.0",
  "parameters": {
    "dockerfile": "Dockerfile.prod"
  }
}
```

### Deleting a Build

```json
{
  "name": "my-app-build",
  "namespace": "default"
}
```

### Deleting a BuildRun

```json
{
  "name": "my-app-build-run-1",
  "namespace": "default"
}
```

## Configuration

The server automatically detects and uses Kubernetes configuration:

1. **In-cluster config** - When running inside a Kubernetes pod
2. **Kubeconfig file** - When running outside the cluster (uses `~/.kube/config` by default)

## Docker Support

You can also run the server in a container:

```bash
docker build -t shipwright-build-mcp-server .
docker run shipwright-build-mcp-server
```

## Development

### Running from Source

```bash
go run main.go
```

### Adding New Tools

To add new tools, modify `main.go`:

1. Define parameter struct for the tool
2. Implement the tool function  
3. Register the tool with the server using `mcp.AddTool()`

## Integration with Shipwright Build

This server is designed to work with Shipwright Build v0.13.0 and later. It uses the `v1beta1` API version which is the current stable version.

### Supported Resources

- **Build** (`builds.shipwright.io/v1beta1`) - Build definitions
- **BuildRun** (`buildruns.shipwright.io/v1beta1`) - Build execution instances  
- **BuildStrategy** (`buildstrategies.shipwright.io/v1beta1`) - Namespace-scoped build strategies
- **ClusterBuildStrategy** (`clusterbuildstrategies.shipwright.io/v1beta1`) - Cluster-scoped build strategies

## Error Handling

The server provides comprehensive error handling with descriptive messages:

- **Validation Errors** - Missing required parameters, invalid formats
- **Kubernetes Errors** - Resource not found, permission denied, API errors
- **Network Errors** - Connection issues, timeout errors

All errors are returned as MCP error responses with helpful context.

## Contributing

Contributions are welcome! Please ensure that:

1. All new tools include proper parameter validation
2. Error handling follows the existing patterns
3. Documentation is updated for new features
4. Code follows Go best practices

## License

This project follows the same license as the main Shipwright Build project. 