# Changelog

All notable changes to the Shipwright Build MCP Server will be documented in this file.

## [v1.1.0] - 2024-01-24

### ✨ Added
- **create_build** tool for creating new Build resources
  - Support for Git and OCI source types
  - Configurable build strategies (BuildStrategy or ClusterBuildStrategy)
  - Parameter passing and timeout configuration
  - Flexible context directory and revision specification

- **create_buildrun** tool for creating new BuildRun resources
  - **Reference Mode** - Create BuildRuns from existing Build resources
  - **Inline Mode** - Create BuildRuns with complete build specifications
  - Service account configuration
  - Parameter overrides and timeout settings
  - Auto-generated names when not specified

### 🔧 Enhanced
- Updated server version to v1.1.0
- Enhanced documentation with detailed creation examples
- Added comprehensive parameter validation for creation tools
- Improved error handling with better context messages

### 📚 Documentation
- Updated README.md with detailed tool documentation and examples
- Enhanced OVERVIEW.md with creation capabilities and use cases
- Added CHANGELOG.md to track version history
- Provided JSON examples for all creation tools

## [v1.0.0] - 2024-01-23

### ✨ Initial Release
- **list_builds** - List builds in a namespace with filtering options
- **get_build** - Get detailed information about a specific build
- **list_buildruns** - List buildruns in a namespace with filtering options
- **get_buildrun** - Get detailed information about a specific buildrun
- **restart_buildrun** - Restart a buildrun by creating a new one
- **list_buildstrategies** - List namespace-scoped build strategies
- **list_clusterbuildstrategies** - List cluster-scoped build strategies

### 🏗️ Infrastructure
- Go 1.22 with MCP Go SDK integration
- Kubernetes controller-runtime client
- Docker support with multi-stage builds
- Makefile for build automation
- Comprehensive documentation and examples 